package hyde

import (
	"path/filepath"
	"strings"

	"github.com/elos/ehttp/templates"
)

var (
	// The canonical "go get"/import
	ImportPath = "github.com/elos/hyde"

	// the directory of hyde specific assets - eventually
	// should be minimal defaults
	assetsDir = filepath.Join(templates.PackagePath(ImportPath), "assets")
)

type (
	// FileMap is a map of path strings to FileNode pointers
	FileMap map[string]*FileNode

	// FileNode represents a watched file, knows of its root directory
	// and its absolute file system path. Root is useful for hyde's
	// navigation sidebar (e.g., don't start at "/" on *nix machines, start
	// at "elos/documentation/"). /Users/Nick/workspace/github.com/elos/documentation
	// would be the root in this case.
	FileNode struct {
		Path, Root string
	}
)

// NewFileNode allocates a NewFileNode, use this explicitly
// rather than a struct literal so that any future initialization
// logic can be applied here
func NewFileNode(path string, root string) *FileNode {
	return &FileNode{
		Path: path,
		Root: root,
	}
}

// Name returns the last element of the FileNode's path.
// (e.g., ~/github.com/elos/documentation/README.md -> "README.md")
func (fn *FileNode) Name() string {
	return filepath.Base(fn.Path)
}

// RelPath trims the FileNode's absolute path with from the root,
// (e.g.,
// if  path := ~/github.com/elos/documentation/server/README.md
// and root := ~/github.com/elos/documentation
// RelPath := "/server/README.md"
func (fn *FileNode) RelPath() string {
	return strings.TrimPrefix(fn.Path, fn.Root)
}
