// Handles exif data from JPEG files

package crate

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
)

const (
	ExifTimeLayout = "2006:01:02 15:04:05" // Time layout for EXIF data type
)

var (
	GPSTimePattern = regexp.MustCompile("\"(\\d+)/\\d+\"")
)

//=============================================================================

// Checks if an Image is a JPEG file
func (img *ImageMeta) IsJPEG() bool {
	if img.IsImage() {
		return img.MimeType == "image/jpeg"
	}

	return false
}

// Get the EXIF Data from the JPEG
func (img *ImageMeta) GetExif() (*ExifHandler, bool) {

	// Ensure that this is a JPEG
	if !img.IsJPEG() {
		return nil, false
	}

	if f, err := os.Open(img.Path); err == nil {
		defer f.Close()
		exif.RegisterParsers(mknote.All...)

		walker := new(ExifHandler)
		walker.tags = make(map[string]string)

		if x, err := exif.Decode(f); err == nil {
			walker.exif = x
			x.Walk(walker)

			return walker, true
		}
	}

	return nil, false
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
		return strings.Trim(val.String(), "\"")
	}

	return ""
}

// Helper function to fetch the GPS date and time
func (ew *ExifHandler) GPSDateTime() (time.Time, error) {
	var dt time.Time

	// Get the GPS Date and Time Stamps
	gpsDateStamp := ew.Get(exif.GPSDateStamp)
	gpsTimeStamp := ew.Get(exif.GPSTimeStamp)

	if gpsTimeStamp != "" {

		// Handle GPS timestamp parsing
		rats := GPSTimePattern.FindAllStringSubmatch(gpsTimeStamp, 3)
		nums := make([]string, len(rats), cap(rats))
		for idx, item := range rats {
			nums[idx] = item[1]
		}

		if len(nums) != 3 {
			return dt, errors.New("Could not parse GPSTimeStamp")
		}

		gpsTimeStr := strings.Join(nums, ":")

		// Handle the GPSDateStamp (if not available use date from exif)
		if gpsDateStamp == "" {

			exifdt, err := ew.exif.DateTime()
			if err != nil {
				return dt, err
			}

			gpsDateStamp = exifdt.Format("2006:01:02")
		}

		gpsDateStr := gpsDateStamp + " " + gpsTimeStr
		return time.ParseInLocation(ExifTimeLayout, gpsDateStr, time.UTC)

	} else {
		return dt, errors.New("No GPSDateTime available in EXIF")
	}
}

// Helper function for the date taken - overloads the exif library DateTime
func (ew *ExifHandler) DateTaken() (time.Time, error) {
	var dt time.Time

	if gpsDate, err := ew.GPSDateTime(); err == nil {
		return gpsDate, nil
	}

	// Get either the DateTimeOriginal or the DateTime
	tag, err := ew.exif.Get(exif.DateTimeOriginal)
	if err != nil {
		tag, err = ew.exif.Get(exif.DateTime)
		if err != nil {
			return dt, err
		}
	}

	if tag.Format() != tiff.StringVal {
		return dt, errors.New("DateTime[Original] not in string format")
	}

	dateStr := strings.TrimRight(string(tag.Val), "\x00")
	return time.Parse(ExifTimeLayout, dateStr)
}

// Helper function for the GPS Coordinates
func (ew *ExifHandler) Coordinates() (float64, float64, error) {
	return ew.exif.LatLong()
}

//=============================================================================
