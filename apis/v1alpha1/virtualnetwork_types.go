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

type AllocationPool struct {
	StartIP string `json:"startIP,omitempty"`

	EndIP string `json:"endIP,omitempty"`
}

// VirtualNetworkSpec defines the desired state of VirtualNetwork
type VirtualNetworkSpec struct {

	// CIDR block to use for this private network in IPv4 or IPv6.
	Cidr string `json:"cidr,omitempty"`

	// AllocationPools are a way of telling DHCP agent what IP to use from the cidr block.
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`

	// Statically defined the IP that gateway(router) should be given in this network.
	GatewayIP string `json:"gatewayIP,omitempty"`

	// Does not creates DHCP agent in this network when true.
	DisableDhcp bool `json:"disableDhcp,omitempty"`

	// The maximum transmission unit(MTU) value to address fragementation.
	Mtu int32 `json:"mtu,omitempty"`
}

// VirtualNetworkStatus defines the observed state of VirtualNetwork
type VirtualNetworkStatus struct {

	// Unique Id at openstack
	ID string `json:"id,omitempty"`

	// Whether virtual network is ready for use or not
	Ready bool `json:"ready"`

	// Name of the virtualrouter that this virtualnetwork is connected to.
	//+kubebuilder:default="none"
	GatewayName string `json:"gatewayName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="CIDR",type="string",JSONPath=".spec.cidr",priority=1
//+kubebuilder:printcolumn:name="GATEWAY",type="string",JSONPath=".status.gatewayName"
//+kubebuilder:printcolumn:name="GATEWAY-IP",type="string",JSONPath=".spec.gatewayIP",priority=1
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:resource:scope=Namespaced
//+kubebuilder:resource:shortName={vn}
// VirtualNetwork is the Schema for the virtualnetworks API
type VirtualNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualNetworkSpec   `json:"spec,omitempty"`
	Status VirtualNetworkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VirtualNetworkList contains a list of VirtualNetwork
type VirtualNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualNetwork `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualNetwork{}, &VirtualNetworkList{})
}
