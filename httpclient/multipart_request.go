package httpclient

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/mgolfam/gogutils/glog"
	"github.com/mgolfam/gogutils/utils"
)

func MultipartData(config HttpConfig, textFields map[string]string, fileFields map[string]string) (*HttpResponse, error) {
	var hresp HttpResponse
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add text fields
	for key, value := range textFields {
		_ = writer.WriteField(key, value)
	}

	// Add file fields
	for key, filePath := range fileFields {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(key, file.Name())
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	_ = writer.Close()

	request, err := http.NewRequest("POST", config.URL, body)
	if err != nil {
		return nil, err
	}

	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	glog.LogL(glog.DEBUG, "http ->", "POST", config.URL)
	startTime := time.Now()
	response, err := client.Do(request)
	elapsedTime := time.Since(startTime)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	headers := make(map[string]string)
	// Parse the response headers into a map
	for key, values := range response.Header {
		// Combine multiple header values with a comma (or your preferred delimiter)
		headers[key] = values[0]
	}
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		glog.LogL(glog.DEBUG, "Error reading response body:", err)
		return nil, err
	}
	// responseBody := buffer.String()
	hresp = HttpResponse{
		StatusCode:  response.StatusCode,
		Headers:     headers,
		Body:        []byte{},
		ElapsedTime: int64(elapsedTime.Milliseconds()),
	}

	if config.LogResponse {
		glog.LogL(glog.DEBUG, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL, hresp.Body)
	} else {
		glog.LogL(glog.DEBUG, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL)
	}

	// Send the HTTP request
	// Convert the response body to a string and print it

	hresp = HttpResponse{
		StatusCode:  response.StatusCode,
		Headers:     headers,
		Body:        []byte{},
		ElapsedTime: int64(elapsedTime.Milliseconds()),
	}

	if config.Cache {
		hresp.CacheTtl = config.CacheTtl
		hresp.CreatedUnix = utils.NowUnixSeconds()
		hresp.SerializeCache(makeHash(config))
	}

	// if config.LogResponse {
	// 	glog.LogL(glog.DEBUG, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL, hresp.Body)
	// } else {
	// 	glog.LogL(glog.DEBUG, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL)
	// }

	return &hresp, nil
}
