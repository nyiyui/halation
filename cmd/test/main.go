package main

import (
	"context"
	"log"
	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/mpv"
	"time"
)

func main() {
	r := &aiz.Runner{Specific: map[string]interface{}{}}
	mpvClient, err := mpv.NewClientUsingSubprocess()
	if err != nil {
		panic(err)
	}
	mpvClient.Register(r)

	show := &aiz.Show{Cues: []aiz.Cue{
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
	}}
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
}
