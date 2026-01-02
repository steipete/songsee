# songsee

Generate modern spectrogram images from audio files.

## Features

- Classic time–frequency spectrograms (FFT/STFT, Hann window)
- Native decoding for WAV + MP3, ffmpeg fallback for everything else
- PNG or JPEG output (default JPG)
- Time slicing via `--start` + `--duration`
- Palette styles: classic, magma, inferno, viridis, gray

## Install

```bash
go install github.com/steipete/songsee/cmd/songsee@latest
```

## Quick start

```bash
songsee track.mp3
songsee track.wav --style magma --width 2048 --height 1024 -o spectro.png
cat track.mp3 | songsee - --style gray --format png
songsee track.mp3 --start 12.5 --duration 8 --output slice.jpg
```

## Usage

```text
songsee [flags] <input>
```

Input can be a file path or `-` for stdin.

### Flags

- `-o, --output` output image path (default: input name + extension)
- `--format` `jpg` or `png` (default: `jpg`)
- `--width`, `--height` output size (default: 1920x1080)
- `--window` FFT window size, power of two (default: 2048)
- `--hop` hop size in samples (default: 512)
- `--min-freq`, `--max-freq` frequency range in Hz
- `--start`, `--duration` time slice in seconds
- `--style` palette name
- `--sample-rate` ffmpeg output sample rate
- `--ffmpeg` explicit ffmpeg path
- `-q, --quiet` suppress stdout output
- `-v, --verbose` verbose stderr output
- `--version` print version

## Output notes

- Default output is `input.jpg` or `input.png` based on `--format`.
- When `--output` has no extension, the chosen format is appended.
- JPEG quality is set to 95.

## Decoding

- WAV/MP3 decode is native.
- For other formats, install `ffmpeg` and songsee will auto-fallback.

## Development

```bash
go test ./... -cover

golangci-lint run
```

Lint config: `.golangci.yml`.

## License

MIT. See `LICENSE`.
