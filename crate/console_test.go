// Testing for the helper console

package crate

import (
	"errors"
)

func ExampleConsoleLog() {
	console := &Console{true}
	console.Log("Testing %d, %d, %d: %s", 1, 2, 3, "mike check")

	// Output:
	// Testing 1, 2, 3: mike check
}

func ExampleInfoDebug() {
	console := &Console{true}
	console.Info("The %s in %s falls gently on the %s.", "rain", "Spain", "plane")

	// Output:
	// The rain in Spain falls gently on the plane.
}

func ExampleInfoDebugSupress() {
	console := &Console{false}
	console.Info("The %s in %s falls gently on the %s.", "rain", "Spain", "plane")

	// Output:
}

func ExampleError() {
	console := &Console{true}
	err := errors.New("something bad happened")
	console.Err("errcode %d", err, 500)

	// Output:
	// ERROR (errcode 500): something bad happened
}
