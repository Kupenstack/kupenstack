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

package utils

type json = map[string]interface{}

// Default json package does not support merging two json data
// This function assumes json data is stored as map[string]interface{} and
// patches `changes` to `original`
func PatchJson(original, changes json) json {

	if original == nil {
		return changes
	}

	for key, val := range changes {

		if defaultVal, ok := original[key]; ok {
			// key already exists.

			switch defaultVal.(type) {
			case map[string]interface{}:
				// When it has more nested values then we don't overwrite
				original[key] = PatchJson(defaultVal.(json), val.(json))
				break
			default:
				original[key] = val
			}

		} else {
			// key was not in `original` so we simply store it.
			original[key] = val
		}
	}
	return original
}
