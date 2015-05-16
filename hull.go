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

type Hull struct {
	serve.Router

	pods map[string]*Pod
}

func NewHull(dirs ...string) *Hull {
	r := builtin.NewRouter()
	return NewHullWithRouter(r, dirs...)
}

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
