package ingress

import (
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/pkg/helm"
)

func Manage(c client.Client, profilename string, log logr.Logger) {
	log = log.WithName("ingress")

	for {

		ok, err := reconcile(c, profilename)
		if err != nil {
			log.Error(err, "")
		}

		if ok {
			break
		}
		time.Sleep(30 * time.Second)
	}
}

func reconcile(c client.Client, profilename string) (bool, error) {

	vals := map[string]interface{}{
		"deployment": map[string]interface{}{
			"mode": "cluster",
			"type": "DaemonSet",
		},
		"network": map[string]interface{}{
			"host_namespace": true,
		},
		"pod": map[string]interface{}{
			"replicas": map[string]interface{}{
				"error_page": 1,
			},
		},
	}

	release, err := helm.GetRelease("kube-system-ingress", "kube-system")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("kube-system-ingress", "osh", "ingress", "kube-system", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	vals = map[string]interface{}{
		"pod": map[string]interface{}{
			"replicas": map[string]interface{}{
				"ingress":    1,
				"error_page": 1,
			},
		},
	}

	release, err = helm.GetRelease("kupenstack-ingress", "kupenstack")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("kupenstack-ingress", "osh", "ingress", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
