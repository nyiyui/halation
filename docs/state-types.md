# State Types

## `nyiyui.ca/halation/osc`

Controls the ColorSource AV lightboard via OSC commands.

Blackout:
```json
{"blackout": true}
```

Channels:
```json
{"channels":[
  {"channelID": 1, "level": 0", "hue": 0, "saturation": 0},
  {"channelID": 40, "level": 100", "hue": 359, "saturation": 100}
]}
```

## `nyiyui.ca/halation/mpv`

Controls an MPV instance using its socket-based IPC system.

**Warning: subject to change**

```json
{
  "extraProperties": {
    "filename": "./Big Buck Bunny.mov",
    "paused": true,
    "mute": true
  }
}
```
