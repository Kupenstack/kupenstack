module github.com/kupenstack/kupenstack

go 1.15

replace github.com/kupenstack/kupenstack => ./

require (
	github.com/go-logr/logr v0.4.0
	github.com/gofrs/flock v0.8.0
	github.com/gophercloud/gophercloud v0.17.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/pkg/errors v0.9.1
	github.com/racker/perigee v0.1.0 // indirect
	github.com/rackspace/gophercloud v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.7.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/apiserver v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.10.1
)
