package libvirt

import (
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/pkg/helm"
)

func Manage(c client.Client, profilename string, log logr.Logger) {
	log = log.WithName("libvirt")

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
		"network": map[string]interface{}{
			"backend": []string{
				"linuxbridge",
			},
		},
		"conf": map[string]interface{}{
			"ceph": map[string]interface{}{
				"enabled": false,
			},
		},
	}

	release, err := helm.GetRelease("libvirt", "kupenstack")
	if err != nil {
		return false, err
	}

	if release == nil {
		result, err := helm.UpgradeRelease("libvirt", "osh", "libvirt", "kupenstack", vals)
		if err != nil {
			return false, err
		}
		if result == nil {
			return false, nil
		}
	}

	return true, nil
}
