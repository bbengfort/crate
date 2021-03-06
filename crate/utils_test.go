package crate_test

import (
	"io/ioutil"
	"os"
	"time"

	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	It("should be able to hash to a sha1 b64 string", func() {
		text := []byte("The small brown fox jumped over the rabbit.")
		hash := "yPdVQEIMrUg13COQXCl69OCG3Sc="
		Ω(Hash(text)).Should(Equal(hash))
	})

	It("should be able to determine the mimetype of a text file", func() {
		// Create the test file
		testfile, _ := ioutil.TempFile("", "ginkgo-")
		testfile.Close()

		// Add data to the test file
		ioutil.WriteFile(testfile.Name(), []byte("hello world, just a normal text file"), 0644)

		Ω(MimeType(testfile.Name())).Should(Equal("text/plain"))

		// Clean up
		os.Remove(testfile.Name())
	})

	It("should be able to identify the hostname of the computer", func() {
		Ω(Hostname()).ShouldNot(BeZero())
	})

	Describe("Ftoa", func() {

		It("should convert a 0 float to a null string", func() {
			var num float64 // zero value float
			Ω(Ftoa(num)).Should(BeZero())
		})

		It("Should convert a float to a string value", func() {
			num := 3.144512341234412
			Ω(Ftoa(num)).Should(Equal("3.144512341234412"))
		})

	})

	It("should be able to convert a time to a JSON representation", func() {
		lt := "Mon Jan 2 15:04:05 -0700 MST 2006"
		dt, _ := time.Parse(lt, "Mon Jan 12 16:51:19 -0500 EST 2015")

		Ω(JSONStamp(dt)).Should(Equal("2015-01-12T21:51:19+00:00"))
	})

})
