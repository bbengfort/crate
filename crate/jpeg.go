// Handles exif data from JPEG files

package crate

import (
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// Checks if a FileMeta is an image
func (fm *FileMeta) IsImage() bool {
	if fm.MimeType == "" {
		fm.MimeType, _ = MimeType(fm.Path)
	}

	return strings.HasPrefix(fm.MimeType, "image/")
}

// Checks if an Image is a JPEG file
func (fm *FileMeta) IsJPEG() bool {
	if fm.IsImage() {
		return fm.MimeType == "image/jpeg"
	}

	return false
}

// Quick helper function
func GetExifString(x *exif.Exif, tag exif.FieldName) string {
	val, _ := x.Get(tag)
	if val != nil {
		return val.String()
	}

	return ""
}

// Get the EXIF Data from the JPEG
func (fm *FileMeta) GetExif() map[string]string {
	data := make(map[string]string)

	if !fm.IsJPEG() {
		return data
	}

	if f, err := os.Open(fm.Path); err == nil {
		defer f.Close()
		exif.RegisterParsers(mknote.All...)

		if x, err := exif.Decode(f); err == nil {
			data["CameraModel"] = GetExifString(x, exif.Model)
			dt, _ := x.DateTime()
			data["DateTaken"] = dt.String()
			lat, lon, _ := x.LatLong()
			data["Latitude"] = Ftoa(lat)
			data["Longitude"] = Ftoa(lon)
		}
	}

	return data
}
