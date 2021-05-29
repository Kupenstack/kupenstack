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
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func readSecretKey(path string) (string, error) {

	encodedData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	value, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return "", err
	}

	valueStr := string(value)
	valueStr = strings.TrimSuffix(valueStr, "\n")
	valueStr = strings.TrimSpace(valueStr)

	return valueStr, nil
}

func GetClient() (*gophercloud.ProviderClient, error) {

	identityEndpoint, err := readSecretKey("/etc/kupenstack/auth/openstack/authUrl")
	if err != nil {
		return nil, err
	}

	username, err := readSecretKey("/etc/kupenstack/auth/openstack/username")
	if err != nil {
		return nil, err
	}

	password, err := readSecretKey("/etc/kupenstack/auth/openstack/password")
	if err != nil {
		return nil, err
	}

	domain, err := readSecretKey("/etc/kupenstack/auth/openstack/domain")
	if err != nil {
		return nil, err
	}

	tenant, err := readSecretKey("/etc/kupenstack/auth/openstack/tenant")
	if err != nil {
		return nil, err
	}

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		Password:         password,
		DomainName:       domain,
		TenantName:       tenant,
	}

	providerClient, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	return providerClient, nil

}

func GetComputeClient() (*gophercloud.ServiceClient, error) {

	providerClient, err := GetClient()
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetIdentityClient() (*gophercloud.ServiceClient, error) {

	providerClient, err := GetClient()
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewIdentityV3(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetImageServiceClient() (*gophercloud.ServiceClient, error) {

	providerClient, err := GetClient()
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewImageServiceV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
