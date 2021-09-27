package v1alpha1

import (
	v1alpha1 "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
)

type Node struct {
	Name string `yaml:"name"`

	Disabled bool `yaml:"disabled"`

	Type string `yaml:"type"`
}

type KupenstackConfigurationSpec struct {
	DefaultProfile v1alpha1.OccpRef `yaml:"defaultProfile"`

	Nodes []Node `yaml:"nodes"`
}

type KupenstackConfiguration struct {
	ApiVersion string `yaml:"apiVersion"`

	Kind string `yaml:"kind"`

	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`

	Spec KupenstackConfigurationSpec `yaml:"spec"`
}
