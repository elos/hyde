package hyde

import (
	"path/filepath"
	"strings"

	"github.com/elos/ehttp/templates"
)

var (
	hydeImportPath = "github.com/elos/hyde"
	assetsDir      = filepath.Join(templates.PackagePath(hydeImportPath), "assets")
)

type (
	FileMap map[string]*FileNode

	FileNode struct {
		Path, Root string
	}
)

func NewFileNode(path string, root string) *FileNode {
	return &FileNode{
		Path: path,
		Root: root,
	}
}

func (fn *FileNode) Name() string {
	return filepath.Base(fn.Path)
}

func (fn *FileNode) RelPath() string {
	return strings.TrimPrefix(fn.Path, fn.Root)
}
