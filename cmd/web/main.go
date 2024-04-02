package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rs/cors"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/web"
)

var autosavePath string
var addr string

func main() {
	flag.StringVar(&autosavePath, "autosave", "./autosave.halation-nm.json", "path to autosave JSON to")
	flag.StringVar(&addr, "addr", ":3939", "bind address")
	flag.Parse()

	s := web.NewServer(initShow())
	s.AutosavePath = autosavePath
	err := s.LoadAutosave()
	if err != nil {
		log.Printf("load autosave failed: %s", err)
	} else {
		log.Print("load autosave ok")
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://halation.nyiyui.ca", "http://halation.nyiyui.ca"},
	})

	http.ListenAndServe(addr, c.Handler(s))
}

func initShow() (*aiz.Runner, *node.NodeRunner, *node.Cuelist) {
	cuelist := node.NewCuelist()
	runner := aiz.NewRunner()
	var err error
	_ = err
	/*
		c := osc.NewDefaultClient()
		c.Register(runner)
		err = c.Blackout()
		if err != nil {
			panic(err)
		}
		log.Printf("osc setup ok")
	*/
	//mpvClient, err := mpv.NewClientUsingSubprocess()
	//if err != nil {
	//	panic(err)
	//}
	//mpvClient.Register(runner)

	nr := node.NewNodeRunner(runner)
	return runner, nr, cuelist
}
