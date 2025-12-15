package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// HMACSHA256 returns a hex-encoded HMAC-SHA256 of data using the given key.
func HMACSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// VerifyHMACSHA256 compares a message MAC with a computed MAC in constant time.
func VerifyHMACSHA256(key, data, expectedMAC []byte) bool {
	mac := HMACSHA256(key, data)
	return hmac.Equal(mac, expectedMAC)
}

// Base64URLEncode encodes data using URL-safe base64 without padding.
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URLDecode decodes URL-safe base64 without padding.
func Base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}


