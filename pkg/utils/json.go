
package utils

type json = map[string]interface{}


// Default json package does not support merging two json data
// This function assumes json data is stored as map[string]interface{} and
// patches `changes` to `original`
func PatchJson(original, changes json) json {

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
