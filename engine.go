package hyde

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/elos/autonomous"
	"github.com/go-fsnotify/fsnotify"
)

type Engine struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	w       *fsnotify.Watcher
	rootDir string
	root    *FileNode
	fmap    *FileMap

	NodeChanges chan *FileNode
}

func (e *Engine) FileMap() FileMap {
	return *e.fmap
}

func (e *Engine) RootDir() string {
	return e.rootDir
}

func NewEngine(rootPath string) (*Engine, error) {
	log.Print("dd")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	log.Print("bb")

	fm := make(FileMap)

	e := &Engine{
		fmap:    &fm,
		w:       watcher,
		rootDir: rootPath,

		Life:        autonomous.NewLife(),
		Stopper:     make(autonomous.Stopper),
		NodeChanges: make(chan *FileNode, 20),
	}
	log.Print("cc")

	e.loadPath(rootPath)
	log.Print("ee")

	return e, nil
}

func (e *Engine) loadPath(path string) *FileNode {
	// log.Print("loadpath: ", path)
	base := filepath.Base(path)
	if len(base) > 0 && string(base[0]) == "." {
		return nil
	}

	node := NewFileNode(path)
	(*e.fmap)[path] = node

	file, err := os.Stat(path)
	if err != nil {
		log.Print(err)
		return node
	}

	if file.IsDir() {
		e.w.Add(path)
		files, _ := ioutil.ReadDir(path)

		for i := range files {
			e.loadPath(filepath.Join(path, files[i].Name()))
		}
	}

	go func() { e.NodeChanges <- node }()
	return node
}

func (e *Engine) Start() {
	e.Life.Begin()

	events := make(chan *fsnotify.Event)
	errors := make(chan error)

	go func() {
		for {
			select {
			case event := <-e.w.Events:
				events <- &event
			case err := <-e.w.Errors:
				errors <- err
			}
		}
	}()

Run:
	for {
		select {
		case event := <-events:
			e.processFileEvent(event)
		case err := <-errors:
			log.Printf("Watcher error:", err)
			go e.Stop()
		case <-e.Stopper:
			break Run
		}
	}

	log.Print("engine died")
	e.w.Close()
	e.Life.End()
}

func (e *Engine) processFileEvent(event *fsnotify.Event) {
	log.Print(event)
	switch event.Op {
	case fsnotify.Create:
		e.loadPath(event.Name)
	case fsnotify.Write:
		e.loadPath(event.Name)
	case fsnotify.Remove:
		delete(*e.fmap, event.Name)
	case fsnotify.Rename:
		delete(*e.fmap, event.Name)
	}
}
