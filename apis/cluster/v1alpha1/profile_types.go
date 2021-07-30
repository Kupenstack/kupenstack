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

type ConfFile struct{}

type ProfileSpec struct {

	// Name of parent profile to inherit this profile from.
	// + optional
	From string `json:"from,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceApiConf ConfFile `json:"glanceApiConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceLoggingConf ConfFile `json:"glanceLoggingConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceApiPasteConf ConfFile `json:"glanceApiPasteConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceRegistryConf ConfFile `json:"glanceRegistryConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceRegistryPasteConf ConfFile `json:"glanceRegistryPasteConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlancePolicyConf ConfFile `json:"glancePolicyConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceApiAuditMapConf ConfFile `json:"glanceApiAuditMapConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	GlanceSwiftStoreConf ConfFile `json:"glanceSwiftStoreConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronApiPasteConf ConfFile `json:"neutronApiPasteConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronPolicyConf ConfFile `json:"neutronPolicyConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronConf ConfFile `json:"neutronConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronLoggingConf ConfFile `json:"neutronLoggingConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronApiAuditMapConf ConfFile `json:"neutronApiAuditMapConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronDhcpAgentConf ConfFile `json:"neutronDhcpAgentConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronL3AgentConf ConfFile `json:"neutronL3AgentConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronMetadataAgent ConfFile `json:"neutronMetadataAgent,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronMl2Conf ConfFile `json:"neutronMl2Conf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronOpenvswitchAgentConf ConfFile `json:"neutronOpenvswitchAgentConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NeutronSriovAgentConf ConfFile `json:"neutronSriovAgentConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NovaApiPasteConf ConfFile `json:"novaApiPasteConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NovaPolicyConf ConfFile `json:"novaPolicyConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NovaConf ConfFile `json:"novaConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NovaLoggingConf ConfFile `json:"novaLoggingConf,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	NovaApiAuditMapConf ConfFile `json:"novaApiAuditMapConf,omitempty"`
}

type ProfileStatus struct {

	// List of all nodes that are using this profile for OpenStack.
	NodeList []string `json:"nodes,omitempty"`

	// NodeCount is string value of format: {Nodes using this profile}/{Total no. of nodes}
	NodeCount string `json:"nodeCount,omitempty"`
}

// skipped //// +kubebuilder:printcolumn:name="NODES",type="string",JSONPath=".status.nodeCount"

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="PARENT",type="string",JSONPath=".spec.from"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type Profile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileSpec   `json:"spec,omitempty"`
	Status ProfileStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type ProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Profile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Profile{}, &ProfileList{})
}
