package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

// CompressData compresses data using gzip compression.
func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecompressData decompresses gzip compressed data.
func Gunzip(compressedData []byte) ([]byte, error) {
	buf := bytes.NewReader(compressedData)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	decompressedData, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}
