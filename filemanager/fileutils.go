package filemanager

import "path/filepath"

func FileExtension(fileName string) string {
	// Get the file extension
	extension := filepath.Ext(fileName)
	// Remove the dot (.) from the extension
	extension = extension[1:]

	return extension
}
