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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) delete(ctx context.Context, cr coreV1.Namespace) (bool, error) {
	log := r.Log.WithValues("project", cr.Name)

	osclient, err := r.OS.GetClient("identity")
	if err != nil {
		return false, err
	}

	if utils.ContainsString(getFinalizers(cr), kubernetesFinalizer) {
		return false, nil
	}

	err = projects.Delete(osclient, cr.Annotations[ExternalIDAnnotation]).ExtractErr()
	if ignoreNotFoundError(err) != nil {
		log.Error(err, msgDeleteFailed)
		return false, err
	}
	log.Info(msgDeleteSuccessful)

	if utils.ContainsString(cr.GetFinalizers(), Finalizer) {
		controllerutil.RemoveFinalizer(&cr, Finalizer)
	}

	err = r.Update(ctx, &cr)
	if err != nil {
		log.Error(err, msgFinalizerRemoveFailed)
		return false, err
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "KupenstackDeleted", "Openstack project mapping deleted.")
	return true, nil
}

// returns list of finailzers from namespace spec
func getFinalizers(cr coreV1.Namespace) []string {

	var finalizers []string
	for _, finalizerName := range cr.Spec.Finalizers {
		finalizers = append(finalizers, string(finalizerName))
	}
	return finalizers
}

func ignoreNotFoundError(err error) error {
	if err == nil {
		return nil
	}

	if err.Error() == "Resource not found" {
		return nil
	} else {
		return err
	}
}
