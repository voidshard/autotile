package autotile

import (
	"path/filepath"

	"github.com/voidshard/tile"
)

// Loader defines some way of getting a tile.Map from a string.
// Loaders are required to be thread safe as we'll be calling
// Map() with many routines.
type Loader interface {
	Map(string) (*tile.Map, error)
}

// FileLoader is the most straight forward kind of loader; it assumes
// that given strings are file paths on disk.
type FileLoader struct {
	root string
}

// NewFileLoader creates a new FileLoader based in `root`
func NewFileLoader(root string) *FileLoader {
	return &FileLoader{root: root}
}

// Map loads tilemap from file on disk
func (f *FileLoader) Map(in string) (*tile.Map, error) {
	return tile.Open(filepath.Join(f.root, in))
}
