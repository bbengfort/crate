// Utility functions for Crate

package crate

import (
	"os"
	"os/signal"
	"syscall"
)

//=============================================================================

// Watch for CTRL+C and terminate the server
func signalHandler() {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)
	<-sigchan

	// On shutdown, stop the Crate watcher
	console.Info("Crate watcher stopped")
	os.Exit(0)

}
