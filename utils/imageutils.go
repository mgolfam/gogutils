package utils

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

// LoadImage loads an image from a file path and returns it with its format.
func LoadImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

// SaveJPEG saves an image to the given path as JPEG with the specified quality (1-100).
func SaveJPEG(path string, img image.Image, quality int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if quality <= 0 || quality > 100 {
		quality = 80
	}
	opts := &jpeg.Options{Quality: quality}
	return jpeg.Encode(f, img, opts)
}

// SavePNG saves an image to the given path as PNG.
func SavePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// SaveGIF saves an image to the given path as GIF.
func SaveGIF(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.Encode(f, img, nil)
}

// ResizeImage returns a new image with the given width and height using ApproxBiLinear scaling.
func ResizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// ThumbnailImage keeps aspect ratio and fits the image into maxWidth x maxHeight.
func ThumbnailImage(src image.Image, maxWidth, maxHeight int) image.Image {
	b := src.Bounds()
	w := b.Dx()
	h := b.Dy()
	if w <= maxWidth && h <= maxHeight {
		return src
	}
	ratioW := float64(maxWidth) / float64(w)
	ratioH := float64(maxHeight) / float64(h)
	ratio := ratioW
	if ratioH < ratio {
		ratio = ratioH
	}
	newW := int(float64(w) * ratio)
	newH := int(float64(h) * ratio)
	return ResizeImage(src, newW, newH)
}

// ConvertImageFormat loads an image from inPath and saves it to outPath in the given format ("jpeg", "png", "gif").
func ConvertImageFormat(inPath, outPath, format string) error {
	img, _, err := LoadImage(inPath)
	if err != nil {
		return err
	}
	switch format {
	case "jpeg", "jpg":
		return SaveJPEG(outPath, img, 80)
	case "png":
		return SavePNG(outPath, img)
	case "gif":
		return SaveGIF(outPath, img)
	default:
		// default to PNG
		return SavePNG(outPath, img)
	}
}


