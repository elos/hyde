package hyde

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/elos/autonomous"
	"github.com/go-fsnotify/fsnotify"
)

// An Engine is a slightly glorified fsnotify.Watcher, it watches
// a root directory recursively (meaning it watches all the directories
// within that too) and sends updates (relatively abstracted - just
// changes and removes) regarding the file system on its channels
// NodeChanges, NodeRemoves.
type Engine struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	w        *fsnotify.Watcher
	RootedAt string
	fmap     *FileMap

	NodeChanges chan *FileNode
	NodeRemoves chan *FileNode
}

// FileMap returns A map of the files this engine is aware off
func (e *Engine) FileMap() FileMap {
	return *e.fmap
}

// NewEngine allocates a new Engine, and recursively loads
// the files at the given path, the given path becomes the
// root directory of the engine
// All initial FileNodes will be sent over NodeChanges
func NewEngine(atPath string) (*Engine, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fm := make(FileMap)

	e := &Engine{
		fmap:     &fm,
		w:        watcher,
		RootedAt: atPath,

		Life:        autonomous.NewLife(),
		Stopper:     make(autonomous.Stopper),
		NodeChanges: make(chan *FileNode, 20),
	}

	e.load(e.RootedAt)

	return e, nil
}

// watch adds a path to the fsnotify.Watcher
func (e *Engine) watch(path string) {
	e.w.Add(path)
}

// load recursively imports a given path and all its visible
// (meaning not .files) descendants into the engine
// sends all new FileNodes over NodeChanges
func (e *Engine) load(path string) {
	file, err := os.Stat(path)
	if err != nil {
		return
	}

	base := filepath.Base(path)
	if len(base) > 0 && string(base[0]) == "." {
		return // cause file is hidden
	}

	if file.IsDir() {
		e.watch(path)
		files, err := ioutil.ReadDir(path)
		if err == nil {
			for _, f := range files {
				e.load(filepath.Join(path, f.Name()))
			}
		}
	}

	node := NewFileNode(path, e.RootedAt)
	(*e.fmap)[path] = node
	go func() { e.NodeChanges <- node }()
}

// remove sends the FileNode over NodeRemoves and
// deletes if from the FileMap
func (e *Engine) remove(path string) {
	node, ok := (*e.fmap)[path]

	if ok {
		go func() { e.NodeRemoves <- node }()
		delete(*e.fmap, path)
	}
}

// Start begins an engine's life, includes listening
// for events from the fsnotify.Watcher
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
			e.process(event)
		case err := <-errors:
			log.Printf("Hyde's fsnotify failed:", err)
			go e.Stop()
		case <-e.Stopper:
			break Run
		}
	}

	e.w.Close()
	e.Life.End()
}

// process routes an fsnotfiy.Event to the appropriate handler
// function.
// Create, Update -> load
// Remove, Rename -> remove
// Note: when a file is renamed fsnotify sends one event
// for the rename, one event for the new name create
func (e *Engine) process(event *fsnotify.Event) {
	switch event.Op {
	case fsnotify.Create:
		e.load(event.Name)
	case fsnotify.Write:
		e.load(event.Name)
	case fsnotify.Remove:
		e.remove(event.Name)
	case fsnotify.Rename:
		e.remove(event.Name)
	}
}
