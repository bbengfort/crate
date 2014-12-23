// TODO: Find way to test stdout

package crate_test

import (
	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	"errors"
)

var _ = Describe("Console", func() {

	var (
		console *Console
	)

	BeforeEach(func() {
		console = new(Console)
		console.Init(true)
	})

	XContext("console in Debug Mode", func() {

		// How to evaluate output to stdout?

		It("should log to stdout", func() {
			console.Log("Testing %d, %d, %d: %s", 1, 2, 3, "mike check")

			// Output:
			// Testing 1, 2, 3: mike check
		})

		It("should print info to stdout", func() {
			console.Info("The %s in %s falls gently on the %s.", "rain", "Spain", "plane")

			// Output:
			// The rain in Spain falls gently on the plane.
		})

		It("should prefix errors with an error code", func() {
			err := errors.New("something bad happened")
			console.Err("errcode %d", err, 500)

			// Output:
			// ERROR (errcode 500): something bad happened
		})

	})

	XContext("console not in Debug Mode", func() {

		// How to evaluate output to stdout?

		BeforeEach(func() {
			console = new(Console)
			console.Init(false)
		})

		It("should supress info statements", func() {
			console.Info("The %s in %s falls gently on the %s.", "rain", "Spain", "plane")

			// Output:
		})

	})

})
