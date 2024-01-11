import requests
from typing import List
from urllib.parse import urljoin


base_url = "http://10.10.0.2:8080"

def send_get(command, base_url=base_url):
    url = urljoin(base_url, command)
    requests.get(url)

def cs_chan_select(x: int):
    "Exclusively select channel x"
    send_get(f'cs/chan/select/{x}')

def cs_chan_add(channel: int):
    send_get(f'cs/chan/add/{channel}')

def cs_chan_subtract(channel: int):
    send_get(f'cs/chan/subtract/{channel}')

def cs_chan_thru(x: int):
    "Select every channel from the last channel selected 'thru' to channel x"
    send_get(f'cs/chan/thru/{x}')

def cs_chan_set(channels: List[int]):
    "Select channels listed and them only."
    # TODO: use cs_chan_thru to optimise requests made
    if len(channels) == 0:
        cs_chan_select(1)
        cs_chan_subtract(1)
    else:
        cs_chan_select(channels[0])
        for i in range(1, len(channels)):
            cs_chan_add(channels[i])

def cs_chan_at(level: int):
    "Set the selected channel's level."
    send_get(f'cs/chan/at/{level}')

def cs_color_hs(hue: int, saturation: int):
    "Set the selected channel's color. hue is 0-360, saturation is 0-100."
    send_get(f'cs/color/hs/{hue}/{saturation}')
