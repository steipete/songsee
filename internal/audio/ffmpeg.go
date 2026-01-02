// Package audio handles decoding audio into mono float samples.
package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os/exec"
)

// DecodeWithFFmpeg uses ffmpeg to decode any input into mono float samples.
func DecodeWithFFmpeg(path string, stdin io.Reader, sampleRate int, ffmpegPath string) (Audio, error) {
	if sampleRate <= 0 {
		sampleRate = 44100
	}
	ffmpeg, err := resolveFFmpeg(ffmpegPath)
	if err != nil {
		return Audio{}, err
	}

	args := []string{"-hide_banner", "-loglevel", "error"}
	if stdin != nil {
		args = append(args, "-i", "pipe:0")
	} else {
		args = append(args, "-i", path)
	}
	args = append(args, "-f", "f32le", "-ac", "1", "-ar", fmt.Sprintf("%d", sampleRate), "-")

	cmd := exec.Command(ffmpeg, args...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return Audio{}, fmt.Errorf("ffmpeg: %v: %s", err, stderr.String())
		}
		return Audio{}, err
	}

	if len(out)%4 != 0 {
		return Audio{}, fmt.Errorf("ffmpeg: unexpected pcm length")
	}

	samples := make([]float64, len(out)/4)
	for i := 0; i < len(samples); i++ {
		bits := binary.LittleEndian.Uint32(out[i*4 : i*4+4])
		samples[i] = float64(math.Float32frombits(bits))
	}

	return Audio{SampleRate: sampleRate, Samples: samples}, nil
}

func resolveFFmpeg(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH")
	}
	return ffmpeg, nil
}
