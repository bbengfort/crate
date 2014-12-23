// Utility functions for Crate

package crate

import (
	"crypto/sha1"
	"encoding/base64"
	"os"
	"os/signal"
	"syscall"

	"github.com/rakyll/magicmime"
)

var Magic *magicmime.Magic

//=============================================================================

func InitMagic() {
	Magic, _ = magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
}

// Use libmagic to determine the MimeType of the file
func MimeType(path string) (string, error) {
	if Magic == nil {
		InitMagic()
	}

	return Magic.TypeByFile(path)
}

// Compute the Base64 encoded SHA1 hash of the data
func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

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
