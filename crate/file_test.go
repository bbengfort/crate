package crate_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/bbengfort/crate/crate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("File", func() {

	var (
		err      error     // Any errors during execution
		testRoot string    // The path to the root directory
		testDir  *Dir      // The Dir object to test upon
		dpath    *Dir      // A Directory object to test on
		fpath    *FileMeta // A file path to test on
		spath    *FileMeta // A static file to test on
	)

	BeforeEach(func() {
		// Set up the test file system
		testRoot, err = ioutil.TempDir("", "ginkgo-")
		os.MkdirAll(filepath.Join(testRoot, "h1", "h2", "h3"), 0755)
		os.MkdirAll(filepath.Join(testRoot, "foo", "bar"), 0755)
		os.MkdirAll(filepath.Join(testRoot, ".secret"), 0755)

		// Write files at each level
		ioutil.WriteFile(filepath.Join(testRoot, ".secret", "time.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "foo", "now.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "foo", "hello.txt"), []byte("Hello world!"), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "foo", "bar", "later.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "aspect.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "object.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "h2", "tom.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "h2", "jerry.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "h2", "h3", "now.txt"), []byte(time.Now().String()), 0644)

		// Write hidden files at several levels
		ioutil.WriteFile(filepath.Join(testRoot, "foo", ".hidden.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, ".scary.txt"), []byte(time.Now().String()), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, "h1", "h2", ".secret.txt"), []byte(time.Now().String()), 0644)

		node, _ := NewPath(testRoot)
		testDir = node.(*Dir)

		dnode, _ := NewPath(filepath.Join(testRoot, "foo"))
		dpath = dnode.(*Dir)
		fnode, _ := NewPath(filepath.Join(testRoot, "foo", "now.txt"))
		fpath = fnode.(*FileMeta)
		snode, _ := NewPath(filepath.Join(testRoot, "foo", "hello.txt"))
		spath = snode.(*FileMeta)
	})

	AfterEach(func() {
		// Remove the test file system
		err = os.RemoveAll(testRoot)
	})

	It("should ensure testDir is actually a Dir", func() {
		dirPath := new(Dir)
		Ω(testDir).Should(BeAssignableToTypeOf(dirPath))
	})

	It("should tell if a path exists or not", func() {
		Ω(PathExists(filepath.Join(testRoot, "foo", "bar"))).Should(BeTrue())
		Ω(PathExists(filepath.Join(testRoot, "foo", "baz"))).Should(BeFalse())
	})

	Describe("Initialization", func() {
		// Test the NewPath function

		It("should create a FileMeta", func() {
			fm := new(FileMeta)
			path := filepath.Join(testRoot, "foo", "now.txt")
			Ω(NewPath(path)).Should(BeAssignableToTypeOf(fm))
		})

		It("should create a Dir", func() {
			dir := new(Dir)
			path := filepath.Join(testRoot, "foo")
			Ω(NewPath(path)).Should(BeAssignableToTypeOf(dir))
		})

		It("should return nil on error", func() {
			path, err := NewPath(filepath.Join(testRoot, "baz"))
			Ω(path).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		})
	})

	Describe("Node", func() {
		It("should be a Path", func() {
			var _ Path = &Node{}
		})

		It("should report if its a directory", func() {

			Ω(dpath.IsDir()).Should(BeTrue())
			Ω(fpath.IsDir()).Should(BeFalse())
		})

		It("should report if its a file", func() {
			Ω(dpath.IsFile()).Should(BeFalse())
			Ω(fpath.IsFile()).Should(BeTrue())
		})

		It("should report if it is hidden", func() {

			hf, _ := NewPath(filepath.Join(testRoot, ".scary.txt"))
			hd, _ := NewPath(filepath.Join(testRoot, ".secret"))

			Ω(hf.IsHidden()).Should(BeTrue())
			Ω(hd.IsHidden()).Should(BeTrue())

			Ω(dpath.IsHidden()).Should(BeFalse())
			Ω(fpath.IsHidden()).Should(BeFalse())

		})

		It("should not report current or parent dirs as hidden", func() {

			cwd, _ := NewPath(".")
			pwd, _ := NewPath("..")

			Ω(cwd.IsHidden()).Should(BeFalse())
			Ω(pwd.IsHidden()).Should(BeFalse())

		})

		It("should return a parent directory", func() {
			Ω(dpath.Dir()).Should(BeEquivalentTo(testDir))
			Ω(fpath.Dir()).Should(BeEquivalentTo(dpath))
		})

		It("should stringify to an absolute path", func() {
			Ω(testDir.String()).Should(Equal(testRoot))
		})

		It("should return a file info on stat", func() {
			Ω(dpath.Stat()).ShouldNot(BeNil())
			Ω(fpath.Stat()).ShouldNot(BeNil())
		})

		It("should return a user object", func() {
			Ω(dpath.User()).ShouldNot(BeNil())
			Ω(fpath.User()).ShouldNot(BeNil())
		})
	})

	Describe("FileMeta", func() {
		It("should be a Path and a FilePath", func() {
			var _ Path = &FileMeta{}
			var _ FilePath = &FileMeta{}
		})

		It("should return a correct extension", func() {
			Ω(fpath.Ext()).Should(Equal(".txt"))
		})

		It("should return a correct basename", func() {
			Ω(fpath.Base()).Should(Equal("now.txt"))
		})

		It("should compute a correct signature", func() {
			Ω(spath.Hash()).Should(Equal("00hq6RNueFa8QiEjhep5cJRHWAI="))
		})

		It("should begin unpopulated", func() {
			Ω(fpath.MimeType).Should(BeZero())
			Ω(fpath.Name).Should(BeZero())
			Ω(fpath.Size).Should(BeZero())
			Ω(fpath.Modified).Should(BeZero())
			Ω(fpath.Signature).Should(BeZero())
			Ω(fpath.Host).Should(BeZero())
			Ω(fpath.Author).Should(BeZero())
		})

		It("should be populated", func() {
			fpath.Populate()
			Ω(fpath.MimeType).ShouldNot(BeZero(), "MimeType must be populated")
			Ω(fpath.Name).ShouldNot(BeZero(), "Name must be populated")
			Ω(fpath.Size).ShouldNot(BeZero(), "Size must be populated")
			Ω(fpath.Modified).ShouldNot(BeZero(), "Modified must be populated")
			Ω(fpath.Signature).ShouldNot(BeZero(), "Signature must be populated")
			Ω(fpath.Host).ShouldNot(BeZero(), "Host must be populated")
			Ω(fpath.Author).ShouldNot(BeZero(), "Author must be populated")
		})
	})

	Describe("Dir", func() {
		It("should be a Path and a DirPath", func() {
			var _ Path = &Dir{}
			var _ DirPath = &Dir{}
		})

		It("should be able to join paths to itself", func() {
			Ω(testDir.Join("test", "to", "here")).Should(Equal(filepath.Join(testRoot, "test/to/here")))
		})

		It("should be able to list a directory", func() {

			Ω(dpath.List()).Should(HaveLen(4))
			Ω(dpath.List()).Should(ContainElement(fpath))
		})

		PIt("should be able to walk a directory", func() {
			// How to test Walk function?
		})
	})

})
