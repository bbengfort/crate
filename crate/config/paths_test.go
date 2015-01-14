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

var _ = Describe("Config Paths", func() {

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

		// Clear the Cache
		ClearPathCache()
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

	It("should correctly get but not initialize the database path", func() {
		var expected string
		if runtime.GOOS == "windows" {
			expected = filepath.Join(testHome, "AppData", "Roaming", "Crate", "filemeta.db")
		} else {
			expected = filepath.Join(testHome, ".crate", "filemeta.db")
		}

		Ω(PathExists(expected)).Should(BeFalse())
		Ω(CrateDatabasePath()).Should(Equal(expected))
		Ω(PathExists(expected)).Should(BeFalse())
	})

	It("should correctly get and initialize the crate configuration", func() {
		var expected string
		if runtime.GOOS == "windows" {
			expected = filepath.Join(testHome, "AppData", "Roaming", "Crate", "config.yaml")
		} else {
			expected = filepath.Join(testHome, ".crate", "config.yaml")
		}

		Ω(PathExists(expected)).Should(BeFalse())
		Ω(CrateConfigPath()).Should(Equal(expected))
		Ω(PathExists(expected)).Should(BeTrue())
	})

	It("should correctly get and initialize the crate logging directory but not the log file", func() {
		var logDir string
		if runtime.GOOS == "windows" {
			logDir = filepath.Join(testHome, "AppData", "Roaming", "Crate", "logs")
		} else {
			logDir = filepath.Join(testHome, ".crate", "logs")
		}

		logPath := filepath.Join(logDir, "events.log")

		Ω(PathExists(logDir)).Should(BeFalse())
		Ω(PathExists(logPath)).Should(BeFalse())

		Ω(CrateLoggingPath()).Should(Equal(logPath))
		Ω(PathExists(logDir)).Should(BeTrue())
		Ω(PathExists(logPath)).ShouldNot(BeTrue())
	})

	It("should not overwite an existing crate configuration", func() {

		path, err := CrateDirectory()
		Ω(err).Should(BeNil())

		// Create a stub config.yaml file that shouldn't be overwritten
		Ω(os.Create(filepath.Join(path, "config.yaml"))).ShouldNot(BeZero())

		// Attempt to get config path (errors expected)
		path, err = CrateConfigPath()
		Ω(path).Should(BeZero())
		Ω(err).ShouldNot(BeNil())

	})

})
