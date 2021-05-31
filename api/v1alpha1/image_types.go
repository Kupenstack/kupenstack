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

// +kubebuilder:validation:Enum=raw;qcow2;iso;vdi;ami;ari;aki;vhd;vmdk
type DiskFormat string

const (
	RAW   DiskFormat = "raw"
	QCOW2 DiskFormat = "qcow2"
	ISO   DiskFormat = "iso"
	VDI   DiskFormat = "vdi"
	AMI   DiskFormat = "ami"
	ARI   DiskFormat = "ari"
	AKI   DiskFormat = "aki"
	VHD   DiskFormat = "vhd"
	VMDK  DiskFormat = "vmdk"
)

type ImageSpec struct {

	// Source contains url to pull image from.
	// +immutable
	Src string `json:"src,omitempty"`

	// Disk format of the image
	Format string `json:"format,omitempty"`

	// +optional
	ContainerFormat string `json:"containerFormat,omitempty"`

	// Minimum disk size in GB required to boot this image.
	// +optional
	MinDisk int32 `json:"minDisk,omitempty"`

	// Minimum ram size in MB required to boot this image.
	// +optional
	MinRam int32 `json:"minRam,omitempty"`
}

type ImageStatus struct {

	// Contains list of all instances using it.
	Usage UsageType `json:"usage"`

	// Unique Id at openstack
	ID string `json:"id,omitempty"`

	// Image is active or not
	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="IN-USE",type="boolean",JSONPath=".status.usage.inUse"
//+kubebuilder:printcolumn:name="READY",type="boolean",JSONPath=".status.ready"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:resource:scope=Cluster
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
