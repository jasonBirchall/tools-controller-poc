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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RStudioSpec defines the desired state of RStudio
type RStudioSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of RStudio. Edit rstudio_types.go to remove/update
	Version string `json:"version"`
	Image   string `json:"image"`
	Name    string `json:"name"`
}

// RStudioStatus defines the observed state of RStudio
type RStudioStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
	Host  string   `json:"host"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RStudio is the Schema for the rstudios API
type RStudio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RStudioSpec   `json:"spec,omitempty"`
	Status RStudioStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RStudioList contains a list of RStudio
type RStudioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RStudio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RStudio{}, &RStudioList{})
}
