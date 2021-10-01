package neutron

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/pkg/helm"
	ksk "github.com/kupenstack/kupenstack/pkg/kupenstack"
	"github.com/kupenstack/kupenstack/pkg/kupenstack/osknode"
)

func Manage(c client.Client, profilename string, log logr.Logger) {
	log = log.WithName("neutron")

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
				if oskNode.Status.DesiredNodeConfiguration["neutron"] != nil {
					vals = oskNode.Status.DesiredNodeConfiguration["neutron"].(map[string]interface{})
				}
			}
		}
	}

	if nodesReady == false {
		return false, nil
	}

	vals["network"] = map[string]interface{}{
		"backend": []string{"linuxbridge"},
	}

	release, err := helm.GetRelease("neutron", "kupenstack")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("neutron", "osh", "neutron", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
