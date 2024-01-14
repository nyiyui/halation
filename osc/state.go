package osc

import (
	"errors"
	"fmt"
	"nyiyui.ca/halation/aiz"
)

type State struct {
	// Blackout specifies whether all lights are off. If Blackout is true, everything else is ignored.
	Blackout bool
	Channels []Channel `json:"channels"`
}

var _ aiz.State = new(State)

func (s *State) Reify(r *aiz.Runner, g aiz.Gradient, prev_ aiz.State) error {
	c_, ok := r.Specific["nyiyui.ca/halation/osc"]
	if !ok {
		return errors.New("runner doesn't have OSC client")
	}
	c, ok := c_.(*Client)
	if !ok {
		return errors.New("runner has wrong client")
	}
	prev, ok := prev_.(*State)
	if !ok || prev.Blackout {
		prev = &State{Channels: nil}
	}
	transitions := map[int][2]Channel{}
	for _, ch := range prev.Channels {
		transitions[ch.ChannelID] = [2]Channel{ch, Channel{}}
	}
	for _, ch := range s.Channels {
		trans, ok := transitions[ch.ChannelID]
		if !ok {
			transitions[ch.ChannelID] = [2]Channel{Channel{}, ch}
		} else {
			transitions[ch.ChannelID] = [2]Channel{trans[0], ch}
		}
	}
	for i, val := range g.Values(g.PreferredResolution()) {
		err := s.applyStep(c, transitions, val)
		if err != nil {
			return fmt.Errorf("step %d (t=%f): %w", i, val, err)
		}
	}
	return nil
}

func (s *State) applyStep(c *Client, transitions map[int][2]Channel, val float32) error {
	for cid, chs := range transitions {
		prev, next := chs[0], chs[1]
		err := c.ChanSelect(cid)
		if err != nil {
			return err
		}
		{ // level
			delta := next.Level - prev.Level
			level := int(float32(delta) * val)
			err := c.ChanAt(level)
			if err != nil {
				return err
			}
		}
		if !(prev.Saturation == 0 && next.Saturation == 0) { // hue, saturation
			deltaH := next.Hue - prev.Hue
			deltaS := next.Saturation - prev.Saturation
			hue := int(float32(deltaH) * val)
			saturation := int(float32(deltaS) * val)
			err := c.ColorHS(hue, saturation)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Channel struct {
	ChannelID  int `json:"channelID"`
	Level      int `json:"level"`
	Hue        int `json:"hue"`
	Saturation int `json:"saturation"`
}

func (c *Client) ApplyState(s State) (err error) {
	if s.Blackout {
		err = c.ChanSelect(1)
		if err != nil {
			return
		}
		err = c.ChanThru(NumberOfChannels)
		if err != nil {
			return
		}
		err = c.ChanAt(0)
		return
	}
	for _, ch := range s.Channels {
		err = c.ChanSelect(ch.ChannelID)
		if err != nil {
			return
		}
		err = c.ChanAt(ch.Level)
		if err != nil {
			return
		}
		if ch.Saturation != 0 {
			err = c.ColorHS(ch.Hue, ch.Saturation)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *State) TypeName() string { return "nyiyui.ca/halation/osc" }
