package render

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
)

// OpenImage opens an image file from disk.
//
// Supported file types are: jpeg, gif, png, bmp.
func OpenImage(filename string) (image.Image, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var img image.Image

	switch filepath.Ext(filename) {
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(fh)
	case ".png":
		img, err = png.Decode(fh)
	case ".bmp":
		img, err = bmp.Decode(fh)
	case ".gif":
		img, err = gif.Decode(fh)
	default:
		return nil, errors.New("unsupported file type")
	}

	return img, err
}

// ImageToRGBA converts a Go image.Image into an image.RGBA.
func ImageToRGBA(input image.Image) *image.RGBA {
	var bounds = input.Bounds()
	var rgba = image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			color := input.At(x, y)
			rgba.Set(x, y, color)
		}
	}
	return rgba
}
