package oskops

import (
	"os"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/oskops/glance"
	"github.com/kupenstack/kupenstack/oskops/horizon"
	"github.com/kupenstack/kupenstack/oskops/ingress"
	"github.com/kupenstack/kupenstack/oskops/keystone"
	"github.com/kupenstack/kupenstack/oskops/libvirt"
	"github.com/kupenstack/kupenstack/oskops/mariadb"
	"github.com/kupenstack/kupenstack/oskops/memcached"
	"github.com/kupenstack/kupenstack/oskops/neutron"
	"github.com/kupenstack/kupenstack/oskops/nova"
	"github.com/kupenstack/kupenstack/oskops/placement"
	"github.com/kupenstack/kupenstack/oskops/rabbitmq"
	"github.com/kupenstack/kupenstack/pkg/helm"
)

func Start(c k8sclient.Client, kupenstackConfig string) {
	log := ctrl.Log.WithName("kupenstack.oskops")

	cfg, err := ReadKupenStackConfiguration(kupenstackConfig)
	if err != nil {
		log.Error(err, "Unable to read KupenstackConfiguration.")
		os.Exit(1)
	}
	profilename := cfg.Spec.DefaultProfile.Name + "." + cfg.Spec.DefaultProfile.Namespace

	err = helm.AddRepoIfNotExist("osh", "https://charts.kupenstack.io")
	if err != nil {
		log.Error(err, "Unable to add helm repo: https://charts.kupenstack.io.")
		os.Exit(1)
	}

	err = helm.UpdateHelmRepos()
	if err != nil {
		log.Error(err, "Failed to update helm repositories.")
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)
	go ingress.Manage(c, profilename, log)
	go mariadb.Manage(c, profilename, log)
	go rabbitmq.Manage(c, profilename, log)
	go memcached.Manage(c, profilename, log)
	go keystone.Manage(c, profilename, log)
	go glance.Manage(c, profilename, log)
	go horizon.Manage(c, profilename, log)
	go nova.Manage(c, profilename, log)
	go neutron.Manage(c, profilename, log)
	go placement.Manage(c, profilename, log)
	go libvirt.Manage(c, profilename, log)
}
