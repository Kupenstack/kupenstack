package nova

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
	log = log.WithName("nova")

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
				if oskNode.Status.DesiredNodeConfiguration["nova"] != nil {
					vals = oskNode.Status.DesiredNodeConfiguration["nova"].(map[string]interface{})
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
	vals["bootstrap"] = map[string]interface{}{
		"wait_for_computes": map[string]interface{}{
			"enabled": true,
		},
	}

	// overrides for enabling separate placement
	vals["manifests"] = map[string]interface{}{
		"cron_job_cell_setup":        false,
		"cron_job_service_cleaner":   false,
		"statefulset_compute_ironic": false,
		"deployment_placement":       false,
		"ingress_placement":          false,
		"job_db_init_placement":      false,
		"job_ks_placement_endpoints": false,
		"job_ks_placement_service":   false,
		"job_ks_placement_user":      false,
		"pdb_placement":              false,
		"secret_keystone_placement":  false,
		"service_ingress_placement":  false,
		"service_placement":          false,
	}

	release, err := helm.GetRelease("nova", "kupenstack")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("nova", "osh", "nova", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
