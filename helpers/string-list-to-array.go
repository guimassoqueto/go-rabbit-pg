package helpers

import (
	"encoding/json"
	"fmt"
	"log"
)
func convertInterfaceToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))
	for i, val := range input {
		strVal, ok := val.(string)
		if !ok {
			strVal = ""
		}
		output[i] = strVal
	}
	return output
}


func StringifiedArrayToArray(stringifiedArray string) []string {
	var parsedArray interface{}

	err := json.Unmarshal([]byte(stringifiedArray), &parsedArray)
	if err != nil {
		fmt.Println("Error parsing stringified array:", err)
	}

	array, ok := parsedArray.([]interface{})
	if !ok {
		log.Panicf("Error converting to array: %s", stringifiedArray)
	}

	return convertInterfaceToStringSlice(array)
}