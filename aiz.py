import lib
from typing import List, Optional
from dataclasses import dataclass


@dataclass
class Cue:
    channels: List["Channel"]

    def apply(self):
        for channel in self.channels:
            channel.apply()

@dataclass
class Channel:
    channel: int
    level: int
    color: Optional["Color"] = None

    def apply(self):
        print(f'apply {self}')
        lib.cs_chan_select(self.channel)
        lib.cs_chan_at(self.level)
        print(f'done {self}')
        if self.color is not None:
            lib.cs_color_hs(self.color.hue, self.color.saturation)

@dataclass
class Color:
    hue: int # 0-360
    saturation: int # 0-100
    def __post_init__(self):
        assert  360 <= self.hue >= 0, "Hue must be between 0-360"
        assert 100 <= self.saturation >= 0, "Saturation must be between 0-100"
        
class Blackout:
    def apply(self):
        lib.cs_chan_select(1)
        lib.cs_chan_thru(40)
        lib.cs_chan_at(0)
