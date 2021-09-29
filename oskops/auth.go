package oskops

import (
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/kupenstack/kupenstack/pkg/openstack"
)

// Temporary fix to dynamically authenticate to openstack.
func AuthenticateOpenstackClient(OSclient *openstack.Client) {

	for {
		time.Sleep(20 * time.Second)

		gopheropts := &gophercloud.AuthOptions{
			IdentityEndpoint: "http://keystone.kupenstack.svc.cluster.local/v3",
			Username:         "admin",
			Password:         "password",
			DomainName:       "Default",
			TenantName:       "admin",
		}
		newClient, err := openstack.New(gopheropts)
		if err != nil {
			continue
		}
		*OSclient = *newClient
	}
}
