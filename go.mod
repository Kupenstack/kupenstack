module github.com/kupenstack/kupenstack

go 1.15

replace github.com/kupenstack/kupenstack => ./

require (
	github.com/go-logr/logr v0.4.0
	github.com/gophercloud/gophercloud v0.17.0
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/tools v0.1.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.21.1
	k8s.io/apiserver v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.7.2
)
