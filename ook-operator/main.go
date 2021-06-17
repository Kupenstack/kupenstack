
package main

import (
	"k8s.io/klog/v2/klogr"
	operator "github.com/kupenstack/kupenstack/ook-operator/pkg"
	"github.com/kupenstack/kupenstack/ook-operator/settings"
)



func main(){

	// Todo: take from arg flags

	settings.Port = ":5000"
	settings.DefaultsDir = "/workspace/ook-operator/settings/"
	settings.ActionsDir = "/workspace/ook-operator/pkg/actions/"
	settings.ConfigDir = "/etc/kupenstack/"
	settings.Log = klogr.New().WithName("ook-operator")

	operator.Serve()
}