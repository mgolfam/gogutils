package utils

import "fmt"

// FlattenJSON flattens a bson object with nested children into a single-level map
func FlattenJSON(data map[string]interface{}) map[string]interface{} {
	flattened := make(map[string]interface{})

	var flattenHelper func(prefix string, value interface{})
	flattenHelper = func(prefix string, value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			for key, val := range v {
				newKey := fmt.Sprintf("%s%s%s", prefix, key, ".")
				flattenHelper(newKey, val)
			}
		default:
			flattened[prefix[:len(prefix)-1]] = value
		}
	}

	for key, val := range data {
		flattenHelper(key+".", val)
	}

	return flattened
}
