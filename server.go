package hyde

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
}

type TreeNav struct {
	This string
	*FileNode
	Subs map[string]TreeNav
}

func BuildTreeNav(e *Engine) TreeNav {
	t := TreeNav{This: filepath.Base(e.RootDir()), Subs: make(map[string]TreeNav)}
	rootPath := e.RootDir()
	for _, v := range e.FileMap() {
		t.add(v, rootPath)
	}
	return t
}

func (t TreeNav) add(fn *FileNode, rootPath string) {
	splits := strings.Split(strings.TrimPrefix(fn.Path, rootPath), "/")

	subpath := rootPath

	tt := t

	for i := range splits {
		_, ok := tt.Subs[splits[i]]
		if ok {
			tt = tt.Subs[splits[i]]
			continue
		}

		name := splits[i]
		file, _ := os.Stat(subpath)
		if file.IsDir() {
			name = name + "/"
		}

		nt := TreeNav{This: name, Subs: make(map[string]TreeNav)}
		tt.Subs[splits[i]] = nt
		tt = nt

		subpath = filepath.Join(subpath, splits[i])
	}

	tt.FileNode = fn
}

func NewServer(e *Engine) *Server {
	r := httprouter.New()

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	return &Server{
		Life:    autonomous.NewLife(),
		Stopper: make(autonomous.Stopper),

		Engine: e,
		Server: s,
		Router: r,
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
	log.Printf("loadNode: %s", fn.Name)
	s.Router.GET("/"+fn.Name, s.TemplateHandler(fn))
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
	e, err := NewEngine(path)
	if err != nil {
		return nil, err
	}

	s := NewServer(e)
	go e.Start()
	return s, nil
}

func Serve(e *Engine) {
	s := new(Server)
	s.Engine = e
	s.Server = new(http.Server)
	s.Router = httprouter.New()
}
