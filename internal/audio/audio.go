// Package audio handles decoding audio into mono float samples.
package audio

import "fmt"

// Audio holds mono samples in [-1,1] range.
type Audio struct {
	SampleRate int
	Samples    []float64
}

// Options controls decoding behavior.
type Options struct {
	SampleRate int
	FFmpegPath string
}

var (
	// ErrUnsupported is returned when no decoder can handle the input.
	ErrUnsupported = fmt.Errorf("unsupported audio format")
)
