package oskops

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/kupenstack/kupenstack/oskops/apis/v1alpha1"
)

func ReadKupenStackConfiguration(filename string) (v1alpha1.KupenstackConfiguration, error) {

	var cfg v1alpha1.KupenstackConfiguration

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.ApiVersion != "kupenstack.io/v1alpha1" {
		return cfg, fmt.Errorf("Invalid apiVersion in %s for KupenStackConfiguration", filename)
	}

	if cfg.Kind != "KupenstackConfiguration" {
		return cfg, fmt.Errorf("Invalid kind in %s for KupenStackConfiguration", filename)
	}

	return cfg, err
}
