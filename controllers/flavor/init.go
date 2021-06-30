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

func (r *Reconciler) init(ctx context.Context, cr kstypes.Flavor) error {
	log := r.Log.WithValues("flavor", cr.Name)

	osclient, err := r.OS.GetClient("compute")
	if err != nil {
		return err
	}

	public := true
	disk := int(cr.Spec.Disk)
	swap := int(cr.Spec.Swap)
	ephemeral := int(cr.Spec.Ephemeral)
	createOpts := flavors.CreateOpts{
		Name:       cr.Name,
		RAM:        int(cr.Spec.Ram),
		VCPUs:      int(cr.Spec.VCPU),
		Disk:       &disk,
		Swap:       &swap,
		IsPublic:   &public,
		Ephemeral:  &ephemeral,
		RxTxFactor: cr.Spec.Rxtx.AsApproximateFloat64(),
	}
	createResult, err := flavors.Create(osclient, createOpts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}
	log.Info(msgCreateSuccessful)

	// update status
	cr.Status.ID = createResult.ID
	err = r.Status().Update(ctx, &cr)
	if err != nil {
		return err
	}

	cr.Annotations[ExternalNameAnnotation] = createResult.Name
	if !utils.ContainsString(cr.GetFinalizers(), Finalizer) {
		controllerutil.AddFinalizer(&cr, Finalizer)
	}

	// update spec
	err = r.Update(ctx, &cr)
	if err != nil {
		return err
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "Flavor created.")
	return nil
}
