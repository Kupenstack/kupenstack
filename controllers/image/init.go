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

package image

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imageimport"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	coreV1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr kstypes.Image) error {
	log := r.Log.WithValues("image", cr.Name)

	osclient, err := r.OS.GetClient("image")
	if err != nil {
		return err
	}

	protected := false
	public := images.ImageVisibilityPublic
	if cr.Spec.ContainerFormat == "" {
		cr.Spec.ContainerFormat = "bare"
	}
	createOpts := images.CreateOpts{
		Name:            cr.Name,
		Visibility:      &public,
		Protected:       &protected,
		ContainerFormat: cr.Spec.ContainerFormat,
		DiskFormat:      cr.Spec.Format,
		MinDisk:         int(cr.Spec.MinDisk),
		MinRAM:          int(cr.Spec.MinRam),
	}
	createResult, err := images.Create(osclient, createOpts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}

	importOpts := imageimport.CreateOpts{
		Name: imageimport.WebDownloadMethod,
		URI:  cr.Spec.Src,
	}
	err = imageimport.Create(osclient, createResult.ID, importOpts).ExtractErr()
	if err != nil {
		log.Error(err, msgUploadFailed)
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

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "Image created.")
	return nil
}
