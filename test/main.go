package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/elos/hyde"
)

var (
	watch = flag.Bool("w", true, "Watch files for changes")
)

func main() {
	p, _ := filepath.Abs("../../../../../../github.com/elos/documentation")
	log.Print(p)
	s, err := hyde.New(p)
	if err != nil {
		log.Fatal(err)
	}
	go s.Start()
	s.WaitStop()
}
