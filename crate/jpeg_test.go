package crate_test

import (
	"time"

	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ImageFromPath(path string) *ImageMeta {
	node, err := NewPath(path)
	Ω(err).Should(BeNil())
	Ω(node).ShouldNot(BeNil())

	img, ok := ConvertImageMeta(node.(*FileMeta))
	Ω(ok).Should(BeTrue())
	Ω(img).ShouldNot(BeNil())

	return img
}

var _ = Describe("Jpeg", func() {

	const (
		coastPath   = "coast.jpg"
		draculaPath = "dracula.txt"
		ferryPath   = "ferry.jpg"
		wharfPath   = "wharf.png"
	)

	var (
		fixtures *Dir       // The root fixtures directory
		coast    *ImageMeta // The path to the coast fixture (JPG with GPS)
		ferry    *ImageMeta // The path to the ferry fixture (JPEG without GPS)
		wharf    *ImageMeta // The path to the wharf fixture (PNG)
	)

	BeforeEach(func() {

		// Locate the fixtures to test on
		if exists, _ := PathExists("./fixtures"); exists {
			fpath, _ := NewPath("./fixtures")
			fixtures = fpath.(*Dir)
		} else if exists, _ := PathExists("../fixtures"); exists {
			fpath, _ := NewPath("../fixtures/")
			fixtures = fpath.(*Dir)
		}

		Ω(fixtures).ShouldNot(BeNil())

		// Create the Images from the fixtures
		coast = ImageFromPath(fixtures.Join(coastPath))
		ferry = ImageFromPath(fixtures.Join(ferryPath))
		wharf = ImageFromPath(fixtures.Join(wharfPath))

	})

	It("should identify a JPEG as an JPEG", func() {
		Ω(coast.IsJPEG()).Should(BeTrue())
	})

	It("should not identify a PNG as an JPEG", func() {
		Ω(wharf.IsJPEG()).ShouldNot(BeTrue())
	})

	It("should not be able to extract an ExifHandler from a PNG", func() {
		exif, ok := wharf.GetExif()
		Ω(ok).Should(BeFalse())
		Ω(exif).Should(BeNil())
	})

	It("should be able to extract an ExifHandler from a JPEG", func() {
		var exif *ExifHandler
		var ok bool

		exif, ok = ferry.GetExif()
		Ω(ok).Should(BeTrue())
		Ω(exif).ShouldNot(BeNil())

		exif, ok = coast.GetExif()
		Ω(ok).Should(BeTrue())
		Ω(exif).ShouldNot(BeNil())
	})

	It("should be able to get tags", func() {
		exif, ok := coast.GetExif()
		Ω(ok).Should(BeTrue())
		Ω(exif).ShouldNot(BeNil())

		// Get the Make and Model
		Ω(exif.Get("Make")).Should(Equal("LGE"))
		Ω(exif.Get("Model")).Should(Equal("Nexus 5"))
	})

	XIt("should be able to extract the date taken", func() {
		// TODO: Figure out how to get this working in Travis
		exif, ok := coast.GetExif()
		Ω(ok).Should(BeTrue())
		Ω(exif).ShouldNot(BeNil())

		taken, _ := time.Parse("2006-01-02T15:04:05-07:00", "2015-01-05T17:57:30+00:00")
		original, err := exif.DateTaken()

		Ω(err).Should(BeNil())
		Ω(original.UTC()).Should(Equal(taken.UTC()))
	})

	It("should be able to extract the GPS coordinates", func() {
		exif, ok := coast.GetExif()
		Ω(ok).Should(BeTrue())
		Ω(exif).ShouldNot(BeNil())

		lat, lon, err := exif.Coordinates()
		Ω(err).Should(BeNil())
		Ω(lat).Should(Equal(31.510427472222222))
		Ω(lon).Should(Equal(-9.774266222222224))
	})

})
