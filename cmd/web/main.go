package main

import (
	"net/http"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/timeutil"
	"nyiyui.ca/halation/web"
)

func main() {
	s := web.NewServer(initShow())
	http.ListenAndServe(":8080", s)
}

func initShow() (*aiz.Runner, *node.NodeRunner) {
	runner := &aiz.Runner{Specific: map[string]interface{}{}}
	runner.Setup()
	//c := osc.NewDefaultClient()
	//c.Register(runner)
	//err = c.Blackout()
	//if err != nil {
	//	panic(err)
	//}
	//mpvClient, err := mpv.NewClientUsingSubprocess()
	//if err != nil {
	//	panic(err)
	//}
	//mpvClient.Register(runner)

	show := &aiz.Show{
		Cues: []aiz.Cue{
			{Name: "0 blackout", SGs: []aiz.SG{
				{State: &osc.State{
					Blackout: true,
				}},
			}},
			{Name: "1 test", SGs: []aiz.SG{
				{State: &osc.State{
					Channels: []osc.Channel{
						{ChannelID: osc.ChannelLx4Multi, Level: 100},
						{ChannelID: osc.ChannelPotA, Level: 100, Hue: 0, Saturation: 100},
					},
				}},
			}},
			{Name: "2 test2", SGs: []aiz.SG{
				{State: &osc.State{
					Channels: []osc.Channel{
						{ChannelID: osc.ChannelLx4Multi, Level: 0},
						{ChannelID: osc.ChannelLx4Red, Level: 100},
						{ChannelID: osc.ChannelPotA, Level: 0, Hue: 0, Saturation: 100},
					},
				}, Gradient: &gradient.LinearGradient{
					Duration_:            timeutil.Duration(1 * time.Second),
					PreferredResolution_: timeutil.Duration(50 * time.Millisecond),
				}},
			}},
		},
	}
	_ = show
	nr := node.NewNodeRunner(runner)
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000"}] = node.NewManual()
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000"}].SetDescription("0 Pre-show")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}] = node.NewSetState(&aiz.SG{State: &osc.State{
		Channels: []osc.Channel{
			{ChannelID: osc.ChannelLx4Multi, Level: 0},
			{ChannelID: osc.ChannelLx4Red, Level: 100},
			{ChannelID: osc.ChannelPotA, Level: 0, Hue: 0, Saturation: 100},
		},
	}})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}].SetListensTo([]node.NodeName{
		node.NodeName{"nyiyui.ca/halation/cmd/web", "000"},
	})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-state"}].SetDescription("0 Red")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}] = node.NewSetState(&aiz.SG{State: &osc.State{
		Channels: []osc.Channel{
			{ChannelID: osc.ChannelLx4Multi, Level: 0},
			{ChannelID: osc.ChannelLx4Red, Level: 100},
			{ChannelID: osc.ChannelPotA, Level: 0, Hue: 0, Saturation: 100},
		},
	}})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}].SetListensTo([]node.NodeName{
		node.NodeName{"nyiyui.ca/halation/cmd/web", "000"},
		node.NodeName{"nyiyui.ca/halation/cmd/web", "001"},
		node.NodeName{"nyiyui.ca/halation/cmd/web", "002"},
	})
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "000-mpv"}].SetDescription("0 Video")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "001"}] = node.NewManual()
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "001"}].SetDescription("1 Emcees")
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "002"}] = node.NewManual()
	nr.NM.Nodes[node.NodeName{"nyiyui.ca/halation/cmd/web", "002"}].SetDescription("2 IDK")
	return runner, nr
}
