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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ToolsSpec defines the desired state of Tools
type ToolsSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=0
	Id string `json:"id"`

	//+kubebuilder:validation:MinLength=0
	Name string `json:"name"`

	//+kubebuilder:validation:MinLength=0
	Version string `json:"version"`

	PreviousVersions []string `json:"previousVersions,omitempty"`

	//+kubebuilder:validation:MinLength=0
	Image string `json:"image"`
}

// ToolsStatus defines the observed state of Tools
type ToolsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// + optional
	Installs int32 `json:"installs"`

	// + optional
	Created *metav1.Time `json:"created,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Tools is the Schema for the tools API
type Tools struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ToolsSpec   `json:"spec,omitempty"`
	Status ToolsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ToolsList contains a list of Tools
type ToolsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tools `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tools{}, &ToolsList{})
}
