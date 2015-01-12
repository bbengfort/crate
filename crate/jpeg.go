// Handles exif data from JPEG files

package crate

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
)

//=============================================================================

// Checks if an Image is a JPEG file
func (img *ImageMeta) IsJPEG() bool {
	if img.IsImage() {
		return img.MimeType == "image/jpeg"
	}

	return false
}

//=============================================================================

// Implements the Walker interface to retrieve all tags
type ExifHandler struct {
	exif *exif.Exif        // The Exif struct being walked
	tags map[string]string // Tags
}

// Implements the Walk function to be a Walker
func (ew *ExifHandler) Walk(name exif.FieldName, tag *tiff.Tag) error {
	ew.tags[string(name)] = tag.String()
	return nil
}

// Retrieves a tag from the exif data in the handler
func (ew *ExifHandler) Get(tag exif.FieldName) string {
	val, _ := ew.exif.Get(tag)
	if val != nil {
		return val.String()
	}

	return ""
}

// Helper function for the date taken
func (ew *ExifHandler) DateTaken() (time.Time, error) {
	return ew.exif.DateTime()
}

// Helper function for the GPS Coordinates
func (ew *ExifHandler) Coordinates() (float64, float64, error) {
	return ew.exif.LatLong()
}

//=============================================================================

// Get the EXIF Data from the JPEG
func (img *ImageMeta) GetExif() *ExifHandler {

	// Ensure that this is a JPEG
	if !img.IsJPEG() {
		return nil
	}

	if f, err := os.Open(img.Path); err == nil {
		defer f.Close()
		exif.RegisterParsers(mknote.All...)

		walker := new(ExifHandler)
		walker.tags = make(map[string]string)

		if x, err := exif.Decode(f); err == nil {
			walker.exif = x
			x.Walk(walker)

			return walker
		}
	}

	return nil
}
