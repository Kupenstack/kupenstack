package mariadb

import (
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/pkg/helm"
)

func Manage(c client.Client, profilename string, log logr.Logger) {
	log = log.WithName("mariadb")

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

	vals := map[string]interface{}{
		"pod": map[string]interface{}{
			"replicas": map[string]interface{}{
				"server":  1,
				"ingress": 1,
			},
		},
		"volume": map[string]interface{}{
			"enabled": false,
			"use_local_path_for_single_pod_cluster": map[string]interface{}{
				"enabled": true,
			},
		},
	}

	release, err := helm.GetRelease("mariadb", "kupenstack")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("mariadb", "osh", "mariadb", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
