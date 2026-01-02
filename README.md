# songsee 🎵👁️

**See your sound.** Generate gorgeous spectrogram visualizations from any audio file.

![9-mode visualization example](example.png)

*9 visualization modes in one image. Yes, we went there.* 🦞

## Features

- 🎨 **9 visualization modes**: spectrogram, mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux
- 🌈 **6 color palettes**: classic, magma, inferno, viridis, gray, and *clawd* (the lobster special 🦞)
- 📊 **Combine modes**: Stack multiple visualizations in one image
- 🎵 **Universal input**: WAV, MP3, or anything ffmpeg can handle
- ⚡ **Fast**: Native Go, no Python dependencies
- 🖼️ **Flexible output**: PNG or JPEG, customizable dimensions

## Install

```bash
go install github.com/steipete/songsee/cmd/songsee@latest
```

Or clone and build:
```bash
git clone https://github.com/steipete/songsee.git
cd songsee
go build -o songsee ./cmd/songsee/
```

## Quick Start

```bash
# Basic spectrogram
songsee track.mp3

# Mel spectrogram with magma palette
songsee track.mp3 --viz mel --style magma

# THE MULTIPASS: all 9 modes combined 🎫
songsee track.mp3 --viz spectrogram,mel,chroma,hpss,selfsim,loudness,tempogram,mfcc,flux --style clawd

# Custom output
songsee track.mp3 --viz hpss,chroma --style inferno -o my_viz.png --width 2560 --height 1440
```

## Visualization Modes

| Mode | What it shows |
|------|---------------|
| `spectrogram` | Classic time × frequency magnitude |
| `mel` | Perceptual frequency scale (how humans hear) |
| `chroma` | 12-bin pitch class (for harmony/key analysis) |
| `hpss` | Harmonic vs percussive separation |
| `selfsim` | Self-similarity matrix (song structure) |
| `loudness` | Volume/energy over time |
| `tempogram` | Tempo variation over time |
| `mfcc` | Timbre fingerprint |
| `flux` | Spectral change (where the action is ⚡) |

## Palettes

- `classic` - Traditional spectrogram colors
- `magma` - Hot gradient 🔥
- `inferno` - Even hotter 🌋
- `viridis` - Scientific green-yellow
- `gray` - Noir mode 🎬
- `clawd` - Ocean depths to lobster red 🦞

## Advanced Options

```
--width         Output width in pixels (default: 1920)
--height        Output height in pixels (default: 1080)
--window        FFT window size (default: 2048)
--hop           Hop size in samples (default: 512)
--start         Start time in seconds
--duration      Duration in seconds
--min-freq      Minimum frequency in Hz
--max-freq      Maximum frequency in Hz
--format        Output format: jpg or png (default: jpg)
```

## Why "songsee"?

Because `spectrogram-generator-cli-tool-v2-final-FINAL` was taken.

Also: **see** your **song**. Get it? 👀🎵

---

Built at inference speed by [@steipete](https://twitter.com/steipete) with help from [Clawd](https://clawd.ai) 🦞

*"EXFOLIATE your audio!"*
