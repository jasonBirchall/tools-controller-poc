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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	toolv1alpha1 "github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
)

// ToolsReconciler reconciles a Tools object
type ToolsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tool.analytical-platform.justice,resources=tools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tool.analytical-platform.justice,resources=tools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tool.analytical-platform.justice,resources=tools/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ToolsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	tool := &toolv1alpha1.Tools{}
	err := r.Get(ctx, req.NamespacedName, tool)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Check if the tool already exists
	namespace := tool.Namespace
	for _, t := range tool.Spec.Tool {
		found := &appsv1.Deployment{}
		err = r.Get(ctx, types.NamespacedName{Name: t.Name, Namespace: namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			dep, svc := r.deploymentForTool(t.Name, namespace, t.Version, t.Image)
			log.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = r.Create(ctx, dep)
			if err != nil {
				log.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}

			log.Log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			err = r.Create(ctx, svc)
			if err != nil {
				log.Log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
				return ctrl.Result{}, err
			}

			err = ctrl.SetControllerReference(tool, dep, r.Scheme)
			if err != nil {
				log.Log.Error(err, "Failed to set controller reference on Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ToolsReconciler) deploymentForTool(name, namespace, version, image string) (*appsv1.Deployment, *corev1.Service) {
	switch name {
	case "jupyter":
		return r.jupyterDeploymentForTool(namespace, version, image)
	case "rstudio":
		return r.jupyterDeploymentForTool(namespace, version, image)
	case "airflow":
		return r.jupyterDeploymentForTool(namespace, version, image)
	default:
		return nil, nil
	}
}

func (r *ToolsReconciler) airflowDeploymentForTool(namespace, version, image string) *appsv1.Deployment {
	return nil
}
func (r *ToolsReconciler) rstudioDeploymentForTool(namespace, version, image string) *appsv1.Deployment {
	return nil
}

func (r *ToolsReconciler) jupyterDeploymentForTool(namespace, version, image string) (*appsv1.Deployment, *corev1.Service) {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jupyter",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "jupyter",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "jupyter",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "jupyter",
							Image: "jupyter/minimal-notebook",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8888,
								},
							},
						},
					},
				},
			},
		},
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jupyter",
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 8888,
				},
			},
			Selector: deploy.Spec.Template.Labels,
		},
	}

	return deploy, service
}

// SetupWithManager sets up the controller with the Manager.
func (r *ToolsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolv1alpha1.Tools{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
