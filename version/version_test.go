// Check the version during testing

package version

import (
	"testing"
)

const (
	ExpectedVersion = "0.0.1"
)

func TestVersion(t *testing.T) {
	if Version() != ExpectedVersion {
		t.Error("Version mismatch package version does not match test version!")
	}
}
