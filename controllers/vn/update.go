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
	// "encoding/json"
	"fmt"

	// "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	// "github.com/rackspace/gophercloud/pagination"
	// coreV1 "k8s.io/api/core/v1"
	// "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	// "github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) update(ctx context.Context, cr kstypes.VirtualNetwork) error {
	log := r.Log.WithValues("virtualnetwork", cr.Name)

	osclient, err := r.OS.GetClient("network")
	if err != nil {
		return err
	}

	network, err := networks.Get(osclient, cr.Status.ID).Extract()
	subnet, err := subnets.Get(osclient, network.Subnets[0]).Extract()

	if network.Status == "ACTIVE" {
		cr.Status.Ready = true
	} else {
		// DOWN, BUILD, ERROR
		cr.Status.Ready = false
	}

	err = r.Status().Update(ctx, &cr)
	if err != nil {
		return err
	}

	if cr.Spec.Cidr == "" || cr.Spec.GatewayIP == "" {
		cr.Spec.Cidr = subnet.CIDR
		cr.Spec.GatewayIP = subnet.GatewayIP

		err = r.Update(ctx, &cr)
		if err != nil {
			return err
		}
	}

	if cr.Spec.Mtu == 0 {

		body, ok := networks.Get(osclient, cr.Status.ID).Result.Body.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Invalid response for network type")
		}
		n, ok := body["network"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Invalid response for network type")
		}

		mtu, _ := n["mtu"].(float64)
		cr.Spec.Mtu = int32(mtu)

		err = r.Update(ctx, &cr)
		if err != nil {
			return err
		}
	}

	log.Info("reconciled")
	return nil
}
