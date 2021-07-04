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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/extendedserverattributes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
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
	// contains name of virtual machine resource at openstack.
	ExternalNameAnnotation = "kupenstack.io/external-virtual-machine-name"

	Finalizer = "kupenstack.io/finalizer"
)

// Log messages
const (
	msgCreateFailed          = "Failed to create virtual machine at openstack."
	msgCreateSuccessful      = "Successfully created virtual machine at openstack."
	msgDeleteFailed          = "Failed to delete virtual machine at openstack."
	msgDeleteSuccessful      = "Successfully deleted virtual machine at openstack."
	msgFinalizerRemoveFailed = "Failed to remove virtual machine finalizer at kubernetes."
)

// Reconciler reconciles a VirtualMachine object
type Reconciler struct {
	client.Client
	OS            *openstack.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=kupenstack.io,resources=virtualmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kupenstack.io,resources=virtualmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kupenstack.io,resources=virtualmachines/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("virtual-machine", req.NamespacedName)

	requeuePeriod := time.Duration(2000000000)

	var cr kstypes.VirtualMachine
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// create
	if cr.Status.ID == "" {

		err := r.init(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "CreateFailed",
				"Virtual Machine create failed. error: %s", err)
		}
		return ctrl.Result{RequeueAfter: 500000000}, err
	}

	// delete
	if !cr.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.delete(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "DeleteFailed",
				"Vitual Machine deletion failed. error: %s", err)
		}
		return ctrl.Result{}, err
	}

	err = r.updateStatus(ctx, cr)

	log.Info("reconciled")
	return ctrl.Result{RequeueAfter: requeuePeriod}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kstypes.VirtualMachine{}).
		Complete(r)
}

// Records kubernetes event for passed custom resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}

func (r *Reconciler) updateStatus(ctx context.Context, cr kstypes.VirtualMachine) error {

	osclient, err := r.OS.GetClient("compute")
	if err != nil {
		return err
	}

	type serverAttributesExt struct {
		servers.Server
		extendedserverattributes.ServerAttributesExt
	}
	var server serverAttributesExt

	err = servers.Get(osclient, cr.Status.ID).ExtractInto(&server)
	if err != nil {
		return err
	}

	cr.Status.Node = server.Host

	cr.Status.State = server.Status
	if server.Status == "ACTIVE" {
		cr.Status.State = "Running"
	}

	networkStatus := ""
	allPages, err := servers.ListAddresses(osclient, cr.Status.ID).AllPages()
	if err != nil {
		return err
	}
	allNetworks, err := servers.ExtractAddresses(allPages)
	if err != nil {
		return err
	}
	for networkName, allAddresses := range allNetworks {

		networkStatus += networkName + "("
		for i, address := range allAddresses {

			if i > 0 {
				networkStatus += ","
			}
			networkStatus += address.Address
		}
		networkStatus += ") "
	}

	cr.Status.IP = networkStatus

	err = r.Status().Update(ctx, &cr)
	return err
}
