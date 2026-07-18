package filemanager

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

// sha256("hello")
const helloSha256 = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

func TestCalculateBinaryChecksum(t *testing.T) {
	if got := CalculateBinaryChecksum([]byte("hello")); got != helloSha256 {
		t.Errorf("CalculateBinaryChecksum = %q, want %q", got, helloSha256)
	}
}

func TestCalculateFileChecksum(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := CalculateFileChecksum(path)
	if err != nil {
		t.Fatalf("CalculateFileChecksum returned error: %v", err)
	}
	if got != helloSha256 {
		t.Errorf("CalculateFileChecksum = %q, want %q", got, helloSha256)
	}

	if _, err := CalculateFileChecksum(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("CalculateFileChecksum(missing) expected an error, got nil")
	}
}

func TestCalculateBase64Checksum(t *testing.T) {
	b64 := base64.StdEncoding.EncodeToString([]byte("hello"))
	got, err := CalculateBase64Checksum(b64)
	if err != nil {
		t.Fatalf("CalculateBase64Checksum returned error: %v", err)
	}
	if got != helloSha256 {
		t.Errorf("CalculateBase64Checksum = %q, want %q", got, helloSha256)
	}

	if _, err := CalculateBase64Checksum("!!!not base64!!!"); err == nil {
		t.Error("CalculateBase64Checksum(invalid) expected an error, got nil")
	}
}

func TestFile2Base64RoundTrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.bin")
	payload := []byte{0x00, 0x01, 0x02, 0xff, 0xfe}
	if err := os.WriteFile(src, payload, 0o644); err != nil {
		t.Fatal(err)
	}

	encoded, err := File2Base64(src)
	if err != nil {
		t.Fatalf("File2Base64 returned error: %v", err)
	}
	if encoded != base64.StdEncoding.EncodeToString(payload) {
		t.Errorf("File2Base64 = %q, unexpected encoding", encoded)
	}

	dst := filepath.Join(dir, "out.bin")
	if err := Base642File(encoded, dst); err != nil {
		t.Fatalf("Base642File returned error: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Errorf("Base642File wrote %v, want %v", got, payload)
	}
}

func TestFileExtension(t *testing.T) {
	cases := map[string]string{
		"photo.png":       "png",
		"archive.tar.gz":  "gz",
		"/a/b/report.pdf": "pdf",
	}
	for in, want := range cases {
		if got := FileExtension(in); got != want {
			t.Errorf("FileExtension(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFileExtensionFromBytes(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want string
	}{
		{"bmp", []byte("BM\x00\x00"), "bmp"},
		{"jpg", []byte("\xFF\xD8\xFF\x00"), "jpg"},
		{"gif", []byte("GIF8"), "gif"},
		{"pdf", []byte("%PDF"), "pdf"},
		{"png", []byte("\x89PNG\r\n\x1A\n"), "png"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := FileExtensionFromBytes(c.data)
			if err != nil {
				t.Fatalf("FileExtensionFromBytes returned error: %v", err)
			}
			if got != c.want {
				t.Errorf("FileExtensionFromBytes = %q, want %q", got, c.want)
			}
		})
	}

	if _, err := FileExtensionFromBytes([]byte{0x01}); err == nil {
		t.Error("FileExtensionFromBytes(too short) expected an error, got nil")
	}
	if _, err := FileExtensionFromBytes([]byte("random!")); err == nil {
		t.Error("FileExtensionFromBytes(unknown) expected an error, got nil")
	}
}

func TestReadWriteFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.txt")

	if err := WriteFile("first", path); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if got != "first" {
		t.Errorf("ReadFile = %q, want %q", got, "first")
	}

	// WriteFile overwrites.
	if err := WriteFile("second", path); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadFile(path); got != "second" {
		t.Errorf("after overwrite ReadFile = %q, want %q", got, "second")
	}

	// AppendFile appends.
	if err := AppendFile("-third", path); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadFile(path); got != "second-third" {
		t.Errorf("after append ReadFile = %q, want %q", got, "second-third")
	}
}

func TestReadWriteFileBytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "blob.bin")
	data := []byte{0xde, 0xad, 0xbe, 0xef}

	if err := WriteFileBytes(path, data, 0o644); err != nil {
		t.Fatalf("WriteFileBytes returned error: %v", err)
	}
	got, err := ReadFileBytes(path)
	if err != nil {
		t.Fatalf("ReadFileBytes returned error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("ReadFileBytes = %v, want %v", got, data)
	}
}

func TestMkDirAndFileDirExist(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deep")
	if FileDirExist(dir) {
		t.Fatal("FileDirExist reported a non-existent path as existing")
	}
	if _, err := MkDir(dir); err != nil {
		t.Fatalf("MkDir returned error: %v", err)
	}
	if !FileDirExist(dir) {
		t.Error("FileDirExist reported a created directory as missing")
	}
}

func TestDeleteFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "temp.txt")
	if err := WriteFile("bye", path); err != nil {
		t.Fatal(err)
	}
	if err := DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile returned error: %v", err)
	}
	if FileDirExist(path) {
		t.Error("file still exists after DeleteFile")
	}
}

func TestFileWriter(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stream.log")
	w, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("NewFileWriter returned error: %v", err)
	}
	if _, err := w.Write([]byte("line1\n")); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("line2\n")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	got, _ := ReadFile(path)
	if got != "line1\nline2\n" {
		t.Errorf("FileWriter wrote %q, want %q", got, "line1\nline2\n")
	}
}
