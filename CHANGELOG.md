# Changelog

## 0.2.0 - 2026-01-02

- Multi-panel visualization grid via `--viz`
- Added feature views: mel, chroma, hpss, selfsim, loudness, tempogram, mfcc
- New heatmap renderer + loudness panel
- Updated docs/spec to cover visualizations

## 0.1.0 - 2026-01-02

- Generate classic spectrogram images from audio files
- Native WAV/MP3 decoding with ffmpeg fallback
- PNG/JPEG output with configurable size
- Time slicing via `--start` and `--duration`
- FFT window + hop size controls
- Frequency range selection (`--min-freq` / `--max-freq`)
- Multiple palette styles (classic, magma, inferno, viridis, gray)
- CLI quality-of-life flags (`--quiet`, `--verbose`, `--version`)
- End-to-end MP3 test with image validation
