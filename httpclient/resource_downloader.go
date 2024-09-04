package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mgolfam/gogutils/glog"
)

// Download downloads a file from a remote URL and saves it to a local file.
func Download(remoteURL string, headers map[string]string, filePath, uagent string) (HttpResponse, error) {
	resp := HttpResponse{}
	// Create or open the local file where the content will be saved
	file, err := os.Create(filePath)
	if err != nil {
		return resp, err
	}
	defer file.Close()

	// Create a new HTTP client with custom headers
	client := &http.Client{}
	request, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		return resp, err
	}

	if uagent == "" {
		request.Header.Set("User-Agent", UserAgent)
	} else {
		request.Header.Set("User-Agent", uagent)
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// Send the HTTP request
	glog.LogL(glog.DEBUG, "http ->", "download", request.URL)
	startTime := time.Now()
	response, err := client.Do(request)
	elapsedTime := time.Since(startTime)
	resp.ElapsedTime = int64(elapsedTime.Milliseconds())
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()

	glog.LogL(glog.DEBUG, "http <-", "download", response.Header)
	// Check if the response status code is not OK (200)
	if response.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	resBody, err := io.ReadAll(response.Body)
	// Copy the response body to the local file
	// _, err = io.Copy(file, response.Body)
	if err != nil {
		return resp, err
	}

	file.Write(resBody)

	if err != nil {
		return resp, err
	}

	return resp, err
}

func saveFile(filePath string) error {
	// fileBytes, err := io.ReadFile(filePath)
	// if err != nil {
	// 	return err
	// }

	// 	w.WriteHeader(http.StatusOK)
	// 	w.Header().Set("Content-Type", "application/octet-stream")
	// 	w.Write(fileBytes)
	// }
	return nil
}
