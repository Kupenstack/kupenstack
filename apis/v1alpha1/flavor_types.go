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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FlavorSpec struct {

	// +optional
	VCPU int32 `json:"vcpu,omitempty"`

	// Ram size in MB
	// +optional
	Ram int32 `json:"ram,omitempty"`

	// Disk size in GB
	// +optional
	Disk int32 `json:"disk,omitempty"`

	// +optional
	Swap int32 `json:"swap,omitempty"`

	// +optional
	Rxtx resource.Quantity `json:"rxtx,omitempty"`

	// +optional
	Ephemeral int32 `json:"ephemeral,omitempty"`
}

type FlavorStatus struct {

	// Contains list of all instances using it.
	Usage UsageType `json:"usage"`

	// Unique Id at openstack
	ID string `json:"id,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="IN-USE",type="boolean",JSONPath=".status.usage.inUse"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:resource:scope=Cluster
type Flavor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlavorSpec   `json:"spec,omitempty"`
	Status FlavorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type FlavorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Flavor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Flavor{}, &FlavorList{})
}
