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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kstypes "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
	ook "github.com/kupenstack/kupenstack/ook-operator/pkg/actions"
	"github.com/kupenstack/kupenstack/pkg/k8s"
	"github.com/kupenstack/kupenstack/pkg/openstack"
)

const (
	Healthy                   = "Healthy"
	Unhealthy                 = "Unhealthy"
	HelmSteupInProgress       = "HelmSteupInProgress"
	IngressSteupInProgress    = "IngressSteupInProgress"
	MariadbSteupInProgress    = "MariadbSteupInProgress"
	RabbitMQSteupInProgress   = "RabbitMQSteupInProgress"
	MemcachedSteupInProgress  = "MemcachedSteupInProgress"
	KeystoneSteupInProgress   = "KeystoneSteupInProgress"
	HorizonSteupInProgress    = "HorizonSteupInProgress"
	GlanceSteupInProgress     = "GlanceSteupInProgress"
	ComputeKitSteupInProgress = "ComputeKitSteupInProgress"
)

// Reconciler reconciles a OOK Cluster
type Reconciler struct {
	client.Client
	OS            *openstack.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
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
		return ctrl.Result{RequeueAfter: 2 * time.Second}, err
	}

	if status == Unhealthy {
		err := r.Update(ctx, cr)
		return ctrl.Result{RequeueAfter: 8 * time.Second}, err
	}

	if status == Healthy {

		// prepare client
		_, err := r.OS.GetClient("compute")
		if err != nil && err.Error() == openstack.MsgConnectionFailed {

			opts := &gophercloud.AuthOptions{
				IdentityEndpoint: "http://keystone.kupenstack.svc.cluster.local/v3",
				Username:         "admin",
				Password:         "password",
				DomainName:       "Default",
				TenantName:       "admin",
			}
			newClient, err := openstack.New(opts)
			if err != nil {
				return ctrl.Result{}, err
			}
			*r.OS = *newClient
		}

		return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
	}

	return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
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

func (r *Reconciler) getProfile(ctx context.Context, cr kstypes.OOKCluster) (*unstructured.Unstructured, error) {

	profileName := ""
	profileNamespace := ""

	// Predict Name/Namespace of profile based on FQDN name provided.
	profileFQDN := strings.Split(cr.Spec.Profile, ".")
	if len(profileFQDN) == 1 {

		profileName = profileFQDN[0]
		profileNamespace = cr.Namespace

	} else {
		profileName = strings.Join(profileFQDN[:len(profileFQDN)-1], ".")
		profileNamespace = profileFQDN[len(profileFQDN)-1]
		ns := &corev1.Namespace{}
		err := r.Get(ctx, types.NamespacedName{Name: profileNamespace}, ns)
		if err != nil && errors.IsNotFound(err) {
			// Maybe profile is in same namespace
			profileNamespace = cr.Namespace
			profileName = cr.Spec.Profile
		}
	}

	// Get profile data
	profile := &unstructured.Unstructured{}
	profile.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cluster.kupenstack.io",
		Kind:    "Profile",
		Version: "v1alpha1",
	})
	err := r.Get(ctx, types.NamespacedName{Name: profileName, Namespace: profileNamespace}, profile)
	if err != nil {

		if errors.IsNotFound(err) {
			r.Eventf(&cr, corev1.EventTypeWarning, "ProfileNotFound",
				"Profile having name %v not found in namespace %v.", profileName, profileNamespace)
		}
		return nil, err
	}

	return profile, nil
}

func (r *Reconciler) getDesiredState(ctx context.Context, cr kstypes.OOKCluster, component string) (string, error) {

	profile, err := r.getProfile(ctx, cr)
	if err != nil {
		return "", err
	}

	var config map[string]interface{}

	switch component {
	case "glance":
		config = getGlanceDesiredState(profile)
	case "nova":
		config = getNovaDesiredState(profile)
	case "neutron":
		config = getNeutronDesiredState(profile)
	default:
		return "", fmt.Errorf("Unknown Component")
	}

	configStr, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(configStr), nil
}

func updateComponent(component string, conf string) error {

	resp, err := http.Post(fmt.Sprintf("http://localhost:5000/%v/apply", component), "application/json",
		bytes.NewBuffer([]byte(conf)))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OOKOperator: %v apply request failed.", component)
	}

	return nil
}

func (r *Reconciler) Update(ctx context.Context, cr kstypes.OOKCluster) error {

	resp, err := http.Get("http://localhost:5000/cluster/apply")
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	status, err := getComponentStatus("helm")
	if status == Unhealthy {
		return updateComponent("helm", "")
	}

	status, err = getComponentStatus("ingress")
	if status == Unhealthy {
		return updateComponent("ingress", "")
	}

	status, err = getComponentStatus("mariadb")
	if status == Unhealthy {
		return updateComponent("mariadb", "")
	}

	status, err = getComponentStatus("rabbitmq")
	if status == Unhealthy {
		return updateComponent("rabbitmq", "")
	}

	status, err = getComponentStatus("memcached")
	if status == Unhealthy {
		return updateComponent("memcached", "")
	}

	status, err = getComponentStatus("keystone")
	if status == Unhealthy {
		return updateComponent("keystone", "")
	}

	status, err = getComponentStatus("horizon")
	if status == Unhealthy {
		return updateComponent("horizon", "")
	}

	status, err = getComponentStatus("glance")
	if status == Unhealthy {

		conf, err := r.getDesiredState(ctx, cr, "glance")
		if err != nil {
			return err
		}

		return updateComponent("glance", conf)
	}

	status, err = getComponentStatus("libvirt")
	if status == Unhealthy {
		err = updateComponent("libvirt", "")
		if err != nil {
			return err
		}
		err = updateComponent("placement", "")
		if err != nil {
			return err
		}

		conf, err := r.getDesiredState(ctx, cr, "nova")
		if err != nil {
			return err
		}
		err = updateComponent("nova", conf)
		if err != nil {
			return err
		}

		conf, err = r.getDesiredState(ctx, cr, "neutron")
		if err != nil {
			return err
		}
		return updateComponent("neutron", conf)
	}
	return nil
}
