package utils

import (
	"errors"
	"strings"
)

const (
	EncodingGzip    = "gzip"
	EncodingDeflate = "deflate"
)

// CompressWithEncoding compresses data using the given content-encoding name.
// Supported: "gzip", "deflate".
func CompressWithEncoding(data []byte, encoding string) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case EncodingGzip:
		return Gzip(data)
	case EncodingDeflate:
		return Deflate(data)
	default:
		return nil, errors.New("unsupported encoding: " + encoding)
	}
}

// DecompressWithEncoding decompresses data using the given content-encoding name.
// Supported: "gzip", "deflate".
func DecompressWithEncoding(data []byte, encoding string) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case EncodingGzip:
		return Gunzip(data)
	case EncodingDeflate:
		return Inflate(data)
	default:
		return nil, errors.New("unsupported encoding: " + encoding)
	}
}


