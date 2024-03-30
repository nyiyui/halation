package main

import (
	"flag"
	"log"
	"net/http"

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
	http.ListenAndServe(addr, s)
}

func initShow() (*aiz.Runner, *node.NodeRunner, *node.Cuelist) {
	cuelist := node.NewCuelist()
	runner := aiz.NewRunner()
	/*
		  var err error
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
	//n00 := node.NodeName{"nyiyui.ca/halation/cmd/web", "00"}
	//n01 := node.NodeName{"nyiyui.ca/halation/cmd/web", "01"}
	//n02 := node.NodeName{"nyiyui.ca/halation/cmd/web", "02"}
	//n03 := node.NodeName{"nyiyui.ca/halation/cmd/web", "03"}
	////n04 := node.NodeName{"nyiyui.ca/halation/cmd/web", "04"}
	//nr.NM.Nodes[n00] = node.NewManual()
	//cuelist.Nodes[0] = n00
	//nr.NM.Nodes[n00].SetDescription("Pre-show")

	//nr.NM.Nodes[n01] = node.NewSetState(&aiz.SG{State: &osc.State{
	//	Channels: []osc.Channel{
	//		{ChannelID: osc.ChannelLeftCentreWall, Level: 100, Hue: 0, Saturation: 0},
	//	},
	//}, Gradient: &gradient.LinearGradient{
	//	Duration_:            timeutil.Duration(3 * time.Second),
	//	PreferredResolution_: timeutil.Duration(100 * time.Millisecond),
	//}})
	//nr.NM.Nodes[n01].BaseNodeRef().Promises = []node.Promise{
	//	{"dummy", n00},
	//}
	//nr.NM.Nodes[n01].SetDescription("wall")

	//nr.NM.Nodes[n02] = node.NewEvalLua(`print("hello from lua")`)
	//nr.NM.Nodes[n02].BaseNodeRef().Promises = []node.Promise{
	//	{"dummy", n01},
	//}
	//nr.NM.Nodes[n02].SetDescription("Lua")
	//cuelist.Nodes[1.0] = n03
	//nr.NM.Nodes[n03] = node.NewManual()
	//nr.NM.Nodes[n03].SetDescription("Emcees")
	return runner, nr, cuelist
}
