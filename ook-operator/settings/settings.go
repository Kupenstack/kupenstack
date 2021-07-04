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

package settings

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Port to server at
var Port string

// Default configuration dir
var DefaultsDir string

// Final Configuration dir. Files in this are the pulled when any
// openstack-helm automation script/plugins are executed.
var ConfigDir string

// Directory containing all executables to automate openstack-helm
var ActionsDir string

var Log logr.Logger

var K8s client.Client