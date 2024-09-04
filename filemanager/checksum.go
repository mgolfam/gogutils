package filemanager

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
)

func CalculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func CalculateBinaryChecksum(input []byte) string {
	hash := sha256.New()
	hash.Write(input)
	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum
}

func CalculateBase64Checksum(base64String string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}

	return CalculateBinaryChecksum(data), nil
}
