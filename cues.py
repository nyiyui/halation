import aiz
from aiz import Cue, Channel, Color
import channels


stage_right_podium = Cue(channels=[
    Channel(channel=channels.left_flood, level=20),
    Channel(channel=channels.left_podium, level=80),
    Channel(channel=channels.left_front, level=80),
])

blackout = aiz.Blackout()
