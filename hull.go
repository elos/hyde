package hyde

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Hull struct {
	router *httprouter.Router
	server *http.Server

	pods map[string]*Pod
}

func NewHull(addr string, dirs ...string) *Hull {
	r := httprouter.New()
	return NewHullWithRouter(addr, r, dirs...)
}

func NewHullWithRouter(addr string, r *httprouter.Router, dirs ...string) *Hull {
	r.ServeFiles("/assets/*filepath", http.Dir(assetsDir))

	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	h := &Hull{
		router: r,
		server: s,

		pods: make(map[string]*Pod),
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

	h.router.GET(filepath.Join("/", p.Name(), "/*subpath"), p.Route)

	h.pods[p.Name()] = p
}

func (h *Hull) Start() {
	log.Fatal(h.server.ListenAndServe())
}

type PodView struct {
	Name, Link string
}

func (h *Hull) index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	t, err := template.ParseFiles(filepath.Join(assetsDir, "templates/root.tmpl"), filepath.Join(assetsDir, "templates/hull.tmpl"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pods := make([]*PodView, len(h.pods))

	var i int = 0
	for k, v := range h.pods {
		pods[i] = &PodView{
			Name: strings.ToUpper(v.Name()),
			Link: filepath.Join("/", k),
		}
		i++
	}

	t.Execute(w, pods)
}
