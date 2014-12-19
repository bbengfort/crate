// Testing for the file and directory handling module

package crate

import (
	"testing"
)

func TestNodeInterface(t *testing.T) {
	var _ Path = &Node{}
}

func TestFileMetaInterface(t *testing.T) {
	var _ Path = &FileMeta{}
	var _ FilePath = &FileMeta{}
}

func TestDirInterface(t *testing.T) {
	var _ Path = &Dir{}
	var _ DirPath = &Dir{}
}
