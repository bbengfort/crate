// Implements File and Directory handling for the Crate watcher

package crate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	String() string             // The string representation of the file
	Dir() *Dir                  // The parent directory of the path
	Stat() (os.FileInfo, error) // Returns the attributes of the path
}

type FilePath interface {
	Ext() string  // The extension (if a file, empty string if not)
	Base() string // The base name of the path
}

type DirPath interface {
	Join(elem ...string) string // Join path elements to the current path
	List() ([]Path, error)      // Return a list of the Paths in the directory
	Walk(walkFn WalkFunc) error // Walk a directory with the walk function
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

	// console := &Console{true}
	// console.Info("%s", stat.Name())

	return strings.HasPrefix(stat.Name(), ".")
}

func (node *Node) Stat() (os.FileInfo, error) {
	return os.Stat(node.Path)
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

//=============================================================================

type FileMeta struct {
	Node
}

func (fm *FileMeta) Ext() string {
	return filepath.Ext(fm.Path)
}

func (fm *FileMeta) Base() string {
	return filepath.Base(fm.Path)
}

//=============================================================================

type Dir struct {
	Node
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
