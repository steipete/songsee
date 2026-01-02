// Package audio handles decoding audio into mono float samples.
package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DecodeFile reads an audio file, decoding WAV/MP3 and falling back to ffmpeg.
func DecodeFile(path string, opts Options) (Audio, error) {
	file, err := os.Open(path)
	if err != nil {
		return Audio{}, err
	}
	defer func() { _ = file.Close() }()

	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".wav" || ext == ".wave" {
		if pcm, ok, err := DecodeWAVIf(file); ok {
			return pcm, err
		}
	}
	if ext == ".mp3" {
		if pcm, ok, err := DecodeMP3If(file); ok {
			return pcm, err
		}
	}

	if pcm, ok, err := DecodeWAVIf(file); ok {
		return pcm, err
	}
	if pcm, ok, err := DecodeMP3If(file); ok {
		return pcm, err
	}

	if opts.SampleRate == 0 {
		opts.SampleRate = 44100
	}
	pcm, err := DecodeWithFFmpeg(path, nil, opts.SampleRate, opts.FFmpegPath)
	if err != nil {
		return Audio{}, fmt.Errorf("%w; ffmpeg fallback failed: %v", ErrUnsupported, err)
	}
	return pcm, nil
}

// DecodeBytes decodes audio data from a byte slice.
func DecodeBytes(data []byte, opts Options) (Audio, error) {
	reader := bytes.NewReader(data)
	if pcm, ok, err := DecodeWAVIf(reader); ok {
		return pcm, err
	}
	reader.Reset(data)
	if pcm, ok, err := DecodeMP3If(reader); ok {
		return pcm, err
	}

	if opts.SampleRate == 0 {
		opts.SampleRate = 44100
	}
	pcm, err := DecodeWithFFmpeg("", bytes.NewReader(data), opts.SampleRate, opts.FFmpegPath)
	if err != nil {
		return Audio{}, fmt.Errorf("%w; ffmpeg fallback failed: %v", ErrUnsupported, err)
	}
	return pcm, nil
}

// DecodeReader decodes audio data from an io.Reader.
func DecodeReader(r io.Reader, opts Options) (Audio, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return Audio{}, err
	}
	return DecodeBytes(data, opts)
}
