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

package keypair

import (
	"context"
	"encoding/base64"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilname "k8s.io/apiserver/pkg/storage/names"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

func (r *Reconciler) init(ctx context.Context, cr kstypes.KeyPair) error {
	log := r.Log.WithValues("keypair", cr.Namespace+"/"+cr.Name)

	osclient, err := r.OS.GetClient("compute")
	if err != nil {
		return err
	}

	createOpts := keypairs.CreateOpts{
		Name:      r.generateName(cr.Name),
		PublicKey: cr.Spec.PublicKey,
	}

	createResult, err := keypairs.Create(osclient, createOpts).Extract()
	if err != nil {
		log.Error(err, msgCreateFailed)
		return err
	}
	log.Info(msgCreateSuccessful)

	privateKeySecretName := ""
	if createResult.PrivateKey != "" {

		immutable := true
		privateKey := base64.StdEncoding.EncodeToString([]byte(createResult.PrivateKey))

		secret := coreV1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    cr.Namespace,
				GenerateName: cr.Name + "-",
			},
			Immutable: &immutable,
			Data: map[string][]byte{
				"privateKey": []byte(privateKey),
			},
		}

		err = ctrl.SetControllerReference(&cr, &secret, r.Scheme)
		if err != nil {
			log.Error(err, msgAddControllerRefFailed)
			return err
		}

		err = r.Create(ctx, &secret)
		if err != nil {
			log.Error(err, msgCreateSecretFailed)
			return err
		}
		log.Info(msgCreateSecretSuccessful)
		r.Eventf(&cr, coreV1.EventTypeNormal, "PrivateKeyCreated",
			"Private key for keypair stored in secret %s.", secret.Name)

		privateKeySecretName = secret.Name
	}

	// update status
	cr.Status.ID = createResult.UserID
	cr.Spec.PublicKey = createResult.PublicKey
	cr.Status.PrivateKey.SecretName = privateKeySecretName

	err = r.Status().Update(ctx, &cr)
	if err != nil {
		return err
	}

	cr.Annotations[ExternalNameAnnotation] = createResult.Name
	if !utils.ContainsString(cr.GetFinalizers(), Finalizer) {
		controllerutil.AddFinalizer(&cr, Finalizer)
	}

	// update spec with publickey and external name
	err = r.Update(ctx, &cr)
	if err != nil {
		return err
	}

	r.Eventf(&cr, coreV1.EventTypeNormal, "Created", "Keypair created.")
	return nil
}

// Appends passed string with a random string suffix.
func (r *Reconciler) generateName(name string) string {

	generatedName := utilname.SimpleNameGenerator.GenerateName(name + "-")

	osclient, _ := r.OS.GetClient("compute")
	allPages, _ := keypairs.List(osclient).AllPages()
	allKeyPairs, _ := keypairs.ExtractKeyPairs(allPages)

	unique := true
	for _, kp := range allKeyPairs {
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
