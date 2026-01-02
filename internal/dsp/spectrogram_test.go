package dsp

import "testing"

func TestComputeSpectrogram(t *testing.T) {
	samples := make([]float64, 4096)
	for i := range samples {
		samples[i] = 0.5
	}
	spec := ComputeSpectrogram(samples, 44100, 1024, 256)
	if spec.Frames <= 0 || spec.Bins <= 0 {
		t.Fatalf("invalid spec size")
	}
	if len(spec.Values) != spec.Frames*spec.Bins {
		t.Fatalf("values len mismatch")
	}
	if spec.Min >= spec.Max {
		t.Fatalf("min/max not set")
	}
	if spec.BinHz <= 0 {
		t.Fatalf("invalid bin hz")
	}
}

func TestHannWindow(t *testing.T) {
	w1 := HannWindow(1)
	if len(w1) != 1 || w1[0] != 1 {
		t.Fatalf("hann size 1")
	}
	w := HannWindow(4)
	if len(w) != 4 {
		t.Fatalf("hann size 4")
	}
	if w[0] != 0 || w[3] != 0 {
		t.Fatalf("hann endpoints")
	}
}

func TestComputeSpectrogramDefaults(t *testing.T) {
	samples := make([]float64, 100)
	spec := ComputeSpectrogram(samples, 0, 0, 0)
	if spec.SampleRate != 44100 {
		t.Fatalf("sample rate default = %d", spec.SampleRate)
	}
	if spec.WindowSize != 2048 {
		t.Fatalf("window default = %d", spec.WindowSize)
	}
	if spec.Frames != 1 {
		t.Fatalf("frames = %d", spec.Frames)
	}
}
