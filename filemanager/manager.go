package filemanager

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/mgolfam/gogutils/glog"
)

// FileWriter is a custom type that implements io.Writer
type FileWriter struct {
	file *os.File
}

func ReadFileBytes(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func WriteFileBytes(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func MkDir(path string) (bool, error) {
	if FileDirExist(path) {
		return true, nil
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		glog.LogL(glog.DEBUG, err)
		return false, err
	}
	return true, err
}

func FileDirExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// path/to/whatever exists

		return true
	}

	return false
}

// ReadFile reads the contents of a file and returns it as a string.
func ReadFile(inFilePath string) (string, error) {
	data, err := ReadFileBytes(inFilePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// AppendFile appends text to an existing file. If the file does not exist, it creates one.
func AppendFile(text string, outFilePath string) error {
	file, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile writes text to a file, overwriting the file if it already exists.
func WriteFile(text string, outFilePath string) error {
	err := WriteFileBytes(outFilePath, []byte(text), 0644)
	if err != nil {
		return err
	}
	return nil
}

// WriteFile writes text to a file, overwriting the file if it already exists.
func WriteFileBin(data []byte, outFilePath string) error {
	err := WriteFileBytes(outFilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// WriteFile writes text to a file, overwriting the file if it already exists.
func Write(file io.Reader, outFilePath string) error {
	out, err := os.Create(outFilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, file)
	defer out.Close()
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile deletes the specified file.
func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

// NewFileWriter creates a new FileWriter
func NewFileWriter(filePath string) (*FileWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: file}, nil
}

// Write implements the io.Writer interface
func (fw *FileWriter) Write(p []byte) (n int, err error) {
	return fw.file.Write(p)
}

// Close closes the underlying file
func (fw *FileWriter) Close() error {
	return fw.file.Close()
}

// WriteTextToFile writes text to a file, overwriting the file if it already exists.
func WriteTextToFile(text string, outFilePath string) error {
	writer, err := NewFileWriter(outFilePath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.WriteString(writer, text)
	return err
}

// WriteBinaryToFile writes binary data to a file, overwriting the file if it already exists.
func WriteBinaryToFile(data []byte, outFilePath string) error {
	writer, err := NewFileWriter(outFilePath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.Write(data)
	return err
}

// JSON file helpers

// ReadJSONFile reads a JSON file into v.
func ReadJSONFile(path string, v interface{}) error {
	data, err := ReadFileBytes(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSONFile writes v as indented JSON to path.
func WriteJSONFile(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return WriteFileBytes(path, data, 0644)
}

// Directory and path utilities

// EnsureDir ensures that the directory at path exists.
func EnsureDir(path string) error {
	_, err := MkDir(path)
	return err
}

// EnsureFileWithDefault ensures that a file exists; if not, writes defaultContent.
func EnsureFileWithDefault(path string, defaultContent []byte) error {
	if FileDirExist(path) {
		return nil
	}
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	return WriteFileBytes(path, defaultContent, 0644)
}

// TempFile creates a temp file in dir (or os.TempDir if empty) with the given pattern.
func TempFile(dir, pattern string) (*os.File, error) {
	if dir == "" {
		return os.CreateTemp("", pattern)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, pattern)
}

// Simple file locking using a lock file.

// Lock represents a file-based lock.
type Lock struct {
	Path string
}

// AcquireLock tries to create a lock file; if it already exists, it waits up to timeout.
// Lock files are simple zero-byte files at lockPath.
func AcquireLock(lockPath string, timeout time.Duration) (*Lock, error) {
	start := time.Now()
	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL, 0644)
		if err == nil {
			f.Close()
			return &Lock{Path: lockPath}, nil
		}

		if timeout > 0 && time.Since(start) > timeout {
			return nil, err
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Release removes the underlying lock file.
func (l *Lock) Release() error {
	if l == nil || l.Path == "" {
		return nil
	}
	return os.Remove(l.Path)
}
