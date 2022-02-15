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

// ToolsSpec defines the desired state of Tools
type ToolsSpec struct {
	// Name specifies the name of the tool
	// +optional
	Name string `json:"name,omitempty"`

	// Id specifies the id of the tool
	// +optional
	Id string `json:"id,omitempty"`

	// Versions specifies all versions of the tool
	// +optional
	Versions []Version `json:"versions,omitempty"`

	// Description specifies the description of the tool
	// +optional
	Description string `json:"description,omitempty"`

	// Homepage specifies the homepage of the tool
	// +optional
	Homepage string `json:"homepage,omitempty"`

	// Documentation specifies the documentation of the tool
	// +optional
	Documentation string `json:"documentation,omitempty"`

	// Repository specifies the repository of the tool
	// +optional
	Repository string `json:"repository,omitempty"`
}

type Version struct {
	Major int `json:"major,omitempty"`
	Minor int `json:"minor,omitempty"`
	Patch int `json:"patch,omitempty"`
}

// ToolsStatus defines the observed state of Tools
type ToolsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Running indicates whether the tool is running
	// +optional
	Running bool `json:"running,omitempty"`

	// Installed indicates whether the tool is installed in a users namespace
	// +optional
	Installed bool `json:"installed,omitempty"`

	// Version indicates the version of the tool
	// +optional
	Version Version `json:"version,omitempty"`

	// Error indicates the error message of the tool
	// +optional
	Error string `json:"error,omitempty"`

	// ErrorTime indicates the time when the error occurred
	// +optional
	ErrorTime metav1.Time `json:"errorTime,omitempty"`

	// ErrorCount indicates the number of times the tool has failed
	// +optional
	ErrorCount int `json:"errorCount,omitempty"`
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
