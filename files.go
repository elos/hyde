package hyde

import (
	"path/filepath"
	"strings"

	"github.com/elos/httpserver/templates"
)

type FileMap map[string]*FileNode

type FileNode struct {
	Path, Root string
}

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

var (
	hydeImportPath = "github.com/elos/hyde"
	assetsDir      = filepath.Join(templates.PackagePath(hydeImportPath), "assets")
)
