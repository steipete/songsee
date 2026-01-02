# 🌊 songsee — FFT so pretty, your ears will be jealous.

![9-mode visualization example](example.png)

## Features

- **9 visualization modes**: spectrogram, mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux
- **6 color palettes**: classic, magma, inferno, viridis, gray, clawd
- **Auto-contrast**: per-panel percentile normalization for readable heatmaps
- **Combine modes**: stack multiple visualizations in one grid image
- **Universal input**: WAV, MP3, or anything ffmpeg can handle
- **Fast**: native Go, no Python dependencies
- **Flexible output**: PNG or JPEG, customizable dimensions

## Install

```bash
brew install steipete/tap/songsee
```

```bash
go install github.com/steipete/songsee/cmd/songsee@latest
```

## Quick Start

```bash
# Basic spectrogram
songsee track.mp3

# Mel spectrogram with magma palette
songsee track.mp3 --viz mel --style magma

# All 9 modes combined
songsee track.mp3 --viz spectrogram,mel,chroma,hpss,selfsim,loudness,tempogram,mfcc,flux

# Custom output
songsee track.mp3 --viz hpss,chroma --style inferno -o viz.png --width 2560 --height 1440
```

## Visualization Modes

| Mode | Description |
|------|-------------|
| `spectrogram` | Time × frequency magnitude |
| `mel` | Perceptual frequency scale |
| `chroma` | 12-bin pitch class |
| `hpss` | Harmonic vs percussive separation |
| `selfsim` | Self-similarity matrix |
| `loudness` | Volume over time |
| `tempogram` | Tempo variation |
| `mfcc` | Timbre fingerprint |
| `flux` | Spectral change detection |

## Palettes

`classic` · `magma` · `inferno` · `viridis` · `gray` · `clawd` 🦞

## Options

```
--output        Output path (default: input name + extension)
--format        jpg or png (default: jpg)
--width         Output width (default: 1920)
--height        Output height (default: 1080)
--window        FFT window size (default: 2048)
--hop           Hop size (default: 512)
--min-freq      Minimum frequency in Hz
--max-freq      Maximum frequency in Hz
--start         Start time in seconds
--duration      Duration in seconds
--style         Palette name
--viz           Visualization list (repeatable or comma-separated)
```

---

Built by [@steipete](https://twitter.com/steipete)
