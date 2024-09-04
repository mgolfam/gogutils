package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}
