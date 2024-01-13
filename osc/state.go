package osc

type State struct {
	// Blackout specifies whether all lights are off. If Blackout is true, everything else is ignored.
	Blackout bool
	Channels []Channel `json:"channels"`
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
