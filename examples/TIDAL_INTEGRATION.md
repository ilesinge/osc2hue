# Tidal Cycles + OSC2Hue Integration Guide

This guide shows how to control Philips Hue lights from Tidal Cycles using OSC2Hue as a bridge.

## Setup

### 1. Prerequisites

- **OSC2Hue**: Running on your system (default: localhost:8080)
- **Tidal Cycles**: Installed and working
- **SuperCollider**: Running with Tidal
- **Philips Hue Bridge**: Connected and configured with OSC2Hue

### 2. Start OSC2Hue

```bash
# Make sure OSC2Hue is running
./osc2hue
# Should show: "OSC Server: 0.0.0.0:8080"
```

### 3. Configure Tidal Cycles

## Modify the BootTidal.hs file

Comment the following line by adding "--" at the start:

```haskell
-- tidal <- startTidal (superdirtTarget {oLatency = 0.05, oAddress = "127.0.0.1", oPort = 57120}) (defaultConfig {cVerbose = True, cFrameTimespan = 1/20})
```

Add this section in the same location (tweak the address, port, and latency if needed):

```haskell
:{
let hueTarget =
      Target {oName = "hue",          -- A friendly name for the target (only used in error messages)
              oAddress = "localhost", -- The target's network address, normally "localhost"
              oPort = 8080,           -- The network port the target is listening on
              oLatency = 0.2,         -- Additional delay, to smooth out network jitter/get things in sync
              oSchedule = Live,       -- The scheduling method
              oWindow = Nothing,      -- Not yet used
              oHandshake = False,     -- SuperDirt specific
              oBusPort = Nothing      -- Also SuperDirt specific
             }
    hueOSC = [OSC "/hue/{light}/set" $ ArgList [
                                                ("x", Just $ VF (-1)),
                                                ("y", Just $ VF (-1)),
                                                ("brightness", Just $ VF (-1)),
                                                ("duration", Just $ VF (-1))
                                   ],
              OSC "/hue/{light}/on" $ ArgList [
                                      ("on", Nothing),
                                      ("duration", Just $ VF (-1))
                                    ]
              ]
    brightness = pF "brightness"
    light = pS "light"
    x = pF "x"
    y = pF "y"
    duration = pF "duration"
    on = pI "on"
    oscmap = [(hueTarget, hueOSC), (superdirtTarget {oLatency = 0.05}, [superdirtShape])]
:}

tidal <- startStream defaultConfig oscmap
```

## Usage

```haskell
-- Comprehensive example
d1 $ light "all*4" # brightness (sine) # duration (range 1 200 (slow 2 sine)) # x saw # y sine -- send 4 update events to all lights with variating brightness, duration and color (x and y)

-- Basic light control patterns

  -- Turn lights on/off in a pattern
d1 $ light "3" # on "<1 0>" # duration 1 -- toggle light #3 between on and off states
  
  -- Control brightness with a sine wave
d1 $ light "3*8" # brightness (slow 4 sine) # duration 200
  
  -- Change colors using XY coordinates
d2 $ light "3*4" # x (range 0 1 $ slow 8 sine) # y (range 0 1 $ slow 6 cosine) # duration 400

-- Advanced patterns

  -- Color cycling through different hues
d3 $ light "all" # x (choose [0.1, 0.3, 0.5, 0.7]) # y (choose [0.1, 0.4, 0.6, 0.8]) # duration 1000

-- Euclidean rhythms for light patterns

  -- Euclidean on/off pattern
d4 $ light "3*10" # euclidFull 3 5 (on "1") (on "0") # duration 10
  
  -- Brightness euclidean with different divisions
d5 $ light "3*16" # euclidFull 5 16 (brightness "0.25") (brightness "0.75") # duration 10
  
  -- Color changes on euclidean hits
d5 $ light "3*16" # euclidFull 5 8 (x "0") (x "0.6") # y "0.3" # duration 20

-- Silence all patterns
hush

-- Custom functions for common patterns
let 
  -- Smooth color transition
  wave lightid dur = duration "300*16" # light lightid # x (slow dur sine) # y (slow (dur * 0.7) cosine)
  
  -- Strobe effect
  strobe lightid speed = brightness (fast speed $ "1 0") # light lightid # duration 10
  
  -- Breathing effect
  breathe lightid speed = duration "500*8" # light lightid # brightness (slow speed $ range 0.2 1 sine) # duration 500

-- Emotion-driven patterns
let
  -- Joy: Bright, warm yellow/orange with fast brightness pulses for energy
  joy = light "all*8" # x 0.4 # y 0.5 # brightness (fast 2 $ range 0.7 1 sine) # duration 200

  -- Calm: Soft green with slow, gentle brightness waves for tranquility
  calm = light "all*8" # x 0.3 # y 0.3 # brightness (slow 8 $ range 0.2 0.6 sine) # duration 2000

  -- Excitement: Fast cycling through lights with random bright colors and quick transitions
  excitement = light (fast 8 $ "<1 2 3 4 5 6>") # x (fast 4 $ choose [0.1,0.7]) # y (fast 4 $ choose [0.1,0.7]) # brightness 1 # duration 10

  -- Anger: Sharp, fast, red pulses
  anger = light "all*16" # x 0.7 # y 0.3 # brightness (fast 8 $ "1 0.1") # duration 25
  
  -- Sadness: Slow, dim, blue fade
  sadness = light "all*16" # x 0.15 # y 0.06 # brightness (slow 8 $ range 0.05 0.3 sine) # duration 100

  -- Fear: Erratic, cold colors, flickering
  fear = struct "t*16" $ light (fast 32 $ choose ["1","2","3","all"]) # x 0.15 # y 0.15 # brightness (fast 16 $ choose [0, 0.1, 0.8]) # duration 10

   -- Disgust: Sickly green, uneven patterns
  disgust = light "3*16" # x 0.25 # y 0.7 # brightness (slow 3 $ range 0.2 0.6 $ irand 10) # duration 800


-- Use custom functions

d1 $ hueColorWave 1 8
d2 $ hueStrobe 2 4
d3 $ hueBreathe 3 6
d4 $ sadness
```

## Performance Tips

- **Don't flood**: Avoid patterns faster than ~10Hz to prevent overwhelming the bridge


