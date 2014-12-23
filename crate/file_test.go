package crate_test

import (
	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("File", func() {

	Describe("Node", func() {
		It("should be a Path", func() {
			var _ Path = &Node{}
		})
	})

	Describe("FileMeta", func() {
		It("should be a Path and a FilePath", func() {
			var _ Path = &FileMeta{}
			var _ FilePath = &FileMeta{}
		})
	})

	Describe("Dir", func() {
		It("should be a Path and a DirPath", func() {
			var _ Path = &Dir{}
			var _ DirPath = &Dir{}
		})
	})

})
