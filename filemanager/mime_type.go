package filemanager

import (
	"fmt"
	"strings"
)

func FileExtensionFromBytes(data []byte) (string, error) {
	if len(data) < 4 {
		return "", fmt.Errorf("insufficient data to determine file extension")
	}

	// Check the magic number or signature for common file types
	switch {
	case strings.HasPrefix(string(data[:2]), "BM"):
		return "bmp", nil
	case string(data[:3]) == "\xFF\xD8\xFF":
		return "jpg", nil
	case string(data[:4]) == "GIF8":
		return "gif", nil
	case string(data[:4]) == "%PDF":
		return "pdf", nil
	case string(data[:8]) == "\x89PNG\r\n\x1A\n":
		return "png", nil
	default:
		return "", fmt.Errorf("unknown file type")
	}
}
