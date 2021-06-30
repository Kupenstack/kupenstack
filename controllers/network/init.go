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

package network

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr kstypes.Network) error {
	log := r.Log.WithValues("network", cr.Name)

	osclient, err := r.OS.GetClient("network")
	if err != nil {
		return err
	}

	shared := true
	activate := true
	createOpts := networks.CreateOpts{
		Name:         cr.Name,
		AdminStateUp: &activate,
		Shared:       &shared,
	}
	createResult, err := networks.Create(osclient, createOpts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}
	log.Info(msgCreateSuccessful)

	cidr := cr.Spec.Cidr
	if cidr == "" {
		cidr = "10.10.0.0/16"
	}
	dhcp := true
	subnetOpts := subnets.CreateOpts{
		Name:       cr.Name,
		NetworkID:  createResult.ID,
		CIDR:       cidr,
		EnableDHCP: &dhcp,
		IPVersion:  4,
	}

	_, err = subnets.Create(osclient, subnetOpts).Extract()
	if err != nil {
		log.Error(err, msgSubnetCreateFailed)
		return err
	}
	log.Info(msgSubnetCreateSuccessful)

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

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "Network created.")
	return nil
}
