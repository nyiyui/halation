package main

import (
	"net/http"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/trigger"
	"nyiyui.ca/halation/web"
)

func main() {
	s := web.NewServer(initShow())
	http.ListenAndServe(":8080", s)
}

func initShow() (*aiz.Runner, *aiz.Show) {
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
					Duration_:            1 * time.Second,
					PreferredResolution_: 50 * time.Millisecond,
				}},
			}},
		},
		Triggers: []aiz.Trigger{
			&trigger.Timed{Delay: 0, CueRequest: aiz.CueRequest{0}},
			&trigger.Timed{Delay: 3 * time.Second, CueRequest: aiz.CueRequest{1}},
			&trigger.Timed{Delay: 6 * time.Second, CueRequest: aiz.CueRequest{2}},
		},
	}
	show.SetupTriggers(runner)
	return runner, show
}
