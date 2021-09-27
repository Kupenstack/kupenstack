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

type ValuesFile struct{}

type KeystoneReplicas struct {

	// Number of keystone-api pods.
	// +kubebuilder:default=1
	// +optional
	Api int32 `json:"api"`
}

type KeystoneConfiguration struct {

	// Configures number of replicas for each pods.
	// +optional
	Replicas KeystoneReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm keystone chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type HorizonReplicas struct {

	// Number of horizon-server pods.
	// +kubebuilder:default=1
	// +optional
	Server int32 `json:"server"`
}

type HorizonConfiguration struct {

	// Whether to disable this component.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`

	// Configures number of replicas for each pods.
	// +optional
	Replicas HorizonReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm horizon chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type GlanceReplicas struct {

	// Number of glance-api pods.
	// +kubebuilder:default=1
	// +optional
	Api int32 `json:"api"`

	// Number of glance-registry pods.
	// +kubebuilder:default=1
	// +optional
	Registry int32 `json:"registry"`
}

type GlanceConfiguration struct {

	// Whether to disable this component.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`

	// Configures number of replicas for each pods.
	// +optional
	Replicas GlanceReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm glance chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type NovaReplicas struct {

	// Number of Nova api-metadata pods.
	// +kubebuilder:default=1
	// +optional
	Metadata int32 `json:"metadata"`

	// Number of Nova ironic pods.
	// +kubebuilder:default=1
	// +optional
	Ironic int32 `json:"ironic"`

	// Number of Nova placement pods.
	// +kubebuilder:default=1
	// +optional
	Placement int32 `json:"placement"`

	// Number of nova-api-osapi pods.
	// +kubebuilder:default=1
	// +optional
	Osapi int32 `json:"osapi"`

	// Number of Nova conductor pods.
	// +kubebuilder:default=1
	// +optional
	Conductor int32 `json:"conductor"`
}

type NovaConfiguration struct {

	// Whether to disable this component.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`

	// Configures number of replicas for each pods.
	// +optional
	Replicas NovaReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm nova chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type NeutronReplicas struct {

	// Number of neutron-server pods.
	// +kubebuilder:default=1
	// +optional
	Server int32 `json:"server"`

	// Number of neutron-ironic-agent pods.
	// +kubebuilder:default=1
	// +optional
	IronicAgent int32 `json:"ironicAgent"`
}

type NeutronConfiguration struct {

	// Whether to disable this component.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`

	// Configures number of replicas for each pods.
	// +optional
	Replicas NeutronReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm neutron chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type PlacementReplicas struct {

	// Number of placement-api pods.
	// +kubebuilder:default=1
	// +optional
	Api int32 `json:"api"`
}

type PlacementConfiguration struct {

	// Whether to disable this component.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`

	// Configures number of replicas for each pods.
	// +optional
	Replicas PlacementReplicas `json:"replicas"`

	// Reference: Values.conf in openstack-helm placement chart.
	// +kubebuilder:pruning:PreserveUnknownFields
	Conf ValuesFile `json:"conf,omitempty"`
}

type OpenStackCloudConfigurationProfileSpec struct {

	// The parent profile to inherit and override in this definition.
	From string `json:"from,omitempty"`

	// Keystone related confs
	Keystone KeystoneConfiguration `json:"keystone,omitempty"`

	// // Horizon related confs
	Horizon HorizonConfiguration `json:"horizon,omitempty"`

	// // Glance related confs
	Glance GlanceConfiguration `json:"glance,omitempty"`

	// // Nova related confs
	Nova NovaConfiguration `json:"nova,omitempty"`

	// // Neutron related confs
	Neutron NeutronConfiguration `json:"neutron,omitempty"`

	// // Placement related confs
	Placement PlacementConfiguration `json:"placement,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName={occp},scope=Namespaced
type OpenStackCloudConfigurationProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenStackCloudConfigurationProfileSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true
type OpenStackCloudConfigurationProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackCloudConfigurationProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackCloudConfigurationProfile{}, &OpenStackCloudConfigurationProfileList{})
}
