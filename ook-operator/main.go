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

package main

import (
	"os"

	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	operator "github.com/kupenstack/kupenstack/ook-operator/pkg"
	"github.com/kupenstack/kupenstack/ook-operator/settings"
)

func main() {

	// Todo: take from arg flags

	settings.Port = ":5000"
	settings.DefaultsDir = "/workspace/ook-operator/settings/"
	settings.ActionsDir = "/workspace/ook-operator/pkg/actions/"
	settings.ConfigDir = "/etc/kupenstack/"
	settings.Log = klogr.New().WithName("ook-operator")

	k8s, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		settings.Log.Error(err, "failed to start ook-operator")
		os.Exit(1)
	}
	settings.K8s = k8s

	operator.Serve()
}
