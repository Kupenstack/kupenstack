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
	"fmt"
	"net/url"
	"strings"

	"github.com/go-logr/logr"
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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterv1alpha1 "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
	"github.com/kupenstack/kupenstack/pkg/k8s"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

// OpenstackNodeReconciler reconciles a OpenstackNode object
type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	EventRecorder record.EventRecorder
}

//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=openstacknodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=openstacknodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.kupenstack.io,resources=openstacknodes/finalizers,verbs=update
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("osknode", req.NamespacedName)

	var cr clusterv1alpha1.OpenstackNode
	err := r.Get(ctx, req.NamespacedName, &cr)
	if err != nil {
		return ctrl.Result{RequeueAfter: 20000000000}, client.IgnoreNotFound(err)
	}

	generatedCfg, err := r.generateDesiredNodeConfiguration(ctx, cr.Spec.Occp.Name, cr.Spec.Occp.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Eventf(&cr, corev1.EventTypeWarning, "OCCPNotFound",
				"Required OpenstackCloudConfigurationProfiles not found for osknode %s.", cr.Name)
		}
		return ctrl.Result{RequeueAfter: 20000000000}, err
	}

	osknode := &unstructured.Unstructured{}
	osknode.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cluster.kupenstack.io",
		Kind:    "OpenstackNode",
		Version: "v1alpha1",
	})
	err = r.Get(ctx, req.NamespacedName, osknode)
	if err != nil {
		return ctrl.Result{RequeueAfter: 20000000000}, client.IgnoreNotFound(err)
	}

	status := make(map[string]interface{})
	if osknode.Object["status"] != nil {
		status = osknode.Object["status"].(map[string]interface{})
	}
	status["desiredNodeConfiguration"] = generatedCfg
	status["generated"] = true
	osknode.Object["status"] = status

	err = r.Status().Update(ctx, osknode)
	if err != nil {
		return ctrl.Result{RequeueAfter: 20000000000}, err
	}

	// get list of desired nodelabels from this osknodes
	// and add them to k8snodes.
	labels := getRequiredLables(osknode)
	err = r.addLabelsToK8sNode(ctx, req.NamespacedName, labels)
	if err != nil {
		return ctrl.Result{RequeueAfter: 20000000000}, err
	}

	return ctrl.Result{RequeueAfter: 20000000000}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1alpha1.OpenstackNode{}).
		Watches(&source.Kind{Type: &clusterv1alpha1.OpenStackCloudConfigurationProfile{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

// Records kubernetes event for passed custom resources.
func (r *Reconciler) Eventf(cr metav1.Object, eventtype, reason, messageFmt string, args ...interface{}) error {
	return k8s.RecordEventf(r.EventRecorder, cr, r.Scheme,
		eventtype, reason, fmt.Sprintf(messageFmt, args...))
}

func (r *Reconciler) addLabelsToK8sNode(ctx context.Context, name types.NamespacedName, labels map[string]string) error {

	var cr corev1.Node
	err := r.Get(ctx, name, &cr)
	if err != nil {
		return err
	}

	for key, value := range labels {
		cr.Labels[key] = value
	}

	return r.Update(ctx, &cr)
}

func getRequiredLables(osknode *unstructured.Unstructured) map[string]string {

	labels := make(map[string]string)
	labels["openstack-control-plane"] = ""
	labels["openstack-compute-node"] = ""

	// node role
	metadata := osknode.Object["metadata"].(map[string]interface{})
	annotations := metadata["annotations"].(map[string]interface{})
	if annotations["node-role"] != nil {
		roles := strings.Split(annotations["node-role"].(string), ",")

		for _, role := range roles {
			role = strings.TrimSpace(role)
			if role == "control" {
				labels["openstack-control-plane"] = "enabled"
			}
			if role == "compute" {
				labels["openstack-compute-node"] = "enabled"
			}
		}
	}

	// profile name
	spec := osknode.Object["spec"].(map[string]interface{})
	occp := spec["openstackCloudConfigurationProfileRef"].(map[string]interface{})
	name := occp["name"].(string)
	namespace := occp["namespace"].(string)
	labels["kupenstack-occp"] = name + "." + namespace

	// invidual openstack component is enabled or not
	status := osknode.Object["status"].(map[string]interface{})
	if status["desiredNodeConfiguration"] != nil {
		cfg := status["desiredNodeConfiguration"].(map[string]interface{})

		key, value := isEnabledLabel(cfg, "keystone")
		labels[key] = value
		key, value = isEnabledLabel(cfg, "glance")
		labels[key] = value
		key, value = isEnabledLabel(cfg, "horizon")
		labels[key] = value
		key, value = isEnabledLabel(cfg, "nova")
		labels[key] = value
		key, value = isEnabledLabel(cfg, "neutron")
		labels[key] = value
		key, value = isEnabledLabel(cfg, "placement")
		labels[key] = value
	}

	// TODO: manage labels for linux-bridge, openvswitch
	// temporary fix:
	labels["linuxbridge"] = "enabled"

	return labels
}

func isEnabledLabel(cfg map[string]interface{}, componentName string) (string, string) {
	labelKey := "kupenstack-" + componentName
	labelValue := ""

	if cfg[componentName] != nil {
		componentDetails := cfg[componentName].(map[string]interface{})

		if componentDetails["disable"] != nil {
			if !componentDetails["disable"].(bool) {
				labelValue = "enabled"
			}
		}
	}

	return labelKey, labelValue
}

func (r *Reconciler) getProfileData(ctx context.Context, name, namespace string) (map[string]interface{}, error) {

	occp := &unstructured.Unstructured{}
	occp.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cluster.kupenstack.io",
		Kind:    "OpenStackCloudConfigurationProfile",
		Version: "v1alpha1",
	})
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, occp)
	if err != nil {
		return nil, err
	}
	if occp.Object["spec"] == nil {
		return nil, nil
	}
	data := occp.Object["spec"].(map[string]interface{})
	return data, err
}

