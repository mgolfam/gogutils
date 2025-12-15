package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mgolfam/gogutils/glog"

	"github.com/mgolfam/gogutils/utils"

	compression "github.com/mgolfam/gogutils/utils/compression"

	"golang.org/x/net/proxy"
)

const UserAgent = "gogutils_client/v0.2.1"

// HTTPClientConfig contains the configuration for the HTTP client.
type HttpConfig struct {
	Method        string
	URL           string
	Headers       map[string]string
	Body          []byte
	Timeout       time.Duration
	LogResponse   bool
	Cache         bool
	RetrieveCache bool
	CacheTtl      int64
	UseProxy      bool
}

type FormDataField struct {
	Name  string
	Value string
	Text  bool
}

type FormDataConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Timeout time.Duration
	Fields  []FormDataField
}

type HttpResponse struct {
	Address     string
	Method      string
	ElapsedTime int64
	StatusCode  int
	Headers     map[string]string
	Body        []byte
	FromCache   bool
	CreatedUnix int64
	CacheTtl    int64
}

func (resp *HttpResponse) IsCacheEXpired() bool {
	return utils.NowUnixSeconds() > resp.CreatedUnix+resp.CacheTtl
}

// SoapConfig contains the configuration for the SOAP client.
type SoapConfig struct {
	URL     string
	Headers map[string]string
	Body    string
	Timeout time.Duration
	LogSoap bool
}

// SoapResponse represents the response from a SOAP call.
type SoapResponse struct {
	ElapsedTime int64
	StatusCode  int
	Headers     map[string]string
	Body        string
}

func (resp *SoapResponse) IsSuccess() bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func getProxyDialer() (proxy.Dialer, error) {
	// proxyUrl := ""
	// if config.Config.App.Proxy.Port == "" {
	// 	proxyUrl = config.Config.App.Proxy.Host
	// } else {
	// 	proxyUrl = config.Config.App.Proxy.Host + ":" + config.Config.App.Proxy.Port
	// }

	// Create a SOCKS5 proxy dialer
	// if config.Config.App.Proxy.ProtocolType == "socks5" {
	// 	dialer, err := proxy.SOCKS5("tcp", proxyUrl, nil, proxy.Direct)
	// 	if err != nil {
	// 		glog.LogL(glog.ERROR, "Error creating SOCKS5 proxy dialer:", err)
	// 	}
	// 	return dialer, err
	// }

	return nil, nil
}

func getTransport() (*http.Transport, error) {
	dialer, err := getProxyDialer()
	if err != nil {
		return nil, err
	}

	// If no proxy is configured, fall back to the default transport.
	if dialer == nil {
		return http.DefaultTransport.(*http.Transport), nil
	}

	// Create a transport that uses the SOCKS5 proxy
	transport := &http.Transport{Dial: dialer.Dial}
	return transport, nil
}

