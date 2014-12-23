package crate_test

import (
	"io/ioutil"
	"os"

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

})