func fetchProfileData(urlPath string) (map[string]interface{}, error) {

	occp, err := utils.ReadYamlFromUrl(urlPath)
	if err != nil {
		return nil, err
	}

	data := occp["spec"].(map[string]interface{})
	return data, nil
}

// generateProfileData() is a recurring function. It reads profile data, and
// if the profile has any parent then it recurs. The resulting data is drived by
// overriding parent data values.
func (r *Reconciler) generateProfileData(ctx context.Context, name, namespace string) (map[string]interface{}, error) {

	var data map[string]interface{}
	var err error

	if isUrl(name) {
		data, err = fetchProfileData(name)
	} else {
		data, err = r.getProfileData(ctx, name, namespace)
	}
	if err != nil {
		return nil, err
	}

	// if no parent then return data
	if data["from"] == nil || data["from"] == "" {
		return data, nil
	}

	// else get parent and merge into profile-data
	parent := data["from"].(string)

	if isUrl(parent) {
		name = parent
		// have not changed namespace to preserve it as default namespace for next function recussion.
	} else {
		profileName, profileNamespace := r.predictNameNamespace(parent)
		name = profileName
		if profileNamespace != "" {
			namespace = profileNamespace
		}
	}

	parentData, err := r.generateProfileData(ctx, name, namespace)
	if err != nil {
		return nil, err
	}
	return utils.PatchJson(parentData, data), err
}

// Predict Name/Namespace of profile based on FQDN name provided.
func (r *Reconciler) predictNameNamespace(profile string) (string, string) {

	profileName := profile
	profileNamespace := ""

	// if dot separated name, then split to get namespace.
	profileFQDN := strings.Split(profile, ".")
	if len(profileFQDN) > 1 {
		profileName = strings.Join(profileFQDN[:len(profileFQDN)-1], ".")
		profileNamespace = profileFQDN[len(profileFQDN)-1]
	}

	// check if namespace exists
	if profileNamespace != "" {
		err := r.Get(context.Background(), types.NamespacedName{Name: profileNamespace}, &corev1.Namespace{})
		if err != nil && errors.IsNotFound(err) {
			// Maybe profile is in same namespace
			profileNamespace = ""
			profileName = profile
		}
	}

	return profileName, profileNamespace
}

func (r *Reconciler) generateDesiredNodeConfiguration(ctx context.Context, name, namespace string) (map[string]interface{}, error) {

	data, err := r.generateProfileData(ctx, name, namespace)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	data["keystone"], err = transformKeys(data["keystone"])
	if err != nil {
		return nil, err
	}

	data["glance"], err = transformKeys(data["glance"])
	if err != nil {
		return nil, err
	}

	data["horizon"], err = transformKeys(data["horizon"])
	if err != nil {
		return nil, err
	}

	data["nova"], err = transformKeys(data["nova"], "metadata", "api_metadata", "ironic", "compute_ironic")
	if err != nil {
		return nil, err
	}

	data["neutron"], err = transformKeys(data["neutron"], "ironicAgent", "ironic_agent")
	if err != nil {
		return nil, err
	}

	data["placement"], err = transformKeys(data["placement"])
	if err != nil {
		return nil, err
	}

	return data, nil
}

func transformKeys(data interface{}, args ...string) (map[string]interface{}, error) {

	if len(args)%2 == 1 {
		return nil, fmt.Errorf("transformKeys() must have even number of keys.")
	}

	if data == nil {
		return nil, nil
	}
	conf := data.(map[string]interface{})

	if conf["replicas"] == nil {
		return data.(map[string]interface{}), nil
	}
	replicas := conf["replicas"].(map[string]interface{})

	for i := range args {
		if i%2 == 1 {
			continue
		}
		// replace keyname of i with i+1
		if replicas[args[i]] != nil {
			replicas[args[i+1]] = replicas[args[i]]
			delete(replicas, args[i])
		}
	}

	conf["pod"] = map[string]interface{}{
		"replicas": replicas,
	}
	delete(conf, "replicas")
	return conf, nil
}

func isUrl(testUrl string) bool {
	_, err := url.ParseRequestURI(testUrl)
	if err == nil {
		return true
	}
	return false
}
