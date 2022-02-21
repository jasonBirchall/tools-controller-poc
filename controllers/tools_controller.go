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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	aptoolsv1 "github.com/jasonbirchall/tools-controller-poc/api/v1"
)

// ToolsReconciler reconciles a Tools object
type ToolsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ap-tools.tools,resources=tools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ap-tools.tools,resources=tools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ap-tools.tools,resources=tools/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tools object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ToolsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Get all tools in the cluster and create a list of all versions.
	tools := &aptoolsv1.ToolsList{}
	if err := r.List(ctx, tools); err != nil {
		return ctrl.Result{}, err
	}

	var (
		latestVersion    string
		previousVersions []string
	)

	for _, tool := range tools.Items {
		if tool.Spec.Version > latestVersion {
			latestVersion = tool.Spec.Version
		} else {
			previousVersions = append(previousVersions, tool.Spec.Version)
		}
	}

	// Amend the latest tool with the list of previous versions.
	for _, tool := range tools.Items {
		if tool.Spec.Version == latestVersion {
			tool.Spec.PreviousVersions = previousVersions
			if err := r.Update(ctx, &tool); err != nil {
				log.Log.Error(err, "failed to update latest tool")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ToolsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aptoolsv1.Tools{}).
		Complete(r)
}
