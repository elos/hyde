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
	p, _ := filepath.Abs("./")
	s, _ := hyde.New(p)
	go s.Start()
	s.WaitStop()
}
