package httpclient

import (
	"errors"
	"strings"
	"time"
)

// ParseCurlCommand parses a curl command string and returns an HttpConfig.
func ParseCurlCommand(curlCommand string) (*HttpConfig, error) {
	// Normalize multi-line commands
	curlCommand = normalizeCurlCommand(curlCommand)

	curlCommand = strings.TrimSpace(curlCommand)
	if !strings.HasPrefix(curlCommand, "curl ") {
		return nil, errors.New("not a valid curl command")
	}

	// Initialize variables
	method := "GET"
	urlAddr := ""
	headers := make(map[string]string)
	var body []byte
	contentTypeSet := false

	// Split the command into tokens
	tokens := splitCurlCommand(curlCommand)

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		switch token {
		case "--location":
			// Ignore the flag; no action needed.
		case "--request", "-X":
			// Extract method
			if i+1 < len(tokens) {
				nextToken := strings.ToUpper(tokens[i+1])
				if isValidMethod(nextToken) {
					method = nextToken
				}
				i++
			}
		case "--header", "-H":
			// Extract headers
			if i+1 < len(tokens) {
				headerParts := strings.SplitN(strings.TrimSpace(tokens[i+1]), ":", 2)
				if len(headerParts) == 2 {
					key := strings.TrimSpace(headerParts[0])
					value := strings.TrimSpace(headerParts[1])
					headers[key] = value
					if strings.ToLower(key) == "content-type" {
						contentTypeSet = true
					}
				}
				i++
			}
		case "--data", "--data-urlencode":
			// Extract body data
			if i+1 < len(tokens) {
				if len(body) > 0 {
					body = append(body, '&')
				}
				body = append(body, strings.Trim(tokens[i+1], `'"`)...)
				i++
				// Automatically set Content-Type if not already set
				if !contentTypeSet {
					headers["Content-Type"] = "application/x-www-form-urlencoded"
					contentTypeSet = true
				}
			}
		default:
			// Assume it's the URL if it starts with http
			if strings.HasPrefix(token, "http") {
				urlAddr = strings.Trim(token, `'"`)
			}
		}
	}

	// Validate the URL
	if urlAddr == "" {
		return nil, errors.New("URL not found in curl command")
	}

	// Create the HttpConfig
	config := &HttpConfig{
		Method:   method,
		URL:      urlAddr,
		Headers:  headers,
		Body:     body,
		Timeout:  30 * time.Second,
		UseProxy: headers["Proxy"] != "",
	}

	return config, nil
}

// normalizeCurlCommand combines multi-line curl commands into a single line.
// normalizeCurlCommand combines multi-line curl commands into a single line.
func normalizeCurlCommand(curlCommand string) string {
	lines := strings.Split(curlCommand, "\n")
	var normalized strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line) // Corrected: Pass the line to TrimSpace
		if strings.HasSuffix(line, "\\") {
			normalized.WriteString(strings.TrimSuffix(line, "\\") + " ")
		} else {
			normalized.WriteString(line + " ")
		}
	}

	return strings.TrimSpace(normalized.String())
}

// isValidMethod checks if the provided method is a valid HTTP method.
func isValidMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, m := range validMethods {
		if m == method {
			return true
		}
	}
	return false
}

// splitCurlCommand splits a curl command into tokens, handling quoted strings correctly.
func splitCurlCommand(curlCommand string) []string {
	var tokens []string
	var token strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(curlCommand); i++ {
		char := curlCommand[i]
		switch char {
		case ' ', '\t':
			if inQuotes {
				token.WriteByte(char)
			} else if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
		case '\'', '"':
			if inQuotes && char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else {
				token.WriteByte(char)
			}
		default:
			token.WriteByte(char)
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}

	return tokens
}
