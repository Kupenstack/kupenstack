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

type OccpRef struct {

	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
}

type OpenstackNodeSpec struct {

	// OpenStack Cloud Configuration Profile(OCCP) used by this node.
	// +kubebuilder:validation:Required
	Occp OccpRef `json:"openstackCloudConfigurationProfileRef"`
}

type OpenstackNodeStatus struct {

	// Whether configuration is generated or not.
	Generated bool `json:"generated,omitempty"`

	// Generated configuration from OCCP.
	// +kubebuilder:pruning:PreserveUnknownFields
	DesiredNodeConfiguration ValuesFile `json:"desiredNodeConfiguration,omitempty"`

	// Status of OpenStack cluster components for this osknode.
	Status string `json:"status,omitempty"`
}

// // +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.status"
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ROLES",type="string",JSONPath=".metadata.annotations.node-role"
//+kubebuilder:printcolumn:name="PROFILE",type="string",JSONPath=".spec.openstackCloudConfigurationProfileRef.name"
//+kubebuilder:resource:shortName={osknode,osknodes},scope=Cluster
type OpenstackNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenstackNodeSpec   `json:"spec"`
	Status OpenstackNodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type OpenstackNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenstackNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenstackNode{}, &OpenstackNodeList{})
}
