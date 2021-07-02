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

package pkg

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/cluster"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/glance"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/horizon"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/ingress"
	"github.com/kupenstack/kupenstack/ook-operator/pkg/actions/helm"
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

	r := mux.NewRouter()
	r.HandleFunc("/helm/apply", helm.Apply)
	r.HandleFunc("/helm/status", helm.Status)

	r.HandleFunc("/cluster/apply", cluster.Apply)
	
	r.HandleFunc("/glance/apply", glance.Apply)
	r.HandleFunc("/glance/status", glance.Status)

	r.HandleFunc("/horizon/apply", horizon.Apply)
	r.HandleFunc("/horizon/status", horizon.Status)

	r.HandleFunc("/ingress/apply", ingress.Apply)
	r.HandleFunc("/ingress/status", ingress.Status)

	r.HandleFunc("/keystone/apply", keystone.Apply)
	r.HandleFunc("/keystone/status", keystone.Status)

	r.HandleFunc("/libvirt/apply", libvirt.Apply)
	r.HandleFunc("/libvirt/status", libvirt.Status)

	r.HandleFunc("/mariadb/apply", mariadb.Apply)
	r.HandleFunc("/mariadb/status", mariadb.Status)

	r.HandleFunc("/memcached/apply", memcached.Apply)
	r.HandleFunc("/memcached/status", memcached.Status)

	r.HandleFunc("/neutron/apply", neutron.Apply)
	r.HandleFunc("/neutron/status", neutron.Status)

	r.HandleFunc("/nova/apply", nova.Apply)
	r.HandleFunc("/nova/status", nova.Status)

	r.HandleFunc("/placement/apply", placement.Apply)
	r.HandleFunc("/placement/status", placement.Status)

	r.HandleFunc("/rabbitmq/apply", rabbitmq.Apply)
	r.HandleFunc("/rabbitmq/status", rabbitmq.Status)

	log.Info("starting.. at PORT " + settings.Port)
	log.Error(http.ListenAndServe(settings.Port, r), "")
}
