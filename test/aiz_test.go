package test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	. "nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/osc"
	"testing"
)

func TestJSON(t *testing.T) {
	s := &Show{
		Cues: []Cue{
			{
				SGs: []SG{
					{&osc.State{
						Blackout: true,
					}, &gradient.LinearGradient{
						Duration_:            1000,
						PreferredResolution_: 100,
					}},
					{&osc.State{
						Channels: []osc.Channel{
							{1, 100, 0, 0},
						},
					}, &gradient.LinearGradient{
						Duration_:            1000,
						PreferredResolution_: 100,
					}},
				},
			},
		},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}
	t.Logf("%s", data)
	var s2 Show
	err = json.Unmarshal(data, &s2)
	if err != nil {
		t.Fatalf("unmarshal: %s", err)
	}
	if !cmp.Equal(*s, s2) {
		t.Log(cmp.Diff(*s, s2))
		t.Fatal("not equal")
	}
}
