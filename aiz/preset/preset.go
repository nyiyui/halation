package preset

import "nyiyui.ca/halation/osc"

var StageRightPodium = osc.State{
	Channels: []osc.Channel{
		{ChannelID: osc.ChannelLeftFlood, Level: 20},
		{ChannelID: osc.ChannelLeftPodium, Level: 80},
		{ChannelID: osc.ChannelLeftFront, Level: 80},
	},
}
