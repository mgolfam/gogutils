package crypt

import (
	"strconv"

	"github.com/mgolfam/gogutils/glog"
)

// EncodeBase36 converts a decimal number to its base36 string representation.
func EncodeBase36(number int64) string {
	return strconv.FormatInt(number, 36)
}

// DecondeBase36 converts a base36-encoded string back to its decimal value.
// NOTE: The function name is intentionally kept (with the typo) for backward compatibility.
//       Prefer using DecodeBase36 in new code.
func DecondeBase36(text string) int64 {
	value, err := strconv.ParseInt(text, 36, 64)
	if err != nil {
		glog.LogL(glog.DEBUG, "base36 conversion failed:", err)
		return -1
	}
	return value
}

// DecodeBase36 is the correctly-spelled helper that delegates to DecondeBase36.
// Kept separate so that existing callers of DecondeBase36 keep working.
func DecodeBase36(text string) int64 {
	return DecondeBase36(text)
}
