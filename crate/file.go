// Implements File and Directory handling for the Crate watcher

package crate

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//=============================================================================

// A file system entity
type Node struct {
	Path string // Current path of the node
}

type Path interface {
	IsDir() bool                // Path is a directory
	IsFile() bool               // Path is a file
	IsHidden() bool             // Path is a hidden file or directory
	Dir() *Dir                  // The parent directory of the path
	Stat() (os.FileInfo, error) // Returns the attributes of the path
	User() (*user.User, error)  // Returns the User object for the path
	String() string             // The string representation of the file
	Byte() []byte               // The byte representation of the JSON
}

type FilePath interface {
	Ext() string  // The extension (if a file, empty string if not)
	Base() string // The base name of the path
	Populate()    // Populates the info on the file path (does a lot of work)
	Info() string // Returns a JSON serialized print of the file info
}

type DirPath interface {
	Join(elem ...string) string // Join path elements to the current path
	List() ([]Path, error)      // Return a list of the Paths in the directory
	Walk(walkFn WalkFunc) error // Walk a directory with the walk function
	Populate()                  // Populates the info on the dir path (does a lot of work)
}

// Type of the Walk Function for DirPath.Walk
type WalkFunc func(path Path, err error) error

//=============================================================================

// Create either a FileMeta or a Dir from a pathname
func NewPath(path string) (Path, error) {
	path = filepath.Clean(path)
	finfo, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if finfo.IsDir() {
		node := new(Dir)
		node.Path = path
		return node, nil
	} else {
		node := new(FileMeta)
		node.Path = path
		return node, nil
	}
}

// Check if a string pathname exists (prerequsite to NewPath)
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

//=============================================================================

func (node *Node) IsDir() bool {
	finfo, _ := node.Stat()
	if finfo != nil {
		return finfo.IsDir()
	}
	return false
}

func (node *Node) IsFile() bool {
	return !node.IsDir()
}

func (node *Node) IsHidden() bool {
	stat, err := node.Stat()
	if err != nil {
		return false
	}

	name := stat.Name()
	if name == "." || name == ".." {
		return false
	}

	return strings.HasPrefix(name, ".")
}

func (node *Node) Stat() (os.FileInfo, error) {
	return os.Stat(node.Path)
}

func (node *Node) User() (*user.User, error) {
	fi, ferr := node.Stat()
	if ferr != nil {
		return nil, ferr
	}

	var uid uint64
	sys := fi.Sys()
	if sys != nil {
		tsys, ok := sys.(*syscall.Stat_t)
		if ok {
			uid = uint64(tsys.Uid)
		}
	} else {
		uid = uint64(os.Geteuid())
	}

	if uid != 0 {
		return user.LookupId(strconv.FormatUint(uid, 10))
	} else {
		return nil, errors.New("unknown user")
	}

}

func (node *Node) Dir() *Dir {
	path := filepath.Dir(node.Path)
	dir := new(Dir)
	dir.Path = path
	return dir
}

func (node *Node) String() string {
	return node.Path
}

func (node *Node) Byte() []byte {
	data, err := json.Marshal(node)
	if err != nil {
		return nil
	}

	return data
}

//=============================================================================

type FileMeta struct {
	Node
	MimeType  string    // The mimetype of the file
	Name      string    // The base name of the file
	Size      int64     // The size of the file in bytes
	Modified  time.Time // The last modified time
	Signature string    // Base64 encoded SHA1 hash of the file
	Host      string    // The hostname of the computer
	Author    string    // The User or username of the file creator
	populated bool      // Indicates if the FileMeta has been populated
}

func (fm *FileMeta) Populate() {

	if fi, err := fm.Stat(); err == nil {
		fm.Name = fi.Name()
		fm.Size = fi.Size()
		fm.Modified = fi.ModTime()
	}

	if user, err := fm.User(); err == nil {
		fm.Author = user.Name
	}

	fm.Host = Hostname()
	fm.MimeType, _ = MimeType(fm.Path)
	fm.Signature, _ = fm.Hash()
	fm.populated = true
}

// Returns the extension of the file
func (fm *FileMeta) Ext() string {
	return filepath.Ext(fm.Path)
}

// Returns the basename of the file (including extension)
func (fm *FileMeta) Base() string {
	return filepath.Base(fm.Path)
}

// Computes the SHA1 hash of the file by using IO copy for memory safety
func (fm *FileMeta) Hash() (string, error) {
	file, err := os.Open(fm.Path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// Returns the byte serialization of the file meta for storage
func (fm *FileMeta) Byte() []byte {
	data, err := json.Marshal(fm)
	if err != nil {
		return nil
	}

	return data
}

// Prints out the info as a JSON indented pretty string
func (fm *FileMeta) Info() string {

	if !fm.populated {
		fm.Populate()
	}

	info, err := json.MarshalIndent(fm, "", "  ")
	if err != nil {
		return ""
	}

	return string(info)
}

//=============================================================================

type Dir struct {
	Node
	Name      string    // The base name of the directory
	Modified  time.Time // The modified time of the directory
	populated bool      // Whether or not the dir has been populated
}

func (dir *Dir) Populate() {
	if fi, err := dir.Stat(); err == nil {
		dir.Name = fi.Name()
		dir.Modified = fi.ModTime()
	}

	dir.populated = true
}

func (dir *Dir) Join(elem ...string) string {
	subdir := filepath.Join(elem...)
	return filepath.Join(dir.Path, subdir)
}

func (dir *Dir) List() ([]Path, error) {

	names, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		return nil, err
	}

	paths := make([]Path, len(names))
	for idx, finfo := range names {
		path := dir.Join(finfo.Name())

		if finfo.IsDir() {
			node := new(Dir)
			node.Path = path
			paths[idx] = node
		} else {
			node := new(FileMeta)
			node.Path = path
			paths[idx] = node
		}

	}

	return paths, nil
}

// Implements a recrusive walk of a directory
func (dir *Dir) Walk(walkFn WalkFunc) error {

	return filepath.Walk(dir.Path, func(path string, finfo os.FileInfo, err error) error {
		if finfo.IsDir() {
			node := new(Dir)
			node.Path = path
			return walkFn(node, err)

		} else {
			node := new(FileMeta)
			node.Path = path
			return walkFn(node, err)
		}

		return nil
	})

}

func (dir *Dir) Byte() []byte {
	data, err := json.Marshal(dir)
	if err != nil {
		return nil
	}

	return data
}
