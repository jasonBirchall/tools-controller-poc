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
	"testing"

	"github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
	toolsv1alpha1 "github.com/jasonbirchall/tools-controller-poc/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestToolReconciler_Reconcile(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		ctx context.Context
		req ctrl.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ctrl.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ToolReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToolReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToolReconciler.Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToolReconciler_createRStudio(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		tool *toolsv1alpha1.Tool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1alpha1.RStudio
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ToolReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.createRStudio(tt.args.tool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToolReconciler.createRStudio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToolReconciler_createJupyterLab(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		tool *toolsv1alpha1.Tool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1alpha1.Jupyterlab
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ToolReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.createJupyterLab(tt.args.tool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToolReconciler.createJupyterLab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToolReconciler_createAirflow(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		tool *toolsv1alpha1.Tool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1alpha1.Airflow
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ToolReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.createAirflow(tt.args.tool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToolReconciler.createAirflow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToolReconciler_SetupWithManager(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		mgr ctrl.Manager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ToolReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.SetupWithManager(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("ToolReconciler.SetupWithManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
