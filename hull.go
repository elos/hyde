package hyde

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

type Hull struct {
	router *httprouter.Router
	server *http.Server

	pods map[string]*Pod
}

func NewHull(addr string, dirs ...string) *Hull {
	r := httprouter.New()

	r.ServeFiles("/assets/*filepath", http.Dir("../assets"))

	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	h := &Hull{
		router: r,
		server: s,

		pods: make(map[string]*Pod),
	}

	for i := range dirs {
		h.AddDir(dirs[i])
	}

	return h
}

func (h *Hull) AddDir(dir string) {
	t, _ := template.ParseFiles("../layout.tmpl")
	e, err := NewEngine(dir)
	if err != nil {
		return
	}

	go e.Start()
	p := NewPod(e, t, nil)
	go p.Start()

	h.router.GET(filepath.Join("/", p.Name(), "*subpath"), h.dispense(p))

	h.pods[p.Name()] = p
}

func (h *Hull) Start() {
	log.Fatal(h.server.ListenAndServe())
}

func (h *Hull) dispense(p *Pod) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		p.Router.ServeHTTP(w, r)
	}
}
