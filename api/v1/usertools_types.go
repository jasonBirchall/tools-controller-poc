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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// UserToolsSpec defines the desired state of UserTools
type UserToolsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	Version string `json:"version,omitempty"`

	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$
	Url string `json:"url,omitempty"`
}

// UserToolsStatus defines the observed state of UserTools
type UserToolsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	RunningReplicas int32 `json:"runningReplicas,omitempty"`

	// +kubebuilder:validation:Required
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// +kubebuilder:validation:Required
	FailedReplicas int32 `json:"failedReplicas,omitempty"`

	// +kubebuilder:validation:Required
	PendingReplicas int32 `json:"pendingReplicas,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$
	CreationDate string `json:"creationDate,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$
	Url string `json:"url,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UserTools is the Schema for the usertools API
type UserTools struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserToolsSpec   `json:"spec,omitempty"`
	Status UserToolsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UserToolsList contains a list of UserTools
type UserToolsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserTools `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UserTools{}, &UserToolsList{})
}
