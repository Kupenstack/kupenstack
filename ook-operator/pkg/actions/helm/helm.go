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

package helm

import (
	"net/http"
	"os"
	"os/exec"
	"encoding/json"

	"github.com/kupenstack/kupenstack/ook-operator/settings"
	pkg "github.com/kupenstack/kupenstack/ook-operator/pkg/actions"
)

func Apply(w http.ResponseWriter, r *http.Request) {
	log := settings.Log.WithValues("action", "apply-helm")

	cmd := exec.Command(settings.ActionsDir + "helm/initCreds")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Error(err, "Failed to apply changes")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	cmd = exec.Command(settings.ActionsDir + "helm/initHelm")
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



func Status(w http.ResponseWriter, r *http.Request) {
	log := settings.Log.WithValues("action", "status-helm")

	status := pkg.Status{Status: "Ok"}

	resp, err := http.Get("http://localhost:8879")
	if err != nil {
		log.Error(err, "Helm serve not running.")
		status.Status = "NotOk"
	}

	if resp == nil || resp.StatusCode != 200 {
		status.Status = "NotOk"
	}

	cmd := exec.Command("helm", "list")
	err = cmd.Run()
	if err != nil {
		log.Error(err, "")
		status.Status = "NotOk"
	}

	statusStr, err := json.Marshal(status)
	if err != nil {
		log.Error(err, "")
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return	    
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statusStr)
}