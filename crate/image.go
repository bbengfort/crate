// Meta data struct specifically for images

package crate

import (
	"encoding/json"
	"errors"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

//=============================================================================

type ImageMeta struct {
	FileMeta
	Width  int               // Width of the image
	Height int               // Height of the image
	Tags   map[string]string // Image tags from the Exif data
}

// Converts a FileMeta into an ImageMeta
func ConvertImageMeta(fm *FileMeta) (*ImageMeta, bool) {

	// If this isn't an image, then return nil
	if !fm.IsImage() {
		return nil, false
	}

	img := new(ImageMeta)
	img.FileMeta = *fm

	return img, true
}

// Popluates the fields on the ImageMeta
func (img *ImageMeta) Populate() {
	img.FileMeta.Populate() // Populate the FileMeta

	if width, height, err := img.Dimensions(); err == nil {
		img.Width = width
		img.Height = height
	}

	exif := img.GetExif()
	if exif != nil {
		img.Tags = exif.tags
	}
}

// Returns the width, hight of the image
func (img *ImageMeta) Dimensions() (int, int, error) {

	if file, err := os.Open(img.Path); err == nil {
		defer file.Close()

		if config, _, err := image.DecodeConfig(file); err == nil {
			return config.Width, config.Height, nil
		}

		return 0, 0, errors.New("Could not decode Image dimensions")
	}

	return 0, 0, errors.New("Could not open Image for reading")
}

// Returns the byte serialization of the file meta for storage
func (img *ImageMeta) Byte() []byte {
	data, err := json.Marshal(img)
	if err != nil {
		return nil
	}

	return data
}

// Prints out the info as a JSON indented pretty string
func (img *ImageMeta) Info() string {
	if !img.populated {
		img.Populate()
	}

	info, err := json.MarshalIndent(img, "", "  ")
	if err != nil {
		return ""
	}

	return string(info)
}
