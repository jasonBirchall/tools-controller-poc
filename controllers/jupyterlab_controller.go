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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
	toolsv1alpha1 "github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
)

// JupyterlabReconciler reconciles a Jupyterlab object
type JupyterlabReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=jupyterlabs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=jupyterlabs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=jupyterlabs/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// +kubebuilder:rbac:groups=v1beta1,resources=ingress,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Jupyterlab object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *JupyterlabReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	jlab := &v1alpha1.Jupyterlab{}
	err := r.Get(ctx, req.NamespacedName, jlab)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("JupyterLab resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Jupyterlab resource")
		return ctrl.Result{}, nil
	}

	deploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: jlab.Name, Namespace: jlab.Namespace}, deploy)
	log.Log.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deployJupyterLabs(jlab)
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
	err = r.Get(ctx, types.NamespacedName{Name: jlab.Name, Namespace: jlab.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		svc := r.serviceJupyterLabs(jlab)
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

	// Update the jupyterlab status with pod names
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(jlab.Namespace),
		client.MatchingLabels(labelsForJupyterlab(jlab.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Log.Error(err, "Failed to list pods", "JupyterLab.Namespace", jlab.Namespace, "JupyterLab.Name", jlab.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	if !reflect.DeepEqual(podNames, jlab.Status.Nodes) {
		jlab.Status.Nodes = podNames
		err := r.Status().Update(ctx, jlab)
		if err != nil {
			log.Log.Error(err, "Failed to update Jupyterlab status")
			return ctrl.Result{}, err
		}
	}

	// Check for ingress resource
	ingress := &v1beta1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: jlab.Name, Namespace: jlab.Namespace}, ingress)
	if err != nil && errors.IsNotFound(err) {
		ing := r.ingressJupyterLabs(jlab)
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
func (r *JupyterlabReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Jupyterlab{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&v1beta1.Ingress{}).
		Complete(r)
}

func (r *JupyterlabReconciler) ingressJupyterLabs(jlab *toolsv1alpha1.Jupyterlab) *v1beta1.Ingress {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jlab.Name,
			Namespace: jlab.Namespace,
			Labels:    labelsForJupyterlab(jlab.Name),
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: jlab.Name + ".cp-2501-1650.cloud-platform.service.justice.gov.uk",
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: jlab.Name,
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(jlab, ingress, r.Scheme)
	return ingress
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}

	return podNames
}

func (r *JupyterlabReconciler) serviceJupyterLabs(m *v1alpha1.Jupyterlab) *corev1.Service {
	port := corev1.ServicePort{
		Name: "jupyterlab",
		Port: 8888,
	}

	selectors := labelsForJupyterlab(m.Name)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}

	ports := []corev1.ServicePort{port}
	service.Spec = corev1.ServiceSpec{
		Ports:    ports,
		Selector: selectors,
	}

	// Set JupyterLab instance as the owner and controller
	ctrl.SetControllerReference(m, service, r.Scheme)
	return service
}

func (r *JupyterlabReconciler) deployJupyterLabs(m *v1alpha1.Jupyterlab) *appsv1.Deployment {
	ls := labelsForJupyterlab(m.Name)
	replicas := m.Spec.Size
	image := m.Spec.Image
	version := m.Spec.Version

	if image == "" {
		image = "jupyterlab"
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
			Replicas: &replicas,
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
							Name:  "jupyterlab",
							Image: image + ":" + version,
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8888,
								},
							},
						},
					},
				},
			},
		},
	}
	// Set Jupyterlab instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForJupyterlab returns the labels for selecting the resources
// belonging to the given jupyterlab CR name.
func labelsForJupyterlab(name string) map[string]string {
	return map[string]string{"app": "jupyterlab", "jupyterlab_cr": name}
}

func selectorsForService(name string) map[string]string {
	return map[string]string{
		"app": name,
	}
}
