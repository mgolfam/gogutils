package utils

import (
	"encoding/json"
	"fmt"

	"github.com/mgolfam/gogutils/glog"
)

func PrintJSON(v interface{}) {
	// Marshal the object to JSON with indentation for readability
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}
	// Print the JSON string
	fmt.Println(string(jsonBytes))
}

func ToJsonString(v interface{}) string {
	// Marshal the object to JSON with indentation for readability
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		glog.LogL(glog.ERROR, err)
		return ""
	}

	return string(jsonBytes)
}
