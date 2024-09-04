package filemanager

import (
	"encoding/base64"
	"io"
	"os"
)

func File2Base64(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	size := fi.Size()
	fileContent := make([]byte, size)

	_, err = io.ReadFull(file, fileContent)
	if err != nil {
		return "", err
	}

	base64Content := base64.StdEncoding.EncodeToString(fileContent)
	return base64Content, nil
}

func Base642File(base64Content string, filePath string) error {
	decodedContent, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(decodedContent)
	if err != nil {
		return err
	}

	return nil
}
