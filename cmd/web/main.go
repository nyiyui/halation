package main

import (
	"log"
	"net/http"
	"nyiyui.ca/halation/timeutil"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/web"
)

func main() {
	s := web.NewServer(initShow())
	http.ListenAndServe(":8080", s)
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
	*/
	log.Printf("osc setup ok")
	//mpvClient, err := mpv.NewClientUsingSubprocess()
	//if err != nil {
	//	panic(err)
	//}
	//mpvClient.Register(runner)

	nr := node.NewNodeRunner(runner)
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000"}] = node.NewManual()
	cuelist.Nodes[0] = node.NodeName{"nyiyui.ca/halation/cmd/web", "000"}
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000"}].SetDescription("Pre-show")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}] = node.NewSetState(&aiz.SG{State: &osc.State{
		Channels: []osc.Channel{
			{ChannelID: osc.ChannelLeftCentreWall, Level: 100, Hue: 0, Saturation: 0},
		},
	}, Gradient: &gradient.LinearGradient{
		Duration_:            timeutil.Duration(3 * time.Second),
		PreferredResolution_: timeutil.Duration(100 * time.Millisecond),
	}})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}].SetListensTo([]node.NodeName{
		node.NodeName{"nyiyui.ca/halation/cmd/web", "000"},
	})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}].SetDescription("wall")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}] = node.NewSetState(&aiz.SG{State: &osc.State{
		Blackout: true,
	}, Gradient: &gradient.LinearGradient{
		Duration_:            timeutil.Duration(3 * time.Second),
		PreferredResolution_: timeutil.Duration(100 * time.Millisecond),
	}})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}].SetListensTo([]node.NodeName{
		node.NodeName{"nyiyui.ca/halation/cmd/web", "001"},
	})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}].SetDescription("blackout")
	cuelist.Nodes[1.0] = node.NodeName{"nyiyui.ca/halation/cmd/web", "001"}
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "001"}] = node.NewManual()
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "001"}].SetDescription("Emcees")
	{
		m := node.NewManual()
		m.SetDescription("Idol")
		nr.NM.Nodes[node.NodeName{"", "002"}] = m
		t := node.NewTimer(1 * time.Second)
		t.SetListensTo([]node.NodeName{
			node.NodeName{"", "002"},
		})
		nr.NM.Nodes[node.NodeName{"", "002-t"}] = t
		nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}].SetListensTo([]node.NodeName{
			node.NodeName{"nyiyui.ca/halation/cmd/web", "001"},
			node.NodeName{"", "002-t"},
		})
	}
	return runner, nr, cuelist
}
