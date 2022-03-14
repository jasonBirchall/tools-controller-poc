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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestJupyterlabReconciler_Reconcile(t *testing.T) {
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
			r := &JupyterlabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("JupyterlabReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JupyterlabReconciler.Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJupyterlabReconciler_SetupWithManager(t *testing.T) {
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
			r := &JupyterlabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if err := r.SetupWithManager(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("JupyterlabReconciler.SetupWithManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJupyterlabReconciler_ingressJupyterLabs(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		jlab *toolsv1alpha1.Jupyterlab
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *v1beta1.Ingress
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JupyterlabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.ingressJupyterLabs(tt.args.jlab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JupyterlabReconciler.ingressJupyterLabs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJupyterlabReconciler_serviceJupyterLabs(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		m *v1alpha1.Jupyterlab
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *corev1.Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JupyterlabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.serviceJupyterLabs(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JupyterlabReconciler.serviceJupyterLabs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJupyterlabReconciler_deployJupyterLabs(t *testing.T) {
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
	}
	type args struct {
		m *v1alpha1.Jupyterlab
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *appsv1.Deployment
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JupyterlabReconciler{
				Client: tt.fields.Client,
				Scheme: tt.fields.Scheme,
			}
			if got := r.deployJupyterLabs(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JupyterlabReconciler.deployJupyterLabs() = %v, want %v", got, tt.want)
			}
		})
	}
}
