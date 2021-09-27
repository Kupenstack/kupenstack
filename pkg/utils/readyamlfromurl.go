package utils

import (
	"bytes"
	"net/http"

	"gopkg.in/yaml.v2"
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

func ReadYamlFromUrl(urlPath string) (map[string]interface{}, error) {

	conf := make(map[interface{}]interface{})

	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	err = yaml.Unmarshal(buf.Bytes(), &conf)
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
