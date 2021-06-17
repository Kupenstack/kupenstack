
package pkg


import (
	"net/http"
	"github.com/gorilla/mux"

	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/cluster"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/glance"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/horizon"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/ingress"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/initializer"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/keystone"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/libvirt"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/mariadb"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/memcached"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/neutron"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/nova"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/placement"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/rabbitmq"
	"github.com/kupenstack/kupenstack/ook-operator/settings"
)


func Serve() {
	log := settings.Log

	// Todo: override global variables from env

	r :=  mux.NewRouter()
	r.HandleFunc("/init", initializer.Apply)
	r.HandleFunc("/cluster/apply", cluster.Apply)
	r.HandleFunc("/glance/apply", glance.Apply)
	r.HandleFunc("/horizon/apply", horizon.Apply)
	r.HandleFunc("/ingress/apply", ingress.Apply)
	r.HandleFunc("/keystone/apply", keystone.Apply)
	r.HandleFunc("/libvirt/apply", libvirt.Apply)
	r.HandleFunc("/mariadb/apply", mariadb.Apply)
	r.HandleFunc("/memcached/apply", memcached.Apply)
	r.HandleFunc("/neutron/apply", neutron.Apply)
	r.HandleFunc("/nova/apply", nova.Apply)
	r.HandleFunc("/placement/apply", placement.Apply)
	r.HandleFunc("/rabbitmq/apply", rabbitmq.Apply)


	log.Info("starting.. at PORT " + settings.Port)
    log.Error(http.ListenAndServe(settings.Port, r), "")
}
