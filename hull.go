package hyde

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/elos/ehttp/builtin"
	"github.com/elos/ehttp/serve"
)

// A hull is the abstraction layer for an entire section of
// documentation, the pods of a hull are subsections of sort
// i.e., http documentation with api, app and security pods.
type Hull struct {
	// Hull satisfies http.Handler interface
	serve.Router
	// the pods this hub is aware of
	pods map[string]*Pod
}

// NewHull is a helper handler for NewHullWithRouter,
// it uses the elos/ehttp/builtin router
// Dirs is the directory roots for each of the hull's pods
func NewHull(dirs ...string) *Hull {
	r := builtin.NewRouter()
	return NewHullWithRouter(r, dirs...)
}

// NewHullWithRouter constructures a hull using the provided
// server.Router. The dirs are the root directories of the hull's
// pods. The hull will imediately serve all hyde assets at
// assets/*filepath. It will also server the hull index at h.index
// It is a panic of you have already handled the "/" route on the Router.
func NewHullWithRouter(r serve.Router, dirs ...string) *Hull {
	r.ServeFiles("/assets/*filepath", http.Dir(assetsDir))

	h := &Hull{
		Router: r,
		pods:   make(map[string]*Pod),
	}

	r.GET("/", h.index)

	for i := range dirs {
		h.AddDir(dirs[i])
	}

	return h
}

// AddDir is a means of adding root directories (as pods) to the hull post-creation
// NewHullWithRouter relies on this as a convenience method
func (h *Hull) AddDir(dir string) {
	t, err := template.ParseFiles(filepath.Join(assetsDir, "templates/root.tmpl"), filepath.Join(assetsDir, "templates/layout.tmpl"))
	if err != nil {
		log.Print(err)
		return
	}

	e, err := NewEngine(dir)
	if err != nil {
		log.Print(err)
		return
	}

	go e.Start()
	p := NewPod(e, t, nil)
	go p.Start()

	h.GET(filepath.Join("/", p.Name(), "/*subpath"), p.Route)

	h.pods[p.Name()] = p
}

type PodView struct {
	Name, Link string
}

func (h *Hull) index(c *serve.Conn) {
	t, err := template.ParseFiles(filepath.Join(assetsDir, "templates/root.tmpl"), filepath.Join(assetsDir, "templates/hull.tmpl"))

	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
		return
	}

	pods := make(map[string]*PodView)

	var i int = 0
	for k, _ := range h.pods {
		pods[k] = &PodView{
			Name: strings.ToUpper(k),
			Link: filepath.Join("/", k),
		}
		i++
	}

	t.Execute(c, pods)
}
