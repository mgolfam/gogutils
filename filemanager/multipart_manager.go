package filemanager

import (
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/mgolfam/gogutils/glog"
	"github.com/mgolfam/gogutils/utils"
)

// FileInfo represents information about a file
type FileInfo struct {
	Path         string
	Name         string
	OriginalName string
	Extension    string
	Size         int64
}

// GetFileInfo retrieves file information from a multipart.FileHeader
func GetFileInfo(fileHeader *multipart.FileHeader) (*FileInfo, error) {
	// Open the file from the file header
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file information
	fileInfo := &FileInfo{
		OriginalName: fileHeader.Filename,
		Extension:    filepath.Ext(fileHeader.Filename),
		Size:         fileHeader.Size,
	}

	return fileInfo, nil
}

func SaveMulitpartFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (*FileInfo, error) {
	fileInfo, err := GetFileInfo(fileHeader)
	if err != nil {
		glog.LogL(glog.DEBUG, err)
		return nil, err
	}

	MkDir(path)

	fileInfo.Path = path
	fileInfo.Name = utils.CleanUUID() + fileInfo.Extension

	fullPath := filepath.Join(path, fileInfo.Name)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	err = WriteBinaryToFile(buf.Bytes(), fullPath)
	if err != nil {
		return nil, err
	}

	return fileInfo, err
}
