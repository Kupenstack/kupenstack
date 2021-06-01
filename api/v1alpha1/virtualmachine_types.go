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

type VirtualMachineSpec struct {
	Image string `json:"image,omitempty"`

	Flavor string `json:"flavor,omitempty"`

	// +optional
	KeyPair string `json:"keyPair,omitempty"`

	// +optional
	Networks []string `json:"network,omitempty"`
}

type VirtualMachineStatus struct {
	ID string `json:"id,omitempty"`

	State string `json:"state,omitempty"`

	IP string `json:"ip,omitempty"`

	// hostname
	Node string `json:"node,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="NODE",type="string",JSONPath=".status.node"
//+kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.state"
//+kubebuilder:printcolumn:name="NETWORKS(IP)",type="string",JSONPath=".status.ip",priority=1
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:resource:shortName={vm}
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}
