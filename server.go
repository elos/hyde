package hyde

import (
	"log"
	"net/http"

	"github.com/elos/autonomous"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	*Engine
	*http.Server
	*httprouter.Router
}

func NewServer(e *Engine) *Server {
	return &Server{
		Life:    autonomous.NewLife(),
		Stopper: make(autonomous.Stopper),

		Engine: e,
		Server: new(http.Server),
		Router: httprouter.New(),
	}
}

func (s *Server) Start() {
	s.Life.Begin()

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

func (s *Server) loadNode(fn *FileNode) {
	log.Print("load node")
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
