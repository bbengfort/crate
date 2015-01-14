// Implements database handling for Crate

package crate

import (
	"encoding/json"

	"github.com/bbengfort/crate/crate/config"
	"github.com/syndtr/goleveldb/leveldb"
	// dbutil "github.com/syndtr/goleveldb/leveldb/util"
)

var db *leveldb.DB // Global var for the storage

//=============================================================================

// Initialize the database in the config location
func InitializeDatabase() error {
	var err error
	var path string

	if path, err = config.CrateDatabasePath(); err == nil {
		if db, err = leveldb.OpenFile(path, nil); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Close the database -- useful for defering close
func CloseDatabase() error {
	return db.Close()
}

//=============================================================================

// Writes the FileMeta to the database, where the key is the SHA1 hash
func (fm *FileMeta) Store() error {

	if !fm.populated {
		fm.Populate()
	}

	return db.Put([]byte(fm.Signature), fm.Byte(), nil)

}

// Writes the ImageMeta to the database, where the key is the SHA1 hash
func (img *ImageMeta) Store() error {

	if !img.populated {
		img.Populate()
	}

	return db.Put([]byte(img.Signature), img.Byte(), nil)

}

//=============================================================================

// Fetch the FileMeta or ImageMeta from the specified hash
func Fetch(key string) (FilePath, error) {

	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	meta := new(FileMeta)
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	meta.populated = true

	if meta.IsImage() {

		img := new(ImageMeta)
		err = json.Unmarshal(data, &img)
		if err != nil {
			return nil, err
		}

		img.populated = true
		return img, nil
	}

	return meta, nil

}

// Return a slice of keys limited by the argument
func FetchKeys(limit int) []string {

	result := make([]string, 0, 0)
	iter := db.NewIterator(nil, nil)
	idx := 0

	for iter.Next() {
		result = append(result, string(iter.Key()))

		idx++
		if idx > limit {
			break
		}

	}

	return result
}
