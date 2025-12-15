package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mgolfam/gogutils/glog"

	"github.com/golang-cz/textcase"
)

func StringExists(target string, slice []string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}

func Substring(text string, maxLength int) string {
	// If the text length is less than or equal to maxLength, return the original text
	if len(text) <= maxLength {
		return text
	}

	// If the text length is greater than maxLength, return a substring of length maxLength
	return text[:maxLength]
}

func Atoi(text string) int {
	i, err := strconv.Atoi(text)
	if err != nil {
		// Handle error if conversion fails
		glog.LogL(glog.DEBUG, "Error:", err)
		return -100
	}

	return i
}

func AtoiPtr(text string) *int {
	i, err := strconv.Atoi(text)
	if err != nil {
		// Handle error if conversion fails
		glog.LogL(glog.DEBUG, "Error:", err)
		return nil
	}

	return &i
}

// toPascalCase converts a string to PascalCase.
func PascalCase(s string) string {
	firstRune, size := utf8.DecodeRuneInString(s)
	pascal := string(unicode.ToTitle(firstRune)) + s[size:]
	return pascal
}

// toCamelCase converts a string to CamelCase.
func CamelCase(s string) string {
	return textcase.CamelCase(s)
}

// toSnakeCase converts a string to SnakeCase.
func SnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	snake := string(result)
	return snake
}

func KebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '-')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	snake := string(result)
	return snake
}

func CheckStringCategory(s string) string {
	hasDigit := false
	hasLetter := false

	for _, char := range s {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsLetter(char) {
			hasLetter = true
		} else {
			// If there are any non-digit and non-letter characters, return "invalid"
			return "invalid"
		}
	}

	if hasDigit && !hasLetter {
		return "all_digits"
	} else if hasLetter && !hasDigit {
		return "all_alphabets"
	} else if hasDigit && hasLetter {
		return "mixed"
	}

	return "invalid" // Shouldn't reach here, just a safety measure
}

// StringPtr returns a pointer to the given string.
// Useful for building structs with optional string fields.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(b bool) *bool {
	return &b
}

// Slugify converts a string to a URL-friendly slug: lowercase, spaces to '-',
// and removes non-alphanumeric characters (except '-').
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// Replace whitespace with dash
	reSpace := regexp.MustCompile(`\s+`)
	s = reSpace.ReplaceAllString(s, "-")
	// Remove invalid chars
	reInvalid := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reInvalid.ReplaceAllString(s, "")
	// Collapse multiple dashes
	reDash := regexp.MustCompile(`-+`)
	s = reDash.ReplaceAllString(s, "-")
	return s
}

// NormalizeSpaces trims leading/trailing spaces and collapses internal
// whitespace runs to a single space.
func NormalizeSpaces(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}

// NormalizeNewlines converts DOS/old Mac newlines to '\n'.
func NormalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// StripHTML removes HTML tags using a simple regex.
// This is not a full HTML sanitizer, but fine for logs and basic cleanup.
func StripHTML(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(s, "")
}

// SecureRandomString generates a cryptographically secure random string of the given length,
// using the provided alphabet. If alphabet is empty, a default URL-safe set is used.
func SecureRandomString(length int, alphabet string) (string, error) {
	if length <= 0 {
		return "", nil
	}
	if alphabet == "" {
		alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	var sb strings.Builder
	l := len(alphabet)
	for _, b := range bytes {
		sb.WriteByte(alphabet[int(b)%l])
	}
	return sb.String(), nil
}

// MustSecureRandomString is like SecureRandomString but panics on error.
func MustSecureRandomString(length int, alphabet string) string {
	s, err := SecureRandomString(length, alphabet)
	if err != nil {
		panic(fmt.Sprintf("SecureRandomString failed: %v", err))
	}
	return s
}