func SendMultipartFormData(config FormDataConfig) (*HttpResponse, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	config.Headers["Content-Type"] = writer.FormDataContentType()
	// files := make([]*os.File, 0)
	for i := 0; i < len(config.Fields); i++ {
		// for loop over all form data fields.
		field := config.Fields[i]

		// we separate files and text fields.
		if !field.Text {
			// creating form data file in bin
			file, err := os.Open(field.Value)
			if err == nil {
				formFile, eff := writer.CreateFormFile(field.Name, field.Value)
				if eff != nil {
					glog.LogL(glog.ERROR, eff)
				}
				_, eff = io.Copy(formFile, file)

				if eff != nil {
					glog.LogL(glog.ERROR, eff)
				}
			}
			defer file.Close()
		} else {
			// creatnig form data field
			formField, eff := writer.CreateFormField(field.Name)
			if eff != nil {
				glog.LogL(glog.ERROR, eff)
			}

			formField.Write([]byte(field.Value))
		}
	}

	err := writer.Close()
	if err != nil {
		glog.LogL(glog.ERROR, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(config.Method, config.URL, payload)
	if err != nil {
		glog.LogL(glog.ERROR, err)
	}

	// Add custom headers to the request
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Send the HTTP request
	glog.LogL(glog.DEBUG, "http multipart ->", config.Method, config.URL)
	startTime := time.Now()
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	elapsedTime := time.Since(startTime)
	return makeResponse(config.Method, config.URL, response, elapsedTime)
}

// RetryConfig defines retry behavior for higher-level helpers.
// It does NOT change the behavior or signature of SendRequest.
type RetryConfig struct {
	Retries    int
	RetryDelay time.Duration
}

// SendRequestWithRetry wraps SendRequest with simple retry logic based on RetryConfig.
// Existing callers can keep using SendRequest; this is fully opt-in.
func SendRequestWithRetry(config HttpConfig, retry RetryConfig) (*HttpResponse, error) {
	// Normalize retries
	maxRetries := retry.Retries
	if maxRetries < 0 {
		maxRetries = 0
	}

	var lastResp *HttpResponse
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := SendRequest(config)
		lastResp, lastErr = resp, err

		// Break on success or non-retryable error/status.
		if err == nil && resp != nil && (resp.StatusCode < 500 || resp.StatusCode >= 600) {
			break
		}

		if attempt == maxRetries {
			break
		}

		if retry.RetryDelay > 0 {
			time.Sleep(retry.RetryDelay)
		}
	}

	return lastResp, lastErr
}

// SendJSON is a convenience helper that sends a JSON body and sets typical headers.
// It keeps the underlying HttpConfig and SendRequest behavior unchanged.
func SendJSON(method, url string, body interface{}, timeout time.Duration, headers map[string]string) (*HttpResponse, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	cfg := HttpConfig{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    jsonBytes,
		Timeout: timeout,
	}

	return SendRequest(cfg)
}

// DecodeJSONBody decodes the response body (assumed JSON) into v.
func DecodeJSONBody(resp *HttpResponse, v interface{}) error {
	if resp == nil {
		return errors.New("nil HttpResponse")
	}
	if len(resp.Body) == 0 {
		return errors.New("empty response body")
	}
	return json.Unmarshal(resp.Body, v)
}

// GetJSON sends a GET request and decodes a JSON response into v.
func GetJSON(url string, timeout time.Duration, headers map[string]string, v interface{}) (*HttpResponse, error) {
	cfg := HttpConfig{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
		Timeout: timeout,
	}
	resp, err := SendRequest(cfg)
	if err != nil {
		return resp, err
	}
	if v != nil {
		if derr := DecodeJSONBody(resp, v); derr != nil {
			return resp, derr
		}
	}
	return resp, nil
}

// PostJSON sends a POST request with a JSON body and decodes a JSON response into v.
func PostJSON(url string, body interface{}, timeout time.Duration, headers map[string]string, v interface{}) (*HttpResponse, error) {
	resp, err := SendJSON(http.MethodPost, url, body, timeout, headers)
	if err != nil {
		return resp, err
	}
	if v != nil {
		if derr := DecodeJSONBody(resp, v); derr != nil {
			return resp, derr
		}
	}
	return resp, nil
}

// Retry presets for convenience.
var (
	RetryFast    = RetryConfig{Retries: 1, RetryDelay: 100 * time.Millisecond}
	RetryDefault = RetryConfig{Retries: 3, RetryDelay: 500 * time.Millisecond}
	RetrySlow    = RetryConfig{Retries: 5, RetryDelay: 2 * time.Second}
)

// SendJSONWithRetryAndDecode combines JSON send, retry, and decode into v.
func SendJSONWithRetryAndDecode(method, url string, body interface{}, timeout time.Duration, headers map[string]string, retry RetryConfig, v interface{}) (*HttpResponse, error) {
	resp, err := SendJSON(method, url, body, timeout, headers)
	if err != nil || v == nil {
		return resp, err
	}

	// Wrap with retry using the existing higher-level helper.
	resp, err = SendRequestWithRetry(HttpConfig{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    resp.Body,
		Timeout: timeout,
	}, retry)
	if err != nil {
		return resp, err
	}

	if derr := DecodeJSONBody(resp, v); derr != nil {
		return resp, derr
	}
	return resp, nil
}

// URL utilities

// MergeHeaders returns a new map with headers from base overridden by overrides.
func MergeHeaders(base, overrides map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}

// AddQueryParams adds the provided query parameters to a URL string.
// If the URL already has query parameters, they are preserved and merged.
func AddQueryParams(rawURL string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return rawURL, nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

func SendRequest(config HttpConfig) (*HttpResponse, error) {
	var hresp HttpResponse
	requestHash := makeHash(config)
	if config.RetrieveCache {
		err := hresp.DeserializeCache(requestHash)
		if err == nil {
			glog.LogL(glog.ERROR, "http ~cache~", config.Method, config.URL)
			return &hresp, err
		}
	}

	// Create a new HTTP client with a custom timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	if config.UseProxy {
		transport, err := getTransport()
		if err != nil {
			return nil, err
		}

		client.Transport = transport
	}

	// Create a request body reader from the string
	var requestBodyReader io.Reader = nil
	if config.Body != nil {
		requestBodyReader = bytes.NewReader(config.Body)
	}

	// Create an HTTP request based on the configuration
	request, err := http.NewRequest(config.Method, config.URL, requestBodyReader)
	if err != nil {
		return nil, err
	}

	// Add custom headers to the request
	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}

	// Send the HTTP request
	if config.UseProxy {
		glog.LogL(glog.INFO, "proxied", "http ->", config.Method, config.URL)
	} else {
		glog.LogL(glog.INFO, "http ->", config.Method, config.URL)
	}
	startTime := time.Now()
	response, err := client.Do(request)
	elapsedTime := time.Since(startTime)

	if err != nil || response == nil {
		return &hresp, err
	}

	// Parse the response headers into a map
	headers := make(map[string]string)
	if response.Header != nil {
		for key, values := range response.Header {
			// Combine multiple header values with a comma (or your preferred delimiter)
			headers[key] = values[0]
		}
	}

	// Read the response body into a byte slice
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glog.LogL(glog.ERROR, "Error reading response body:", err)
		return nil, err
	}

	// Convert the response body to a string and print it
	// responseBody := string(body)
	responseBody := getBody(headers, body)

	hresp = HttpResponse{
		Address:     config.URL,
		Method:      config.Method,
		StatusCode:  response.StatusCode,
		Headers:     headers,
		Body:        responseBody,
		ElapsedTime: int64(elapsedTime.Milliseconds()),
	}

	if config.Cache {
		hresp.CacheTtl = config.CacheTtl
		hresp.CreatedUnix = utils.NowUnixSeconds()
		hresp.SerializeCache(requestHash)
	}

	if config.LogResponse {
		glog.LogL(glog.INFO, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL, hresp.Body)
	} else {
		glog.LogL(glog.INFO, "http <-", hresp.ElapsedTime, hresp.StatusCode, config.Method, config.URL)
	}

	return &hresp, nil
}

