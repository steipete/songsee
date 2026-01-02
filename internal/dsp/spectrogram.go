// Package dsp provides spectral analysis utilities.
package dsp

import (
	"math"
)

// Spectrogram contains log-magnitude FFT frames.
type Spectrogram struct {
	Frames     int
	Bins       int
	Values     []float64
	Min        float64
	Max        float64
	SampleRate int
	WindowSize int
	HopSize    int
	BinHz      float64
}

// HannWindow returns a Hann window of length n.
func HannWindow(n int) []float64 {
	w := make([]float64, n)
	if n == 1 {
		w[0] = 1
		return w
	}
	for i := 0; i < n; i++ {
		w[i] = 0.5 - 0.5*math.Cos(2*math.Pi*float64(i)/float64(n-1))
	}
	return w
}

// ComputeSpectrogram computes a log-magnitude spectrogram.
func ComputeSpectrogram(samples []float64, sampleRate, windowSize, hopSize int) Spectrogram {
	if windowSize <= 0 {
		windowSize = 2048
	}
	if hopSize <= 0 {
		hopSize = windowSize / 4
	}
	if hopSize <= 0 {
		hopSize = 1
	}
	if sampleRate <= 0 {
		sampleRate = 44100
	}

	frames := 1
	if len(samples) > windowSize {
		frames = 1 + (len(samples)-windowSize+hopSize-1)/hopSize
	}
	bins := windowSize/2 + 1
	values := make([]float64, frames*bins)

	window := HannWindow(windowSize)
	minVal := math.Inf(1)
	maxVal := math.Inf(-1)
	eps := 1e-9

	frame := make([]complex128, windowSize)
	for f := 0; f < frames; f++ {
		start := f * hopSize
		for i := 0; i < windowSize; i++ {
			idx := start + i
			if idx < len(samples) {
				frame[i] = complex(samples[idx]*window[i], 0)
			} else {
				frame[i] = 0
			}
		}
		FFTInPlace(frame)
		for b := 0; b < bins; b++ {
			re := real(frame[b])
			im := imag(frame[b])
			mag := math.Sqrt(re*re + im*im)
			db := 20 * math.Log10(mag+eps)
			values[f*bins+b] = db
			if db < minVal {
				minVal = db
			}
			if db > maxVal {
				maxVal = db
			}
		}
	}

	binHz := float64(sampleRate) / float64(windowSize)
	return Spectrogram{
		Frames:     frames,
		Bins:       bins,
		Values:     values,
		Min:        minVal,
		Max:        maxVal,
		SampleRate: sampleRate,
		WindowSize: windowSize,
		HopSize:    hopSize,
		BinHz:      binHz,
	}
}
