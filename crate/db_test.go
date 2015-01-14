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

var _ = Describe("Db", func() {

	const (
		coastPath   = "coast.jpg"
		draculaPath = "dracula.txt"
	)

	var (
		err      error      // Any errors in directory creation
		testRoot string     // Test directory to store temp fixtures
		testHome string     // Fake home directory in temp directory
		fixtures *Dir       //  The Directory containing test fixtures
		coast    *ImageMeta // An ImageMeta for storage
		dracula  *FileMeta  // A FileMEta for storage
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

		// Locate the fixtures to test on
		if exists, _ := PathExists("./fixtures"); exists {
			fpath, _ := NewPath("./fixtures")
			fixtures = fpath.(*Dir)
		} else if exists, _ := PathExists("../fixtures"); exists {
			fpath, _ := NewPath("../fixtures/")
			fixtures = fpath.(*Dir)
		}

		Ω(fixtures).ShouldNot(BeNil())

		// Add the Path fixtures for use in testing
		coastp, _ := NewPath(fixtures.Join(coastPath))
		coast, _ = ConvertImageMeta(coastp.(*FileMeta))
		Ω(coast).ShouldNot(BeNil())
		draculap, _ := NewPath(fixtures.Join(draculaPath))
		dracula = draculap.(*FileMeta)
		Ω(dracula).ShouldNot(BeNil())

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

	It("should initialize the database in the data directory", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		dbpath, err := config.CrateDatabasePath()
		Ω(err).Should(BeNil())
		Ω(config.PathExists(dbpath)).Should(BeTrue())
	})

	It("should be able to store a filemeta in the database", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(dracula.Store()).Should(BeNil())
		Ω(FetchKeys(2)).Should(ContainElement(dracula.Signature))
	})

	It("should be able to store an imagemeta in the database", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(coast.Store()).Should(BeNil())
		Ω(FetchKeys(2)).Should(ContainElement(coast.Signature))
	})

	It("should be able to fetch a filemeta in the database", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(dracula.Store()).Should(BeNil())

		test, err := Fetch(dracula.Signature)
		Ω(err).Should(BeNil())
		_, ok := test.(*FileMeta)
		Ω(ok).Should(BeTrue())

	})

	It("should be able to fetch an imagemeta in the database", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(coast.Store()).Should(BeNil())

		test, err := Fetch(coast.Signature)
		Ω(err).Should(BeNil())
		img, ok := test.(*ImageMeta)
		Ω(ok).Should(BeTrue())
		Ω(img.Width).ShouldNot(BeZero())
		Ω(img.Height).ShouldNot(BeZero())
		Ω(img.Tags).ShouldNot(BeZero())
	})

	It("should be able to fetch keys from the database", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(FetchKeys(100)).Should(BeEmpty())

		Ω(coast.Store()).Should(BeNil())
		Ω(dracula.Store()).Should(BeNil())

		Ω(FetchKeys(100)).Should(HaveLen(2))
	})

	It("should be able to limit key fetch", func() {
		Ω(InitializeDatabase()).Should(BeNil())
		defer CloseDatabase()

		Ω(coast.Store()).Should(BeNil())
		Ω(dracula.Store()).Should(BeNil())

		Ω(FetchKeys(1)).Should(HaveLen(1))
	})

})
