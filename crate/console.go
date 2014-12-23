// Handles writing information to stdout and stderr for crate

package crate

import (
	"fmt"
	"os"
)

var console *Console

//=============================================================================

// Handles the writing to the console
type Console struct {
	debug bool
}

// Initializes a new Console object
func (console *Console) Init(debug bool) {
	console.debug = debug
}

// Writes a format string and arguments as a newline to stdout
func (console *Console) Log(str string, args ...interface{}) {
	out := fmt.Sprintf(str, args...)
	fmt.Println(out)
}

// Logs to the default logger if debug is set to true
func (console *Console) Info(str string, args ...interface{}) {
	if console.debug {
		console.Log(str, args...)
	}
}

// Writes error messages out to the default logger
func (console *Console) Err(str string, err error, args ...interface{}) {
	str = fmt.Sprintf("ERROR (%s): %s", str, err)
	fmt.Printf(str+"\n", args...)
}

// Fatal, writes error message out and quits
func (console *Console) Fatal(str string, args ...interface{}) {
	str = fmt.Sprintf("FATAL: %s", str)
	console.Log(str, args...)
	os.Exit(1)
}
