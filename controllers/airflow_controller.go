/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
	toolsv1alpha1 "github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
)

// AirflowReconciler reconciles a Airflow object
type AirflowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=airflows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=airflows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=airflows/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(jason): Modify the Reconcile function to compare the state specified by
// the Airflow object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *AirflowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	airflow := &v1alpha1.Airflow{}
	err := r.Get(ctx, req.NamespacedName, airflow)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("Airflow resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Airflow resource")
		return ctrl.Result{}, nil
	}

	deploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: airflow.Name, Namespace: airflow.Namespace}, deploy)
	log.Log.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deployAirflow(airflow)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Log.Error(err, "Failed to create a new deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, err
	} else if err != nil {
		log.Log.Error(err, "Failed to get deployment")
		return ctrl.Result{}, err
	}

	service := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: airflow.Name, Namespace: airflow.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		svc := r.serviceAirflow(airflow)
		log.Log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)

		err = r.Create(ctx, svc)
		if err != nil {
			log.Log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Check for ingress resource
	ingress := &v1beta1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: airflow.Name, Namespace: airflow.Namespace}, ingress)
	if err != nil && errors.IsNotFound(err) {
		ing := r.ingressAirflow(airflow)
		log.Log.Info("Creating a new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
		err = r.Create(ctx, ing)
		if err != nil {
			log.Log.Error(err, "Failed to create new Ingress")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Log.Error(err, "Failed to get Ingress")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AirflowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Airflow{}).
		Complete(r)
}

func (r *AirflowReconciler) ingressAirflow(m *v1alpha1.Airflow) *v1beta1.Ingress {
	ing := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: m.Name + "." + m.Namespace + ".svc.cluster.local",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: m.Name,
										ServicePort: intstr.FromInt(8787),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Set airflow instance as the owner and controller
	controllerutil.SetControllerReference(m, ing, r.Scheme)
	return ing
}

func (r *AirflowReconciler) serviceAirflow(m *v1alpha1.Airflow) *corev1.Service {
	labels := labelsForAirflow(m.Name)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8787,
					TargetPort: intstr.FromInt(8787),
				},
			},
		},
	}
	// Set airflow instance as the owner and controller
	controllerutil.SetControllerReference(m, svc, r.Scheme)
	return svc
}

func labelsForAirflow(name string) map[string]string {
	return map[string]string{"app": "airflow", "airflow_cr": name}
}

func (r *AirflowReconciler) deployAirflow(m *v1alpha1.Airflow) *appsv1.Deployment {
	ls := labelsForAirflow(m.Name)
	image := m.Spec.Image
	version := m.Spec.Version

	if image == "" {
		image = "apache/airflow"
	}

	if version == "" {
		version = "latest"
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "airflow",
							Image: image + ":" + version,
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8787,
								},
							},
						},
					},
				},
			},
		},
	}
	// Set Airflow instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
