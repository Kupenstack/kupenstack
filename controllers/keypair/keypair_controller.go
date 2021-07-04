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
	"fmt"

	"github.com/go-logr/logr"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kstypes "github.com/kupenstack/kupenstack/apis/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/k8s"
	"github.com/kupenstack/kupenstack/pkg/openstack"
)

const (
	// contains name of keypair resource at openstack.
	ExternalNameAnnotation = "kupenstack.io/external-keypair-name"

	Finalizer = "kupenstack.io/finalizer"
)

// Log messages
const (
	msgCreateFailed           = "Failed to create keypair resource at openstack."
	msgCreateSuccessful       = "Successfully created keypair resource at openstack."
	msgAddControllerRefFailed = "Failed to add controller reference to secret."
	msgCreateSecretFailed     = "Failed to create k8s-secret for keypair."
	msgCreateSecretSuccessful = "Successfully created k8s-secret for keypair."
	msgDeleteFailed           = "Failed to delete keypair resource at openstack."
	msgDeleteSuccessful       = "Successfully deleted keypair resource at openstack."
	msgFinalizerRemoveFailed  = "Failed to remove keypair finalizer at kubernetes."
)

// Reconciler reconciles a KeyPair object
type Reconciler struct {
	client.Client
	OS            *openstack.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=kupenstack.io,resources=keypairs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kupenstack.io,resources=keypairs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kupenstack.io,resources=keypairs/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("keypair", req.NamespacedName)

	var cr kstypes.KeyPair
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create
	if cr.Status.ID == "" {

		err = r.init(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "CreateFailed",
				"Keypair create failed. error: %s", err)
		}
		return ctrl.Result{RequeueAfter: 1000000000}, err
	}

	// delete
	if !cr.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.delete(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "DeleteFailed",
				"Keypair deletion failed. error: %s", err)
		}
		return ctrl.Result{}, err
	}

	if len(cr.Status.Usage.InstanceList) > 0 {
		cr.Status.Usage.InUse = true
	} else {
		cr.Status.Usage.InUse = false
	}

	err = r.Status().Update(ctx, &cr)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled")
	return ctrl.Result{RequeueAfter: 3000000000}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kstypes.KeyPair{}).
		Owns(&coreV1.Secret{}).
		Complete(r)
}

// Records kubernetes event for passed custom resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}
