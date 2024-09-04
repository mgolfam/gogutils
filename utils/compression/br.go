//go:build linux
// +build linux

package utils

import (
	"bytes"
	"io"

	"github.com/google/brotli/go/cbrotli"
)

func BrCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := cbrotli.NewWriter(&buf, cbrotli.WriterOptions{})
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BrDecompress(data []byte) ([]byte, error) {
	reader := cbrotli.NewReader(bytes.NewReader(data))
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return decompressed, nil
}
