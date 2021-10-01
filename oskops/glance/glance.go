package glance

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/pkg/helm"
	ksk "github.com/kupenstack/kupenstack/pkg/kupenstack"
	"github.com/kupenstack/kupenstack/pkg/kupenstack/osknode"
)

func Manage(c client.Client, profilename string, log logr.Logger) {
	log = log.WithName("glance")

	for {

		_, err := reconcile(c, profilename)
		if err != nil {
			log.Error(err, "")
			time.Sleep(10 * time.Second)
			continue
		}

		time.Sleep(30 * time.Second)
	}
}

func reconcile(c client.Client, profilename string) (bool, error) {

	ok, err := ksk.OccpExists(c, profilename)
	if !ok || err != nil {
		return ok, err
	}

	osknodeList, err := osknode.GetList(context.Background(), c)
	if err != nil {
		return false, err
	}

	nodesReady := false
	vals := make(map[string]interface{})
	for _, n := range osknodeList.Items {
		oskNode, err := osknode.AsStruct(&n)
		if err != nil {
			return false, err
		}

		occp := oskNode.Spec.Occp.Name + "." + oskNode.Spec.Occp.Namespace
		if occp == profilename {

			nodesReady = oskNode.Status.Generated

			if oskNode.Status.DesiredNodeConfiguration != nil {
				if oskNode.Status.DesiredNodeConfiguration["glance"] != nil {
					vals = oskNode.Status.DesiredNodeConfiguration["glance"].(map[string]interface{})
				}
			}
		}
	}

	if nodesReady == false {
		return false, nil
	}

	vals["storage"] = "pvc"

	release, err := helm.GetRelease("glance", "kupenstack")
	if err != nil {
		return false, err
	}

	// create pv if not exists
	pv := &core.PersistentVolume{}
	err = c.Get(context.Background(), types.NamespacedName{Name: "glance-pv"}, pv)
	if err != nil {
		if errors.IsNotFound(err) {
			pv := &core.PersistentVolume{
				ObjectMeta: metav1.ObjectMeta{
					Name: "glance-pv",
				},
				Spec: core.PersistentVolumeSpec{
					StorageClassName: "general",
					Capacity: core.ResourceList{
						"storage": resource.MustParse("2Gi"),
					},
					AccessModes: []core.PersistentVolumeAccessMode{
						"ReadWriteOnce",
					},
					PersistentVolumeSource: core.PersistentVolumeSource{
						HostPath: &core.HostPathVolumeSource{
							Path: "/mnt/glance",
						},
					},
				},
			}
			err = c.Create(context.Background(), pv)
			if err != nil {
				return false, err
			}
		}
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("glance", "osh", "glance", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
