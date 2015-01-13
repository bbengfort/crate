package crate_test

import (
	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image", func() {

	const (
		coastPath   = "coast.jpg"
		draculaPath = "dracula.txt"
		ferryPath   = "ferry.jpg"
		wharfPath   = "wharf.png"
	)

	var (
		fixtures *Dir // The root fixtures directory
		coast    Path // The path to the coast fixture (JPG with GPS)
		dracula  Path // The path to the dracula fixture (TEXT)
		ferry    Path // The path to the ferry fixture (JPEG without GPS)
		wharf    Path // The path to the wharf fixture (PNG)
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

		// Add the Path fixtures for use in testing
		coast, _ = NewPath(fixtures.Join(coastPath))
		dracula, _ = NewPath(fixtures.Join(draculaPath))
		ferry, _ = NewPath(fixtures.Join(ferryPath))
		wharf, _ = NewPath(fixtures.Join(wharfPath))

	})

	It("should identify a JPEG as an image", func() {
		fm := coast.(*FileMeta)
		Ω(fm.IsImage()).Should(BeTrue())
	})

	It("should identify a PNG as an image", func() {
		fm := wharf.(*FileMeta)
		Ω(fm.IsImage()).Should(BeTrue())
	})

	It("should not identify TEXT as an image", func() {
		fm := dracula.(*FileMeta)
		Ω(fm.IsImage()).Should(BeFalse())
	})

	It("should convert an image FileMeta to an ImageMeta", func() {
		fm := coast.(*FileMeta)
		var img *ImageMeta

		converted, ok := ConvertImageMeta(fm)

		Ω(ok).Should(BeTrue())
		Ω(converted).Should(BeAssignableToTypeOf(img))
	})

	It("should not convert a text FileMeta to an ImageMeta", func() {
		fm := dracula.(*FileMeta)

		converted, ok := ConvertImageMeta(fm)

		Ω(ok).ShouldNot(BeTrue())
		Ω(converted).Should(BeNil())
	})

	It("should be able to compute the dimensions of a JPEG", func() {
		img, ok := ConvertImageMeta(ferry.(*FileMeta))
		Ω(ok).Should(BeTrue())
		Ω(img).ShouldNot(BeNil())

		width, height, err := img.Dimensions()
		Ω(err).Should(BeNil())
		Ω(width).Should(Equal(4288))
		Ω(height).Should(Equal(3216))
	})

	It("should be able to compute the dimensions of a PNG", func() {
		img, ok := ConvertImageMeta(wharf.(*FileMeta))
		Ω(ok).Should(BeTrue())
		Ω(img).ShouldNot(BeNil())

		width, height, err := img.Dimensions()
		Ω(err).Should(BeNil())
		Ω(width).Should(Equal(2737))
		Ω(height).Should(Equal(1354))
	})

	It("should begin unpopulated", func() {
		img, ok := ConvertImageMeta(ferry.(*FileMeta))
		Ω(ok).Should(BeTrue())
		Ω(img).ShouldNot(BeNil())

		// Test FileMeta specific fields for unpopulation
		// Note that MimeType will be populated as part of the conversion
		Ω(img.Name).Should(BeZero())
		Ω(img.Size).Should(BeZero())
		Ω(img.Modified).Should(BeZero())
		Ω(img.Signature).Should(BeZero())
		Ω(img.Host).Should(BeZero())
		Ω(img.Author).Should(BeZero())

		// Test Image specific fields for unpopulation
		Ω(img.Width).Should(BeZero())
		Ω(img.Height).Should(BeZero())
		Ω(img.Tags).Should(BeZero())
	})

	It("should be populated", func() {
		img, ok := ConvertImageMeta(ferry.(*FileMeta))
		Ω(ok).Should(BeTrue())
		Ω(img).ShouldNot(BeNil())

		// Populate!
		img.Populate()

		// Test FileMeta specific fields for population
		Ω(img.MimeType).ShouldNot(BeZero(), "MimeType must be populated")
		Ω(img.Name).ShouldNot(BeZero(), "Name must be populated")
		Ω(img.Size).ShouldNot(BeZero(), "Size must be populated")
		Ω(img.Modified).ShouldNot(BeZero(), "Modified must be populated")
		Ω(img.Signature).ShouldNot(BeZero(), "Signature must be populated")
		Ω(img.Host).ShouldNot(BeZero(), "Host must be populated")
		Ω(img.Author).ShouldNot(BeZero(), "Author must be populated")

		// Test Image specific fields for population
		Ω(img.Width).ShouldNot(BeZero(), "Width must be populated")
		Ω(img.Height).ShouldNot(BeZero(), "Height must be populated")
		Ω(img.Tags).ShouldNot(BeZero(), "Tags must be populated")
	})

})
