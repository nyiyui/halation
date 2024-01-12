import cues
from cues import blackout
import aiz
from aiz import Cue, Channel, Color
import channels


cuelist = {
        "0 Cues": blackout,
        "1 MCs": cues.stage_right_podium,
        "2 Love on Top": Cue(channels=[
            Channel(channel=channels.lx4_pink_purple, level=100),
            Channel(channel=channels.lx4_yellow, level=100),
        ]),
        "3 Upsweat": Cue(channels=[
            Channel(channel=channels.lx4_white, level=50),
            Channel(channel=channels.lx4_pink_purple, level=100),
        ]), # Sabrina says ok?????
        "4 No ce ce": Cue(channels=[
            Channel(channel=channels.lx4_blue, level=100),
        ]),
        "5 Bang Bang": Cue(channels=[
            Channel(channel=channels.lx4_pink_purple, level=100),
            Channel(channel=channels.lx4_white, level=30),
            Channel(channel=channels.lx4_yellow, level=50),
        ]),
        "6 Greased Lightning": Cue(channels=[
            Channel(channel=channels.lx4_blue, level=100),
            Channel(channel=channels.lx4_red, level=70),
        ]),
        "7 I'm Still Standing": Cue(channels=[
        ]),
        "8 Feeling Good": Cue(channels=[
            Channel(channel=channels.lx4_pink_purple, level=100),
            Channel(channel=channels.lx4_white, level=30),
            Channel(channel=channels.lx4_yellow, level=50),
        ]),
        "9 AfroFun": Cue(channels=[
            Channel(channel=channels.lx4_white, level=30),
            Channel(channel=channels.lx4_red, level=70),
            Channel(channel=channels.lx4_yellow, level=50),
        ]),
        "10 Pinga Ga Pori": Cue(channels=[
            Channel(channel=channels.lx4_white, level=30),
            Channel(channel=channels.lx4_yellow, level=50),
        ]),
        "11 Afro-carribe an dance mix": Cue(channels=[
            Channel(channel=channels.lx4_white, level=30),
            Channel(channel=channels.lx4_red, level=70),
            Channel(channel=channels.lx4_yellow, level=50),
        ]),
        "12 MCs": cues.stage_right_podium,
}

def testrun():
    import time
    for key, cue in cuelist.items():
        print(key)
        cue.apply()
        time.sleep(1)
        cues.blackout.apply()
        time.sleep(0.5)

if __name__ == '__main__':
    testrun()
