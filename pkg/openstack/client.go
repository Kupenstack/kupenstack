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

package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

const (
	msgInvalidAuthOptions = "must provide non-nil gophercloud.AuthOptions to client.New()"
	MsgConnectionFailed   = "Failed to connect to openstack."
)

// Client is a gophercloud go-client warper that sends creates/list/delete/update requests
// for openstack resources to openstack services/components. It lazily initializes
// new gophercloud ServiceClients at the time they are used.
type Client struct {
	// Todo: cache requests in client.

	provider *gophercloud.ProviderClient

	// Client for each service
	clientList map[string]*gophercloud.ServiceClient
}

// New returns a new Client using the provided openstack authentication config.
func New(config *gophercloud.AuthOptions) (*Client, error) {
	return newClient(config)
}

func newClient(config *gophercloud.AuthOptions) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf(msgInvalidAuthOptions)
	}

	providerClient, err := openstack.AuthenticatedClient(*config)
	if err != nil {
		return nil, err
	}

	c := &Client{
		provider:   providerClient,
		clientList: make(map[string]*gophercloud.ServiceClient),
	}

	return c, nil
}

// client.GetClient() returns valid gophercloud.ServiceClient based on `Type` of service.
// If ServiceClient does not exists then it lazily initializes it.
//
// Valid values for Type:
//   * "compute"
//   * "identity"
//   * "image"
//   * "network"
func (client *Client) GetClient(Type string) (*gophercloud.ServiceClient, error) {

	if client.clientList[Type] != nil {
		return client.clientList[Type], nil
	}

	if client.provider == nil {
		return nil, fmt.Errorf(MsgConnectionFailed)
	}

	var err error

	switch Type {
	case "compute":
		client.clientList[Type], err = openstack.NewComputeV2(client.provider,
			gophercloud.EndpointOpts{})
	case "identity":
		client.clientList[Type], err = openstack.NewIdentityV3(client.provider,
			gophercloud.EndpointOpts{})
	case "image":
		client.clientList[Type], err = openstack.NewImageServiceV2(client.provider,
			gophercloud.EndpointOpts{})
	case "network":
		client.clientList[Type], err = openstack.NewNetworkV2(client.provider,
			gophercloud.EndpointOpts{})
	default:
		return nil, fmt.Errorf(MsgConnectionFailed)
	}

	if err != nil {
		return nil, fmt.Errorf(MsgConnectionFailed)
	}

	return client.clientList[Type], nil
}
