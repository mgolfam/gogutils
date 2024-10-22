package httpclient

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mgolfam/gogutils/filemanager"
)

func headersToString(headers map[string]string) string {
	var result strings.Builder

	for key, value := range headers {
		result.WriteString(key)
		result.WriteString(": ")
		result.WriteString(value)
		result.WriteString("\n")
	}

	return result.String()
}

func makeHash(config HttpConfig) string {
	// Concatenate URL, headers, and request body
	fullRequestString := config.URL + headersToString(config.Headers) + string(config.Body)

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
