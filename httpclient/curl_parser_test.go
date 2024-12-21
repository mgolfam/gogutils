package httpclient

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCurlCommand(t *testing.T) {
	tests := []struct {
		name         string
		curlCommand  string
		expectedConf *HttpConfig
		expectError  bool
	}{
		{
			name:        "Basic GET request",
			curlCommand: `curl https://example.com`,
			expectedConf: &HttpConfig{
				Method:   "GET",
				URL:      "https://example.com",
				Headers:  map[string]string{},
				Body:     nil,
				Timeout:  30 * time.Second,
				UseProxy: false,
			},
			expectError: false,
		},
		{
			name:        "POST request with headers and body",
			curlCommand: `curl -X POST https://example.com -H "Content-Type: application/json" --data '{"key":"value"}'`,
			expectedConf: &HttpConfig{
				Method: "POST",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body:     []byte(`{"key":"value"}`),
				Timeout:  30 * time.Second,
				UseProxy: false,
			},
			expectError: false,
		},
		{
			name:        "GET request with query parameters",
			curlCommand: `curl 'https://example.com?key=value'`,
			expectedConf: &HttpConfig{
				Method:   "GET",
				URL:      "https://example.com?key=value",
				Headers:  map[string]string{},
				Body:     nil,
				Timeout:  30 * time.Second,
				UseProxy: false,
			},
			expectError: false,
		},
		{
			name:        "PUT request with custom header",
			curlCommand: `curl -X PUT https://example.com -H "Authorization: Bearer token"`,
			expectedConf: &HttpConfig{
				Method: "PUT",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
				Body:     nil,
				Timeout:  30 * time.Second,
				UseProxy: false,
			},
			expectError: false,
		},
		{
			name:        "DELETE request with proxy",
			curlCommand: `curl -X DELETE https://example.com --proxy http://proxy.example.com:8080`,
			expectedConf: &HttpConfig{
				Method: "DELETE",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Proxy": "http://proxy.example.com:8080",
				},
				Body:     nil,
				Timeout:  30 * time.Second,
				UseProxy: true,
			},
			expectError: false,
		},
		{
			name:        "POST with URL-encoded data",
			curlCommand: `curl -X POST https://example.com --data-urlencode 'key=value'`,
			expectedConf: &HttpConfig{
				Method: "POST",
				URL:    "https://example.com",
				Headers: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
				Body:     []byte("key=value"),
				Timeout:  30 * time.Second,
				UseProxy: false,
			},
			expectError: false,
		},
		{
			name:         "Request with missing URL",
			curlCommand:  `curl -X GET`,
			expectedConf: nil,
			expectError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf, err := ParseCurlCommand(test.curlCommand)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Compare expected and actual HttpConfig
			if conf.Method != test.expectedConf.Method {
				t.Errorf("Expected Method: %s, Got: %s", test.expectedConf.Method, conf.Method)
			}
			if conf.URL != test.expectedConf.URL {
				t.Errorf("Expected URL: %s, Got: %s", test.expectedConf.URL, conf.URL)
			}
			if !reflect.DeepEqual(conf.Headers, test.expectedConf.Headers) {
				t.Errorf("Expected Headers: %+v, Got: %+v", test.expectedConf.Headers, conf.Headers)
			}
			if string(conf.Body) != string(test.expectedConf.Body) {
				t.Errorf("Expected Body: %s, Got: %s", string(test.expectedConf.Body), string(conf.Body))
			}
			if conf.Timeout != test.expectedConf.Timeout {
				t.Errorf("Expected Timeout: %s, Got: %s", test.expectedConf.Timeout, conf.Timeout)
			}
			if conf.UseProxy != test.expectedConf.UseProxy {
				t.Errorf("Expected UseProxy: %v, Got: %v", test.expectedConf.UseProxy, conf.UseProxy)
			}
		})
	}
}
