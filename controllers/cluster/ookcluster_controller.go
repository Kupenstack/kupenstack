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

package cluster

import (
	"context"
	"net/http"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/gophercloud/gophercloud"

	kstypes "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/k8s"
	"github.com/kupenstack/kupenstack/pkg/openstack"
	ook "github.com/kupenstack/kupenstack/ook-operator/pkg/actions"
)


const(
	Healthy = "Healthy"
	Unhealthy = "Unhealthy"
	HelmSteupInProgress = "HelmSteupInProgress"
	IngressSteupInProgress = "IngressSteupInProgress"
	MariadbSteupInProgress = "MariadbSteupInProgress"
	RabbitMQSteupInProgress = "RabbitMQSteupInProgress"
	MemcachedSteupInProgress = "MemcachedSteupInProgress"
	KeystoneSteupInProgress = "KeystoneSteupInProgress"
	HorizonSteupInProgress = "HorizonSteupInProgress"
	GlanceSteupInProgress = "GlanceSteupInProgress"
	ComputeKitSteupInProgress = "ComputeKitSteupInProgress"
)

// Reconciler reconciles a OOK Cluster
type Reconciler struct {
	client.Client
	OS     *openstack.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=ookclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=ookclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=ookclusters/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("ookcluster", req.NamespacedName)

	var cr kstypes.OOKCluster
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	status, err := GetStatus()
	if err != nil {
		return ctrl.Result{}, err
	}

	if status != cr.Status.Status {
		cr.Status.Status = status
		err = r.Status().Update(ctx, &cr)
		return ctrl.Result{RequeueAfter: 2*time.Second}, err
	}

	if status == Unhealthy {
		err := r.Update(ctx, req)
		return ctrl.Result{RequeueAfter: 8*time.Second}, err
	}

	if status == Healthy {

		// prepare client
		_, err := r.OS.GetClient("compute")
		if err != nil && err.Error() == openstack.MsgConnectionFailed {	

			opts := &gophercloud.AuthOptions{
					  	IdentityEndpoint: "http://keystone.kupenstack.svc.cluster.local/v3",
					  	Username: "admin",
					  	Password: "password",
					  	DomainName: "Default",
					  	TenantName: "admin",
					}
			newClient, err := openstack.New(opts)
			if err != nil{
				return ctrl.Result{}, err
			}
			*r.OS = *newClient
		}

		return ctrl.Result{RequeueAfter: 3*time.Second}, nil
	}

	return ctrl.Result{RequeueAfter: 3*time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kstypes.OOKCluster{}).
		Complete(r)
}

// Records kubernetes event for passed custom resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}


func getComponentStatus(component string) (string, error) {

	s := ook.Status{}

	resp, err := http.Get(fmt.Sprintf("http://localhost:5000/%v/status", component))
	if err != nil {
		return "", err
	}
	json.NewDecoder(resp.Body).Decode(&s)
	if s.Status == "NotOk" || s.Status == "" {
		return Unhealthy, nil
	}
	if s.Status == "InProgress" {

		status := Unhealthy

		switch component {
		case "helm":
			status = HelmSteupInProgress	
		case "ingress":
			status = IngressSteupInProgress
		case "mariadb":
			status = MariadbSteupInProgress
		case "rabbitmq":
			status = RabbitMQSteupInProgress
		case "memcached":
			status = MemcachedSteupInProgress
		case "keystone":
			status = KeystoneSteupInProgress
		case "horizon":
			status = HorizonSteupInProgress
		case "glance":
			status = GlanceSteupInProgress
		case "libvirt":
			status = ComputeKitSteupInProgress
		case "placement":
			status = ComputeKitSteupInProgress
		case "nova":
			status = ComputeKitSteupInProgress
		case "neutron":
			status = ComputeKitSteupInProgress
		default:
			return "", fmt.Errorf("Unknown Status")
		}
		return status, nil
	}

	return "Ok", nil
}

func GetStatus() (string, error) {

	status, err := getComponentStatus("helm")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("ingress")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("mariadb")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("rabbitmq")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("memcached")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("keystone")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("horizon")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("glance")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("libvirt")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("placement")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("nova")
	if status != "Ok" {
		return status, err
	}

	status, err = getComponentStatus("neutron")
	if status != "Ok" {
		return status, err
	}

	return Healthy, nil 
}

func updateComponent(component string) error {

	resp, err := http.Get(fmt.Sprintf("http://localhost:5000/%v/apply", component))
	if err != nil {
		return err
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OOKOperator: %v apply request failed.", component)
	}
	
	return nil
}

func (r *Reconciler) Update(ctx context.Context, req ctrl.Request) error {

	resp, err := http.Get("http://localhost:5000/cluster/apply")
	if err != nil  || resp.StatusCode != http.StatusOK {
		return err
	}

	status, err := getComponentStatus("helm")
	if status == Unhealthy {
		return updateComponent("helm")
	}

	status, err = getComponentStatus("ingress")
	if status == Unhealthy {
		return updateComponent("ingress")
	}

	status, err = getComponentStatus("mariadb")
	if status == Unhealthy {
		return updateComponent("mariadb")
	}

	status, err = getComponentStatus("rabbitmq")
	if status == Unhealthy {
		return updateComponent("rabbitmq")
	}

	status, err = getComponentStatus("memcached")
	if status == Unhealthy {
		return updateComponent("memcached")
	}

	status, err = getComponentStatus("keystone")
	if status == Unhealthy {
		return updateComponent("keystone")
	}

	status, err = getComponentStatus("horizon")
	if status == Unhealthy {
		return updateComponent("horizon")
	}

	status, err = getComponentStatus("glance")
	if status == Unhealthy {
		return updateComponent("glance")
	}

	status, err = getComponentStatus("libvirt")
	if status == Unhealthy {
		err = updateComponent("libvirt")
		if err != nil {
			return err
		}
		err = updateComponent("placement")
		if err != nil {
			return err
		}
		err = updateComponent("nova")
		if err != nil {
			return err
		}
		return updateComponent("neutron")
	}
	return nil
}