func makeResponse(method string, url string,
	response *http.Response, elapsedTime time.Duration) (*HttpResponse, error) {
	if response == nil {
		return nil, errors.New("http.response is null")
	}
	// Parse the response headers into a map
	headers := make(map[string]string)
	if len(response.Header) > 0 {
		for key, values := range response.Header {
			// Combine multiple header values with a comma (or your preferred delimiter)
			headers[key] = values[0]
		}
	}

	// Read the response body into a byte slice
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glog.LogL(glog.ERROR, "Error reading response body:", err)
		return nil, err
	}

	responseBody := getBody(headers, body)

	// Convert the response body to a string and print it
	// responseBody := string(body)

	hresp := HttpResponse{
		StatusCode:  response.StatusCode,
		Headers:     headers,
		Body:        responseBody,
		ElapsedTime: int64(elapsedTime.Milliseconds()),
	}
	glog.LogL(glog.INFO, "http <-", hresp.ElapsedTime, hresp.StatusCode, method, url, hresp.Body)
	return &hresp, nil
}

func getBodyString(headers map[string]string, body []byte) string {
	// compression check
	contentEncoding := headers["Content-Encoding"]
	if contentEncoding == "" {
		return string(body)
	} else if contentEncoding == "gzip" {
		dec, err := compression.Gunzip(body)
		if err != nil {
			return ""
		}
		return string(dec)
	} else if contentEncoding == "deflate" {
		dec, err := compression.Inflate(body)
		if err != nil {
			return ""
		}
		return string(dec)
	}

	return ""
}

func getBody(headers map[string]string, body []byte) []byte {
	// compression check
	contentEncoding := headers["Content-Encoding"]
	if contentEncoding == "" {
		return body
	} else if contentEncoding == "gzip" {
		dec, err := compression.Gunzip(body)
		if err != nil {
			return nil
		}
		return dec
	} else if contentEncoding == "deflate" {
		dec, err := compression.Inflate(body)
		if err != nil {
			return nil
		}
		return dec
	}

	return nil
}

func mkHeader(headers map[string][]string) map[string]string {
	header := make(map[string]string)
	if len(headers) > 0 {
		for key, values := range headers {
			// Combine multiple header values with a comma (or your preferred delimiter)
			header[key] = strings.Join(values, "")
		}
	}
	return header
}

func SoapCall(config SoapConfig) (*SoapResponse, error) {
	client := &http.Client{
		Timeout: config.Timeout,
	}

	request, err := http.NewRequest("POST", config.URL, bytes.NewBufferString(config.Body))
	if err != nil {
		return nil, err
	}

	// Set SOAPAction header if needed
	if soapAction, ok := config.Headers["SOAPAction"]; ok {
		request.Header.Set("SOAPAction", soapAction)
	}

	// Add custom headers to the request
	for key, value := range config.Headers {
		request.Header.Set(key, value)
	}

	glog.LogL(glog.INFO, "SOAP ->", config.URL)

	startTime := time.Now()
	response, err := client.Do(request)
	if err != nil {
		glog.LogL(glog.ERROR, "Error making SOAP request:", err)
		return nil, err
	}
	defer response.Body.Close()

	elapsedTime := time.Since(startTime)

	// Read the response body into a byte slice
	body, err := io.ReadAll(response.Body)
	if err != nil {
		glog.LogL(glog.ERROR, "Error reading response body:", err)
		return nil, err
	}

	soapResp := SoapResponse{
		Headers:    mkHeader(response.Header),
		StatusCode: response.StatusCode,
		Body:       string(body),
	}
	soapResp.ElapsedTime = int64(elapsedTime.Milliseconds())

	if config.LogSoap {
		glog.LogL(glog.INFO, "SOAP <-", soapResp.ElapsedTime, soapResp.StatusCode, config.URL, soapResp.Body)
	} else {
		glog.LogL(glog.INFO, "SOAP <-", soapResp.ElapsedTime, soapResp.StatusCode, config.URL)
	}

	return &soapResp, nil
}
