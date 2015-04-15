package hyde

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/elos/autonomous"
	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday"
)

type Server struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	*Engine
	*http.Server
	*httprouter.Router
	registeredRoutes map[string]httprouter.Handle
}

type TreeNav struct {
	Name, Link string
	*FileNode
	Subs map[string]TreeNav
}

func BuildTreeNav(e *Engine) TreeNav {
	t := TreeNav{Subs: make(map[string]TreeNav)}
	rootPath := e.RootDir()
	for _, v := range e.FileMap() {
		t.add(v, rootPath)
	}
	return t
}

func (t TreeNav) add(fn *FileNode, rootPath string) {
	routeName := strings.TrimPrefix(fn.Path, rootPath)
	splits := strings.Split(routeName, "/")

	tt := t
	for _, level := range splits {
		_, ok := tt.Subs[level]

		if ok {
			tt = tt.Subs[level]
			continue
		}

		nt := TreeNav{
			Name:     level,
			Link:     routeName,
			FileNode: fn,
			Subs:     make(map[string]TreeNav),
		}

		tt.Subs[level] = nt
		tt = nt
	}

	tt.FileNode = fn
}

func NewServer(e *Engine) *Server {
	r := httprouter.New()

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	r.ServeFiles("/assets/*filepath", http.Dir("../assets"))

	return &Server{
		Life:    autonomous.NewLife(),
		Stopper: make(autonomous.Stopper),

		Engine:           e,
		Server:           s,
		Router:           r,
		registeredRoutes: make(map[string]httprouter.Handle),
	}
}

type Page struct {
	Nav     TreeNav
	Content template.HTML
}

func (s *Server) Start() {
	s.Life.Begin()
	go s.Listen()

Run:
	for {
		select {
		case fn := <-s.Engine.NodeChanges:
			s.loadNode(fn)
		case <-s.Stopper:
			break Run
		}
	}
}

func (s *Server) Listen() {
	log.Fatal(s.Server.ListenAndServe())
}

func (s *Server) loadNode(fn *FileNode) {
	log.Print(fn.Path)
	routeName := strings.TrimPrefix(fn.Path, s.Engine.RootDir())
	if routeName == "" {
		return
	}

	_, ok := s.registeredRoutes[routeName]
	s.registeredRoutes[routeName] = s.TemplateHandler(fn)
	if !ok {
		s.Router.GET(routeName, s.Handle(routeName))
	}
}

func (s *Server) Handle(path string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h, ok := s.registeredRoutes[path]
		if ok {
			h(w, r, p)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (s *Server) TemplateHandler(fn *FileNode) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout, err = template.ParseFiles("../layout.tmpl")
		if err != nil {
			log.Print(err)
			return
		}

		input, _ := ioutil.ReadFile(fn.Path)
		output := blackfriday.MarkdownCommon(input)
		layout.Execute(w, Page{
			Nav:     BuildTreeNav(s.Engine),
			Content: template.HTML(output),
		})
	}
}

func New(path string) (*Server, error) {
	log.Print(path)
	log.Print("1")
	e, err := NewEngine(path)
	log.Print("2")
	if err != nil {
		return nil, err
	}
	log.Print("3")

	s := NewServer(e)
	log.Print("3")
	go e.Start()
	return s, nil
}

func Serve(e *Engine) {
	s := new(Server)
	s.Engine = e
	s.Server = new(http.Server)
	s.Router = httprouter.New()
}
