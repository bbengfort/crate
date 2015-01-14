package crate_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/bbengfort/crate/crate"
	"github.com/bbengfort/crate/crate/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {

	var (
		err      error  // Any errors in directory creation
		testRoot string // Test directory to store temp fixtures
		testHome string // Fake home directory in temp directory
		logPath  string // Test directory to store temp fixtures
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

		logPath, err = config.CrateLoggingPath()
		Ω(err).Should(BeNil())

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
		config.ClearPathCache()
	})

	It("should log to a file", func() {
		err := InitializeLoggers(LevelInfo)
		Ω(err).Should(BeNil())
		defer CloseLoggers()

		Log("test log message", LevelInfo)
		Ω(config.PathExists(logPath)).Should(BeTrue())

		data, err := ioutil.ReadFile(logPath)
		Ω(err).Should(BeNil())

		logPattern := `INFO    \[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}\]: test log message\n`
		Ω(string(data)).Should(MatchRegexp(logPattern))
	})

	It("should not log if the level is less than minimum level", func() {
		err := InitializeLoggers(LevelWarn)
		Ω(err).Should(BeNil())
		defer CloseLoggers()

		Log("test log message", LevelInfo)
		Log("test warning message", LevelWarn)
		Ω(config.PathExists(logPath)).Should(BeTrue())

		data, err := ioutil.ReadFile(logPath)
		Ω(err).Should(BeNil())

		logPattern := `WARNING \[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}\]: test warning message\n`
		Ω(string(data)).Should(MatchRegexp(logPattern))
	})

	It("should be able to convert a string to a LogLevel", func() {

		Ω(LevelFromString("debug")).Should(Equal(LevelDebug))
		Ω(LevelFromString("info")).Should(Equal(LevelInfo))
		Ω(LevelFromString("warning")).Should(Equal(LevelWarn))
		Ω(LevelFromString("error")).Should(Equal(LevelError))
		Ω(LevelFromString("fatal")).Should(Equal(LevelFatal))
		Ω(LevelFromString("")).Should(Equal(LevelInfo))

	})

	It("should be able to log with format", func() {
		err := InitializeLoggers(LevelInfo)
		Ω(err).Should(BeNil())
		defer CloseLoggers()

		Log("test message %d", LevelInfo, 1)
		Ω(config.PathExists(logPath)).Should(BeTrue())

		data, err := ioutil.ReadFile(logPath)
		Ω(err).Should(BeNil())

		logPattern := `INFO    \[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[-+]\d{2}:\d{2}\]: test message 1\n`
		Ω(string(data)).Should(MatchRegexp(logPattern))
	})

})
