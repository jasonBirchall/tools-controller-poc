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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aptoolsv1 "github.com/jasonbirchall/tools-controller-poc/api/v1"
)

// UserToolsReconciler reconciles a UserTools object
type UserToolsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ap-tools.tools,resources=usertools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ap-tools.tools,resources=usertools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ap-tools.tools,resources=usertools/finalizers,verbs=update

// Markers for application deployment
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=ingress,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UserTools object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *UserToolsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var userTools aptoolsv1.UserTools
	if err := r.Get(ctx, req.NamespacedName, &userTools); err != nil {
		log.Error(err, "unable to fetch UserTools")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Reconciling UserTools", "UserTools.Namespace", userTools.Namespace, "UserTools.Name", userTools.Name)

	// Create nginx deployment
	err := r.CreateDeployment(userTools)
	if err != nil {
		log.Error(err, "unable to create deployment")
		return ctrl.Result{}, err
	}
	log.Info("Created deployment")

	// Create nginx service
	err = r.CreateService(userTools)
	if err != nil {
		log.Error(err, "unable to create service")
		return ctrl.Result{}, err
	}
	log.Info("Created service")

	userTools.Status.Url = "http://" + userTools.Name + "." + userTools.Namespace + ".svc.cluster.local"

	err = r.Status().Update(ctx, &userTools)
	if err != nil {
		log.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}

	log.Info("Updated status")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserToolsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aptoolsv1.UserTools{}).
		Complete(r)
}

func (r *UserToolsReconciler) CreateDeployment(userTools aptoolsv1.UserTools) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userTools.Name,
			Namespace: userTools.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"deployment": userTools.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"deployment": userTools.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	if err := r.Create(context.Background(), deployment); err != nil {
		return err
	}

	return nil
}

func (r *UserToolsReconciler) CreateService(userTools aptoolsv1.UserTools) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userTools.Name,
			Namespace: userTools.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"deployment": userTools.Name},
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}

	if err := r.Create(context.Background(), service); err != nil {
		return err
	}

	return nil
}
