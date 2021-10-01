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
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type OpenstackNodeSpec struct {

	// OpenStack Cloud Configuration Profile(OCCP) used by this node.
	Occp OccpRef `json:"openstackCloudConfigurationProfileRef"`
}

type OpenstackNodeStatus struct {

	// Whether configuration is generated or not.
	Generated bool `json:"generated,omitempty"`

	// Generated configration from OCCP.
	DesiredNodeConfiguration map[string]interface{} `json:"desiredNodeConfiguration,omitempty"`

	// Status of OpenStack cluster components for this osknode.
	Status string `json:"status,omitempty"`
}

type OpenstackNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenstackNodeSpec   `json:"spec"`
	Status OpenstackNodeStatus `json:"status,omitempty"`
}
