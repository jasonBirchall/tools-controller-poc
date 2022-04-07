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
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
	toolsv1alpha1 "github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
)

// ToolReconciler reconciles a Tool object
type ToolReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=tools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=tools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.analytical-platform.justice.gov.uk,resources=tools/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(jasonBirchall): Modify the Reconcile function to compare the state specified by
// the Tool object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ToolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	tool := &toolsv1alpha1.Tool{}
	if err := r.Get(ctx, req.NamespacedName, tool); err != nil {
		if errors.IsNotFound(err) {
			log.Log.Info("Airflow resource not found")
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "Failed to get Airflow resource")
		return ctrl.Result{}, nil
	}

	if tool.Name == "airflow" {
		log.Log.Info("Reconciling Airflow")
		airflow := &v1alpha1.Airflow{}
		err := r.Get(ctx, types.NamespacedName{Name: tool.Name, Namespace: tool.Namespace}, airflow)
		if err != nil && errors.IsNotFound(err) {
			af := r.createAirflow(tool)
			if err := r.Create(ctx, af); err != nil {
				log.Log.Error(err, "Failed to create Airflow resource")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}

	if tool.Name == "jupyterlabs" {
		log.Log.Info("Reconciling Jupyterlabs")
		jupyterLab := &v1alpha1.Jupyterlab{}
		err := r.Get(ctx, types.NamespacedName{Name: tool.Name, Namespace: tool.Namespace}, jupyterLab)
		if err != nil && errors.IsNotFound(err) {
			af := r.createJupyterLab(tool)
			if err := r.Create(ctx, af); err != nil {
				log.Log.Error(err, "Failed to create RStudio resource")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}

	if tool.Name == strings.ToLower(tool.Name) {
		log.Log.Info("Reconciling rstudio")
		rstudio := &v1alpha1.RStudio{}
		err := r.Get(ctx, types.NamespacedName{Name: tool.Name, Namespace: tool.Namespace}, rstudio)
		if err != nil && errors.IsNotFound(err) {
			af := r.createRStudio(tool)
			if err := r.Create(ctx, af); err != nil {
				log.Log.Error(err, "Failed to create RStudio resource")
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *ToolReconciler) createRStudio(tool *toolsv1alpha1.Tool) *v1alpha1.RStudio {
	rstudio := &v1alpha1.RStudio{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: v1alpha1.RStudioSpec{
			Image: tool.Spec.ImageName,
		},
	}
	ctrl.SetControllerReference(tool, rstudio, r.Scheme)
	return rstudio
}

func (r *ToolReconciler) createJupyterLab(tool *toolsv1alpha1.Tool) *v1alpha1.Jupyterlab {
	jupyterlab := &v1alpha1.Jupyterlab{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: v1alpha1.JupyterlabSpec{
			ImageName: tool.Spec.ImageName,
		},
	}
	ctrl.SetControllerReference(tool, jupyterlab, r.Scheme)
	return jupyterlab
}

func (r *ToolReconciler) createAirflow(tool *toolsv1alpha1.Tool) *v1alpha1.Airflow {
	airflow := &v1alpha1.Airflow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tool.Name,
			Namespace: tool.Namespace,
		},
		Spec: v1alpha1.AirflowSpec{
			Image: tool.Spec.ImageName,
			Name:  tool.Spec.Name + tool.Spec.Version,
		},
	}
	ctrl.SetControllerReference(tool, airflow, r.Scheme)
	return airflow
}

// SetupWithManager sets up the controller with the Manager.
func (r *ToolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.Tool{}).
		Owns(&toolsv1alpha1.Jupyterlab{}).
		Owns(&toolsv1alpha1.Airflow{}).
		Owns(&toolsv1alpha1.RStudio{}).
		Complete(r)
}
