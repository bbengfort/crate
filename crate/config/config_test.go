package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/bbengfort/crate/crate/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		err      error  // Any errors in directory creation
		testRoot string // Test directory to store temp fixtures
	)

	BeforeEach(func() {
		// Setup the temp test root directory
		testRoot, err = ioutil.TempDir("", "ginkgo-")
		Ω(err).Should(BeNil())
	})

	AfterEach(func() {
		// Remove the test file system
		err = os.RemoveAll(testRoot)
		Ω(err).Should(BeNil())
	})

	It("should create a new config with defaults", func() {
		conf := New()

		// Test the defaults
		Ω(conf.Debug).Should(BeFalse())  // Debug is false
		Ω(conf.Notify).Should(BeEmpty()) // Notify is empty list
	})

	It("should be able to dump a config to disk", func() {
		conf := New()
		conf.Debug = true
		conf.Notify = append(conf.Notify, "joe@example.com")
		conf.Notify = append(conf.Notify, "jane@example.com")
		out := filepath.Join(testRoot, "config.yaml")

		Ω(conf.Dump(out)).Should(BeNil())
		Ω(PathExists(out)).Should(BeTrue())
	})

	It("should be able to load a config from disk", func() {
		conf := New()
		conf.Debug = true
		conf.Notify = append(conf.Notify, "joe@example.com")
		out := filepath.Join(testRoot, "config.yaml")

		Ω(conf.Dump(out)).Should(BeNil())
		Ω(PathExists(out)).Should(BeTrue())

		// Load the config
		config, err := Load(out)
		Ω(err).Should(BeNil())
		Ω(config).ShouldNot(BeNil())
		Ω(config.Debug).Should(BeTrue())
		Ω(config.Notify).Should(HaveLen(1))
		Ω(config.Notify).Should(ContainElement("joe@example.com"))
	})

})
