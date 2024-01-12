from aiz import Cue, Channel, Color
import time
import cues



def test_channel(c):
    Cue(channels=[
        Channel(channel=c, level=100),
    ]).apply()
    time.sleep(1)
    cues.blackout.apply()


if __name__ == '__main__':
    Cue(channels=[
        Channel(channel=31, level=100),
    ]).apply()
    #cues.stage_right_podium.apply()
    time.sleep(1)
    cues.blackout.apply()
