// Utility functions for Crate

package crate

import (
	"crypto/sha1"
	"encoding/base64"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rakyll/magicmime"
)

const (
	UnknownHost = "unknown"
	JSONLayout  = "2006-01-02T15:04:05-07:00"
)

var (
	hostname string
	Magic    *magicmime.Magic
)

//=============================================================================

func InitMagic() error {
	var err error
	Magic, err = magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
	return err
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

// Return the hostname of the machine
func Hostname() string {
	var err error
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			hostname = UnknownHost
		}
	}

	return hostname
}

// Convert a Float to a string
func Ftoa(num float64) string {
	if num == 0.0 {
		return ""
	}

	return strconv.FormatFloat(num, 'f', -1, 64)
}

// Convert a time.Time to a JSON timestamp
func JSONStamp(t time.Time) string {
	if !t.IsZero() {
		return t.UTC().Format(JSONLayout)
	}

	return ""
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
