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

package project

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kupenstack/kupenstack/pkg/k8s"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kupenstack/kupenstack/pkg/openstack"
)

const (
	// contains name of project/tenant at openstack.
	ExternalNameAnnotation = "kupenstack.io/external-project-name"

	// contains id of project/tenant at openstack.
	ExternalIDAnnotation = "kupenstack.io/external-project-id"

	Finalizer = "kupenstack.io/finalizer"

	kubernetesFinalizer = "kubernetes"
)

// Log messages
const (
	msgCreateFailed          = "Failed to create project/tenant at openstack."
	msgCreateSuccessful      = "Successfully created project/tenant at openstack."
	msgDeleteFailed          = "Failed to delete project/tenant at openstack."
	msgDeleteSuccessful      = "Successfully deleted project/tenant at openstack."
	msgFinalizerRemoveFailed = "Failed to remove kupenstack finalizer at kubernetes namespace."
)

// Reconciler reconciles a KeyPair object
type Reconciler struct {
	client.Client
	OS            *openstack.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=v1,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=v1,resources=namespaces/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("project", req.NamespacedName.Name)

	var cr coreV1.Namespace
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if cr.Annotations[ExternalNameAnnotation] == "" || cr.Annotations[ExternalIDAnnotation] == "" {
		err = r.init(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "KupenstackCreateFailed",
				"Openstack project create failed. error: %s", err)
		}
		return ctrl.Result{}, err
	}

	if !cr.ObjectMeta.DeletionTimestamp.IsZero() {
		ok, err := r.delete(ctx, cr)
		if err != nil {
			r.Eventf(&cr, coreV1.EventTypeWarning, "KupenstackDeleteFailed",
				"Openstack project deletion failed. error: %s", err)
		}
		if !ok {
			return ctrl.Result{RequeueAfter: 500000000}, err
		}
		return ctrl.Result{}, err
	}

	log.Info("reconciled")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {

	c, err := controller.New("kupenstack-controller", mgr, controller.Options{
		Reconciler: r,
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &coreV1.Namespace{}}, &handler.EnqueueRequestForObject{})
	return err
}

// Records kubernetes event for passed kubernetes resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}
