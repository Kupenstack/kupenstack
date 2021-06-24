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

package initializer

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/kupenstack/kupenstack/ook-operator/settings"
)

func Apply(w http.ResponseWriter, r *http.Request) {
	log := settings.Log.WithValues("action", "init")

	cmd := exec.Command(settings.ActionsDir + "init/initCreds")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Error(err, "Failed to apply changes")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	cmd = exec.Command(settings.ActionsDir + "init/initHelm")
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		log.Error(err, "Failed to apply changes")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
