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

package vm

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	coreV1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	kpctrl "github.com/kupenstack/kupenstack/controllers/keypair"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr kstypes.VirtualMachine) error {
	log := r.Log.WithValues("virtual-machine", cr.Namespace+"/"+cr.Name)

	osclient, err := r.OS.GetClient("compute")
	if err != nil {
		return err
	}

	if cr.Spec.Networks == nil {
		cr.Spec.Networks = append(cr.Spec.Networks, "default")
	}

	err = r.createNetworksIfNotExists(ctx, cr.Spec.Networks)
	if err != nil {
		return err
	}

	img, err := r.getImageResource(ctx, cr)
	if err != nil {
		return err
	}
	if !img.Status.Ready {
		return nil
	}

	keypair, err := r.getKeyPairResource(ctx, cr)
	if err != nil {
		return err
	}

	// we postpone vm creation when:
	// 		1. User has not created KeyPair resource.
	//		2. KeyPair resource is not ready.
	if cr.Spec.KeyPair != "" && keypair.Status.ID == "" {
		return nil
	}

	flavor, err := r.getFlavorResource(ctx, cr)
	if err != nil {
		return err
	}
	if flavor.Status.ID == "" {
		return nil
	}

	networks, err := r.getNetworkResource(ctx, cr)
	if err != nil {
		return err
	}
	if networks == nil {
		return nil
	}

	// get network ids
	var serverNetworks []servers.Network
	for _, network := range networks {
		serverNetworks = append(serverNetworks, servers.Network{UUID: network.Status.ID})
	}

	serverCreateOpts := servers.CreateOpts{
		Name:      cr.Name,
		ImageRef:  img.Status.ID,
		FlavorRef: flavor.Status.ID,
		Networks:  serverNetworks,
	}

	createOpts := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
	}

	if keypair.Status.ID != "" && keypair.Annotations != nil {
		createOpts.KeyName = keypair.Annotations[kpctrl.ExternalNameAnnotation]
	}

	createResult, err := servers.Create(osclient, createOpts).Extract()
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

	if !utils.ContainsString(img.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name) {
		img.Status.Usage.InstanceList = append(img.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name)
		r.Status().Update(ctx, &img)
	}

	if !utils.ContainsString(flavor.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name) {
		flavor.Status.Usage.InstanceList = append(flavor.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name)
		r.Status().Update(ctx, &flavor)
	}

	if !utils.ContainsString(keypair.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name) {
		keypair.Status.Usage.InstanceList = append(keypair.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name)
		r.Status().Update(ctx, &keypair)
	}

	for _, network := range networks {

		if !utils.ContainsString(network.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name) {
			network.Status.Usage.InstanceList = append(network.Status.Usage.InstanceList, cr.Namespace+"/"+cr.Name)
			r.Status().Update(ctx, &network)
		}
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "Virtual Machine created.")
	return nil
}

func (r *Reconciler) createNetworksIfNotExists(ctx context.Context, networkList []string) error {

	for _, networkName := range networkList {

		var temp kstypes.Network
		err := r.Get(ctx, types.NamespacedName{Name: networkName}, &temp)

		if apierrors.IsNotFound(err) {

			newNetwork := kstypes.Network{
				ObjectMeta: metav1.ObjectMeta{
					Name: networkName,
					Annotations: map[string]string{
						"kupenstack.io/auto-delete": "enable",
					},
				},
			}

			err := r.Create(ctx, &newNetwork)
			if err != nil {
				return err
			}

		} else if err != nil {
			return err
		}

	}

	return nil
}

func (r *Reconciler) getImageResource(ctx context.Context, cr kstypes.VirtualMachine) (kstypes.Image, error) {

	var img kstypes.Image
	err := r.Get(ctx, types.NamespacedName{Name: cr.Spec.Image}, &img)
	if apierrors.IsNotFound(err) {
		r.Eventf(&cr, coreV1.EventTypeWarning, "ImageNotFound",
			"Image %s not found.", cr.Spec.Image)
		return img, nil
	}
	return img, err
}

func (r *Reconciler) getKeyPairResource(ctx context.Context, cr kstypes.VirtualMachine) (kstypes.KeyPair, error) {

	var keypair kstypes.KeyPair

	if cr.Spec.KeyPair == "" {
		return keypair, nil
	}

	err := r.Get(ctx, types.NamespacedName{Name: cr.Spec.KeyPair, Namespace: cr.Namespace}, &keypair)
	if apierrors.IsNotFound(err) {
		r.Eventf(&cr, coreV1.EventTypeWarning, "KeyPairNotFound",
			"Key pair %s not found.", cr.Spec.KeyPair)
		return keypair, nil
	}
	return keypair, err
}

func (r *Reconciler) getFlavorResource(ctx context.Context, cr kstypes.VirtualMachine) (kstypes.Flavor, error) {

	var flavor kstypes.Flavor
	err := r.Get(ctx, types.NamespacedName{Name: cr.Spec.Flavor}, &flavor)
	if apierrors.IsNotFound(err) {
		r.Eventf(&cr, coreV1.EventTypeWarning, "FlavorNotFound",
			"Flavor %s not found.", cr.Spec.Flavor)
		return flavor, nil
	}

	return flavor, err
}

func (r *Reconciler) getNetworkResource(ctx context.Context, cr kstypes.VirtualMachine) ([]kstypes.Network, error) {

	var networks []kstypes.Network
	for _, networkName := range cr.Spec.Networks {

		var temp kstypes.Network
		err := r.Get(ctx, types.NamespacedName{Name: networkName}, &temp)
		if err != nil {
			return nil, err
		}

		if temp.Status.ID == "" {
			return nil, nil
		} else {
			networks = append(networks, temp)
		}
	}
	return networks, nil
}
