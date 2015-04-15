package hyde

import "path/filepath"

type FileMap map[string]*FileNode

type FileNode struct {
	Path, Name string
}

func NewFileNode(path string) *FileNode {
	return &FileNode{
		Name: filepath.Base(path),
		Path: path,
	}
}
