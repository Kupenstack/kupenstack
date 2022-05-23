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

package vn

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr kstypes.VirtualNetwork) error {
	log := r.Log.WithValues("virtualnetwork", cr.Name)

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

	opts := mtu.CreateOptsExt{
		CreateOptsBuilder: createOpts,
		MTU:               int(cr.Spec.Mtu),
	}

	createResult, err := networks.Create(osclient, opts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}
	log.Info(msgCreateSuccessful)

	cidr := cr.Spec.Cidr
	IPversion := 4
	if cidr == "" {

		// TODO: If CIDR is not provided on creation then any 10.*.*.0/24 block is
		// automatically assigned if it does not overlap with any existing VN
		cidr = newRandomIP()

	} else {
		_, ipnetwork, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Error(err, msgCreateFailed)
			return err
		}
		if ipnetwork.IP.To4() == nil {
			IPversion = 6
		}
	}

	aps, err := getAllocationPools(cr.Spec.AllocationPools, cidr)
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}

	dhcp := !cr.Spec.DisableDhcp
	subnetOpts := subnets.CreateOpts{
		Name:            cr.Name,
		NetworkID:       createResult.ID,
		CIDR:            cidr,
		EnableDHCP:      &dhcp,
		IPVersion:       gophercloud.IPVersion(IPversion),
		AllocationPools: aps,
	}

	if cr.Spec.GatewayIP != "" {
		subnetOpts.GatewayIP = &cr.Spec.GatewayIP
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

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "VirtualNetwork created.")
	return nil
}

func newRandomIP() string {
	newSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(newSource)
	return fmt.Sprintf("10.%d.%d.0/24", random.Intn(255), random.Intn(255))
}

func getAllocationPools(aps []kstypes.AllocationPool, cidr string) ([]subnets.AllocationPool, error) {

	var allocationPools []subnets.AllocationPool
	for _, ap := range aps {

		var allocationPool subnets.AllocationPool
		var err error

		if ap.StartIP != "" {
			allocationPool.Start, err = translateTargetIPInCIDR(ap.StartIP, cidr)
		} else {
			allocationPool.Start, err = getStartIPOf(cidr)
		}
		if err != nil {
			return nil, err
		}

		if ap.EndIP != "" {
			allocationPool.End, err = translateTargetIPInCIDR(ap.EndIP, cidr)
		} else {
			allocationPool.End, err = getEndIPOf(cidr)
		}
		if err != nil {
			return nil, err
		}

		allocationPools = append(allocationPools, allocationPool)
	}
	return allocationPools, nil
}
