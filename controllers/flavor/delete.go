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

package flavor

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) delete(ctx context.Context, cr kstypes.Flavor) error {
	log := r.Log.WithValues("flavor", cr.Name)

	osclient, err := r.OS.GetClient("compute")
	if err != nil {
		return err
	}

	err = flavors.Delete(osclient, cr.Status.ID).ExtractErr()
	if ignoreNotFoundError(err) != nil {
		log.Error(err, msgDeleteFailed)
		return err
	}
	log.Info(msgDeleteSuccessful)

	if utils.ContainsString(cr.GetFinalizers(), Finalizer) {
		controllerutil.RemoveFinalizer(&cr, Finalizer)
	}

	err = r.Update(ctx, &cr)
	if err != nil {
		log.Error(err, msgFinalizerRemoveFailed)
		return err
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "Deleted", "Flavor deleted.")
	return nil
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
