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

type PrivateKeyType struct {

	// Name of secret in same namespace.
	SecretName string `json:"secret,omitempty"`
}

type UsageType struct {
	InstanceList []string `json:"instances,omitempty"`

	// set to true when keypair is being used by any instance.
	InUse bool `json:"inUse"`
}

type KeyPairSpec struct {

	// Contains public key from ssh key pairs. This is an optional field.
	// if not provided, it will be updated with a automatically generated
	// ssh public key, and private key will be made avaible through secret
	// referenced in status.
	// +optional
	// +immutable
	PublicKey string `json:"publicKey,omitempty"`
}

type KeyPairStatus struct {

	// if publickey is not provided in spec on create, then privatekey field
	// holds reference to k8s-secret that store base64 encoded private key.
	// secrets can be deleted when not needed.
	PrivateKey PrivateKeyType `json:"privateKey,omitempty"`

	// Contains list of all instances using it.
	Usage UsageType `json:"usage"`

	// Unique Id at openstack
	ID string `json:"id,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="IN-USE",type="boolean",JSONPath=".status.usage.inUse"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="PRIVATE-KEY",type="string",JSONPath=".status.privateKey.secret",priority=1
type KeyPair struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyPairSpec   `json:"spec,omitempty"`
	Status KeyPairStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type KeyPairList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeyPair `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeyPair{}, &KeyPairList{})
}
