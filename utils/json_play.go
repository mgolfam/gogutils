package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

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

// Deep map / JSON helpers

// GetByPath retrieves a value from a nested map using a dot-separated path.
// Example: path "user.address.city".
func GetByPath(m map[string]interface{}, path string) (interface{}, bool) {
	if m == nil || path == "" {
		return nil, false
	}
	parts := strings.Split(path, ".")
	current := interface{}(m)
	for _, p := range parts {
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		val, exists := asMap[p]
		if !exists {
			return nil, false
		}
		current = val
	}
	return current, true
}

// SetByPath sets a value on a nested map using a dot-separated path,
// creating intermediate maps as needed.
func SetByPath(m map[string]interface{}, path string, value interface{}) error {
	if m == nil || path == "" {
		return errors.New("map and path are required")
	}
	parts := strings.Split(path, ".")
	current := m
	for i, p := range parts {
		if i == len(parts)-1 {
			current[p] = value
			return nil
		}
		next, ok := current[p].(map[string]interface{})
		if !ok {
			next = make(map[string]interface{})
			current[p] = next
		}
		current = next
	}
	return nil
}

// JSONDiff represents differences between two JSON-like maps.
type JSONDiff struct {
	Added   map[string]interface{}
	Removed map[string]interface{}
	Changed map[string][2]interface{} // [old, new]
}

// DiffJSONMap computes a shallow diff between two maps (no recursion).
func DiffJSONMap(a, b map[string]interface{}) JSONDiff {
	diff := JSONDiff{
		Added:   make(map[string]interface{}),
		Removed: make(map[string]interface{}),
		Changed: make(map[string][2]interface{}),
	}

	for k, av := range a {
		if bv, ok := b[k]; !ok {
			diff.Removed[k] = av
		} else if fmt.Sprintf("%v", av) != fmt.Sprintf("%v", bv) {
			diff.Changed[k] = [2]interface{}{av, bv}
		}
	}
	for k, bv := range b {
		if _, ok := a[k]; !ok {
			diff.Added[k] = bv
		}
	}
	return diff
}

// MergeJSONMap merges src into dst with "last write wins".
func MergeJSONMap(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// Simple validators

// ValidateRequired checks that all given paths exist and are non-empty strings.
func ValidateRequired(m map[string]interface{}, paths []string) error {
	for _, p := range paths {
		v, ok := GetByPath(m, p)
		if !ok || v == nil || fmt.Sprint(v) == "" {
			return fmt.Errorf("required field missing or empty: %s", p)
		}
	}
	return nil
}

// ValidateRegex checks that value at path matches the given regex.
func ValidateRegex(m map[string]interface{}, path, pattern string) error {
	v, ok := GetByPath(m, path)
	if !ok {
		return fmt.Errorf("field not found: %s", path)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if !re.MatchString(fmt.Sprint(v)) {
		return fmt.Errorf("field %s does not match pattern", path)
	}
	return nil
}
