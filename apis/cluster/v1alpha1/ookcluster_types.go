/*
Copyright 2021 The Kupenstack Authors.

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

// +kubebuilder:validation:Enum=compute;control
type NodeType string

const (
	Compute NodeType = "compute"
	Control NodeType = "control"
)

type NodeProfile struct {

	// Name of node. Same as name of Kubernetes node.
	Name string `json:"name,omitempty"`

	// Name of profile to use for this node. Default value: default. For profile
	// in different namespace use dns notation of type profilename.namespace
	// +optional
	Profile string `json:"profile,omitempty"`

	// Type of node: "control" or "compute"
	// +optional
	// +kubebuilder:default=compute
	Type NodeType `json:"type,omitempty"`

	// Default value is false. To remove node from cluster set disabled to true.
	// +optional
	Disabled bool `json:"disabled,omitempty"`
}

type OOKClusterSpec struct {
	Profile string `json:"profile,omitempty"`

	// By default all kubernetes nodes are part of OpenStack cluster with
	// default configurations. Nodes defines a list of all nodes that have overriden
	// default configurations.
	// +optional
	// Nodes []NodeProfile `json:"nodes,omitempty"`
}

type OOKClusterStatus struct {
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type OOKCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OOKClusterSpec   `json:"spec,omitempty"`
	Status OOKClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type OOKClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OOKCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OOKCluster{}, &OOKClusterList{})
}
