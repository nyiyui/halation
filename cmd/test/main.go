package main

import (
	"context"
	"log"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/osc"
)

func main() {
	var err error
	r := &aiz.Runner{Specific: map[string]interface{}{}}
	c := osc.NewDefaultClient()
	c.Register(r)
	err = c.Blackout()
	if err != nil {
		panic(err)
	}
	//mpvClient, err := mpv.NewClientUsingSubprocess()
	//if err != nil {
	//	panic(err)
	//}
	//mpvClient.Register(r)

	show := &aiz.Show{Cues: []aiz.Cue{
		/*
			{Name: "0 paused", SGs: []aiz.SG{
				{State: &mpv.State{
					FilePath:   "./big_buck_bunny_480p_h264.mov",
					Paused:     mpv.Ptr(true),
					Position:   mpv.Ptr(0),
					Fullscreen: mpv.Ptr(false),
				}},
			}},
			{Name: "1 playing", SGs: []aiz.SG{
				{State: &mpv.State{
					FilePath:   "./big_buck_bunny_480p_h264.mov",
					Paused:     mpv.Ptr(false),
					Position:   mpv.Ptr(60),
					Fullscreen: mpv.Ptr(false),
				}},
			}},
		*/
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
				PreferredResolution_: 100 * time.Millisecond,
			}},
		}},
	}}
	log.Printf("cue 0 start")
	show.ApplyCue(r, 0, context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("cue 0 end")
	time.Sleep(1 * time.Second)
	log.Printf("cue 1 start")
	show.ApplyCue(r, 1, context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("cue 1 end")
	time.Sleep(1 * time.Second)
	log.Printf("cue 2 start")
	show.ApplyCue(r, 2, context.Background())
	if err != nil {
		panic(err)
	}
}
