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
	// contains name of flavor resource at openstack.
	ExternalNameAnnotation = "kupenstack.io/external-flavor-name"

	Finalizer = "kupenstack.io/finalizer"
)

// Log messages
const (
	msgCreateFailed          = "Failed to create flavor resource at openstack."
	msgCreateSuccessful      = "Successfully created flavor resource at openstack."
	msgDeleteFailed          = "Failed to delete flavor resource at openstack."
	msgDeleteSuccessful      = "Successfully deleted flavor resource at openstack."
	msgFinalizerRemoveFailed = "Failed to remove flavor finalizer at kubernetes."
)

// Reconciler reconciles a Flavor object
type Reconciler struct {
	client.Client
	OS            *openstack.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=kupenstack.io,resources=flavors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kupenstack.io,resources=flavors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kupenstack.io,resources=flavors/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("flavor", req.NamespacedName)

	var cr kstypes.Flavor
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create
	if cr.Status.ID == "" {

		err = r.init(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "CreateFailed",
				"Flavor create failed. error: %s", err)
		}
		return ctrl.Result{RequeueAfter: 1000000000}, err
	}

	// delete
	if !cr.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.delete(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "DeleteFailed",
				"Flavor deletion failed. error: %s", err)
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
		For(&kstypes.Flavor{}).
		Complete(r)
}

// Records kubernetes event for passed custom resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}
