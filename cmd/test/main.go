package main

import (
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
	s := &mpv.State{
		FilePath:   "./big_buck_bunny_480p_h264.mov",
		Paused:     true,
		Position:   60,
		Fullscreen: false,
	}
	err = s.Reify(r, nil, nil)
	if err != nil {
		panic(err)
	}
	log.Printf("state 1 done")
	time.Sleep(3 * time.Second)
	s = &mpv.State{
		FilePath:   "./big_buck_bunny_480p_h264.mov",
		Paused:     false,
		Position:   0,
		Fullscreen: false,
	}
	err = s.Reify(r, nil, nil)
	if err != nil {
		panic(err)
	}
}
