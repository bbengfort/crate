package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"

	. "github.com/bbengfort/crate/crate/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		err      error  // Any errors in directory creation
		testRoot string // Test directory to store temp fixtures
		testHome string // Fake home directory in temp directory
	)

	BeforeEach(func() {

		// Setup the temp test root directory
		testRoot, err = ioutil.TempDir("", "ginkgo-")
		Ω(err).Should(BeNil())

		// Setup the fake User home directory for testing
		testHome = filepath.Join(testRoot, "Users", "jdoe")
		err = os.MkdirAll(testHome, 0755)
		Ω(err).Should(BeNil())

		if runtime.GOOS == "windows" {
			err = os.Setenv("USERPROFILE", testHome)
			Ω(err).Should(BeNil())
		} else {
			err = os.Setenv("HOME", testHome)
			Ω(err).Should(BeNil())
		}

	})

	AfterEach(func() {
		// Remove the test file system
		err = os.RemoveAll(testRoot)
		Ω(err).Should(BeNil())

		// Unset the environment variables
		if runtime.GOOS == "windows" {
			err = os.Unsetenv("USERPROFILE")
			Ω(err).Should(BeNil())
		} else {
			err = os.Unsetenv("HOME")
			Ω(err).Should(BeNil())
		}
	})

	It("should correctly get the home directory", func() {
		Ω(homedir.Dir()).Should(Equal(testHome))
	})

	It("should correctly get and initialize the crate directory", func() {
		var expected string
		if runtime.GOOS == "windows" {
			expected = filepath.Join(testHome, "AppData", "Roaming", "Crate")
		} else {
			expected = filepath.Join(testHome, ".crate")
		}

		Ω(PathExists(expected)).Should(BeFalse())
		Ω(CrateDirectory()).Should(Equal(expected))
		Ω(PathExists(expected)).Should(BeTrue())
	})

})
