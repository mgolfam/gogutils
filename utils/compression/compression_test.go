package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestGzipRoundTrip(t *testing.T) {
	original := []byte(strings.Repeat("gogutils gzip test payload ", 100))

	compressed, err := Gzip(original)
	if err != nil {
		t.Fatalf("Gzip returned error: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("Gzip returned empty output")
	}

	got, err := Gunzip(compressed)
	if err != nil {
		t.Fatalf("Gunzip returned error: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("gzip round trip mismatch: got %d bytes, want %d", len(got), len(original))
	}
}

func TestDeflateRoundTrip(t *testing.T) {
	original := []byte(strings.Repeat("gogutils deflate test payload ", 100))

	compressed, err := Deflate(original)
	if err != nil {
		t.Fatalf("Deflate returned error: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("Deflate returned empty output")
	}

	got, err := Inflate(compressed)
	if err != nil {
		t.Fatalf("Inflate returned error: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("deflate round trip mismatch: got %d bytes, want %d", len(got), len(original))
	}
}

func TestGzipCompressesRepetitiveData(t *testing.T) {
	original := bytes.Repeat([]byte("A"), 10000)
	compressed, err := Gzip(original)
	if err != nil {
		t.Fatalf("Gzip returned error: %v", err)
	}
	if len(compressed) >= len(original) {
		t.Errorf("Gzip did not shrink highly repetitive data: %d >= %d", len(compressed), len(original))
	}
}

func TestGunzipInvalidInput(t *testing.T) {
	if _, err := Gunzip([]byte("not gzip data")); err == nil {
		t.Error("Gunzip expected an error on invalid input, got nil")
	}
}
