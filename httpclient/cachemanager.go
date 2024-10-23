package httpclient

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mgolfam/gogutils/filemanager"
)

// Function to convert headers map to a sorted string
func headersToString(headers map[string]string) string {
	var result strings.Builder
	// Extract keys from the map
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	// Sort the keys
	sort.Strings(keys)

	// Iterate over the sorted keys
	for _, key := range keys {
		result.WriteString(key)
		result.WriteString(": ")
		result.WriteString(headers[key])
		result.WriteString("\n")
	}

	return result.String()
}

func makeHash(config HttpConfig) string {
	// Concatenate URL, headers, and request body
	hashBody := sha256.Sum256([]byte(string(config.Body)))
	fullRequestString := config.URL + headersToString(config.Headers) + fmt.Sprintf("%x", hashBody)

	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(fullRequestString))
	return fmt.Sprintf("%x", hash)
}

// SaveToFile serializes the HttpResponse to a JSON file.
func (resp *HttpResponse) SerializeCache(hash string) error {
	path := "http-cache/" + hash + ".json"
	filemanager.MkDir("http-cache/")
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	return filemanager.WriteFileBytes(path, data, 0644)
}

// LoadFromFile deserializes the HttpResponse from a JSON file.
func (resp *HttpResponse) DeserializeCache(hash string) error {
	path := "http-cache/" + hash + ".json"
	filemanager.MkDir("http-cache/")
	data, err := filemanager.ReadFileBytes(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, resp)

	if resp.IsCacheEXpired() {
		filemanager.DeleteFile(path)
		return errors.New("cache has been expired.")
	}

	resp.FromCache = true
	return err
}
