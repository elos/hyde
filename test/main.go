package main

import (
	"flag"
	"path/filepath"

	"github.com/elos/hyde"
)

var (
	watch = flag.Bool("w", true, "Watch files for changes")
)

func main() {
	p, _ := filepath.Abs("../../../../../../github.com/elos/documentation")
	agents := filepath.Join(p, "agents")
	data := filepath.Join(p, "data")
	http := filepath.Join(p, "http")
	server := filepath.Join(p, "server")

	p1, _ := filepath.Abs("../../models/")
	/*
		log.Print(p)
		s, err := hyde.New(p)
		if err != nil {
			log.Fatal(err)
		}
		go s.Start()
		s.WaitStop()
	*/

	hull := hyde.NewHull(":3000", p, agents, data, http, server, p1)
	hull.Start()
}
