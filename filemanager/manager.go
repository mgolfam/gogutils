package filemanager

import (
	"io"
	"io/fs"
	"os"

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
