package helpers

import (
	"encoding/json"
	"fmt"
)

func StringfiedJsonKeysToArray(stringfiedJSON string) []string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(stringfiedJSON), &data)
	FailOnError(err, fmt.Sprintf("Error parsing JSON: %s", err))
	
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}