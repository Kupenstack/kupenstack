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

package project

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	coreV1 "k8s.io/api/core/v1"
	utilname "k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr coreV1.Namespace) error {
	log := r.Log.WithValues("project", cr.Name)

	osclient, err := r.OS.GetClient("identity")
	if err != nil {
		return err
	}

	createOpts := projects.CreateOpts{
		Name:        r.generateName(cr.Name),
		Description: "kubernetes-namespace=" + cr.Name,
		Tags:        []string{"kupenstack"},
	}

	createResult, err := projects.Create(osclient, createOpts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}
	log.Info(msgCreateSuccessful)

	if cr.Annotations == nil {
		cr.Annotations = make(map[string]string)
	}

	cr.Annotations[ExternalNameAnnotation] = createResult.Name
	cr.Annotations[ExternalIDAnnotation] = createResult.ID

	if !utils.ContainsString(cr.GetFinalizers(), Finalizer) {
		controllerutil.AddFinalizer(&cr, Finalizer)
	}

	err = r.Update(ctx, &cr)
	if err != nil {
		return err
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "KupenstackCreated", "Openstack project mapping created.")
	return nil
}

// Appends passed string with a random string suffix.
func (r *Reconciler) generateName(name string) string {

	generatedName := utilname.SimpleNameGenerator.GenerateName(name + "-")

	osclient, err := r.OS.GetClient("identity")
	if err != nil {
		return ""
	}

	allPages, err := projects.ListAvailable(osclient).AllPages()
	if err != nil {
		return ""
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return ""
	}

	unique := true
	for _, kp := range allProjects {
		if kp.Name == generatedName {
			unique = false
		}
	}

	if unique {
		return generatedName
	} else {
		// try again with func recursion
		return r.generateName(name)
	}

}
