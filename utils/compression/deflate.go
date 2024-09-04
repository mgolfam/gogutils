package utils

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

// CompressData compresses data using deflate compression.
func Deflate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecompressData decompresses deflate compressed data.
func Inflate(compressedData []byte) ([]byte, error) {
	buf := bytes.NewReader(compressedData)
	reader := flate.NewReader(buf)
	defer reader.Close()

	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decompressedData, nil
}
