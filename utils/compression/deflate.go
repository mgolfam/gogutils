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

	if _, err = writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	// Close flushes any buffered compressed data into buf; it must run
	// before buf.Bytes() is read, so it cannot be deferred.
	if err := writer.Close(); err != nil {
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
