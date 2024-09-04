package utils

import (
	"strconv"
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
