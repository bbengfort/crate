// Handles writing information to stdout and stderr for crate

package crate

import (
	"fmt"
	"os"
)

//=============================================================================

// Handles the writing to the console
type Console struct {
	debug bool
}

// Writes a format string and arguments as a newline to stdout
func (console *Console) log(str string, args ...interface{}) {
	out := fmt.Sprintf(str, args...)
	fmt.Println(out)
}

// Logs to the default logger if debug is set to true
func (console *Console) info(str string, args ...interface{}) {
	if console.debug {
		console.log(str, args...)
	}
}

// Writes error messages out to the default logger
func (console *Console) err(str string, err error, args ...interface{}) {
	str = fmt.Sprintf("ERROR (%s): %s", str, err)
	fmt.Printf(str+"\n", args...)
}

// Fatal, writes error message out and quits
func (console *Console) fatal(str string, args ...interface{}) {
	str = fmt.Sprintf("FATAL: %s", str)
	console.log(str, args...)
	os.Exit(1)
}
