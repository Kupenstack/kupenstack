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

package actions

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"

	"github.com/kupenstack/kupenstack/ook-operator/settings"
	"github.com/kupenstack/kupenstack/pkg/utils"
)

// Converts nested map[interface{}]interface{} datatype to map[string]interface{} datatype
// Temporary fix for reading nested yaml to map[string]interface{} type
func convertToStringInterface(m map[interface{}]interface{}) map[string]interface{} {

	newJson := make(map[string]interface{})
	for key, value := range m {
		switch value.(type) {
		case map[interface{}]interface{}:
			newJson[key.(string)] = convertToStringInterface(value.(map[interface{}]interface{}))
		default:
			newJson[key.(string)] = value
		}
	}
	return newJson
}

func readDefaultConfig(filename string) (map[string]interface{}, error) {

	conf := make(map[interface{}]interface{})

	defaultFile, err := ioutil.ReadFile(settings.DefaultsDir + filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(defaultFile, &conf)
	if err != nil {
		return nil, err
	}

	// Temporary:
	// Because of Bug[1] yamls are read into map[interface{}]interface{}
	// To temporarily fix this we have defined function: convertToStringInterface()
	// Once this issue is resolved we can delete this function and directly use
	// map[string]interface{}
	//
	// [1]: https://github.com/go-yaml/yaml/issues/139
	//
	return convertToStringInterface(conf), nil
}

// Before running OOK-automation scripts/plugins we have to prepare values to
// apply on charts. This functions places appropiate files with settings in
// appropiate location.
//
// Actions:
//  1. Reads default values from `filenames` in default settings directory.
//  2. Reads new Values from http.request
//  3. Patches new values to default Values
//  4. Saves final prepared values to target localtion with same `filename`.
//  Now values are ready for our automation script/plugins to consume.
func PrepareOOKValues(r *http.Request, filenames []string) error {

	args := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	for _, configFile := range filenames {

		defaultConf, err := readDefaultConfig(configFile)
		if err != nil {
			return err
		}

		confName := configFile[:len(configFile)-5]

		var patch map[string]interface{}
		if args[confName] != nil {
			patch = args[confName].(map[string]interface{})
		}

		jsonData := utils.PatchJson(defaultConf, patch)

		yamlData, err := yaml.Marshal(jsonData)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(settings.ConfigDir+configFile, yamlData, 0664) // -rw-rw-r--
		if err != nil {
			return err
		}

		fmt.Printf("\nSetting config: %v\n----\n%v\n----\n", confName, string(yamlData))
	}

	return nil
}
