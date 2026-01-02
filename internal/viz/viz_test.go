package viz

import (
	"image/color"
	"math"
	"testing"

	"github.com/steipete/songsee/internal/dsp"
)

func TestParseList(t *testing.T) {
	out, err := ParseList([]string{"spectrogram,mel", "chroma"})
	if err != nil {
		t.Fatalf("ParseList: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("unexpected list size")
	}
	defaults, err := ParseList(nil)
	if err != nil {
		t.Fatalf("ParseList default: %v", err)
	}
	if len(defaults) != 1 || defaults[0] != Spectrogram {
		t.Fatalf("unexpected default list")
	}
	if _, err := ParseList([]string{"nope"}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestRenderAllKinds(t *testing.T) {
	ctx := NewContext(testSamples(), 44100, 512, 128)
	opts := RenderOptions{
		Width:   120,
		Height:  80,
		Palette: colorRGBA,
	}
	kinds := []Kind{Spectrogram, Mel, Chroma, MFCC, HPSS, SelfSim, Loudness, Tempogram, Flux}
	for _, kind := range kinds {
		img, err := Render(kind, ctx, opts)
		if err != nil {
			t.Fatalf("Render %s: %v", kind, err)
		}
		if img.Bounds().Dx() != opts.Width || img.Bounds().Dy() != opts.Height {
			t.Fatalf("Render %s size mismatch", kind)
		}
	}
}

func TestKindsHelp(t *testing.T) {
	if KindsHelp() == "" {
		t.Fatalf("expected help text")
	}
}

func TestRenderUnknown(t *testing.T) {
	ctx := NewContext(testSamples(), 44100, 512, 128)
	_, err := Render(Kind("nope"), ctx, RenderOptions{Width: 10, Height: 10, Palette: colorRGBA})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPercentiles(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5}
	minVal, maxVal := percentileRange(values, 0.2, 0.8)
	if minVal <= 1 || maxVal >= 5 {
		t.Fatalf("unexpected percentile range")
	}
	val := percentileValue(values, 0.5)
	if val != 3 {
		t.Fatalf("unexpected percentile value")
	}
	if idx := percentileIndex(values, -1); idx != 0 {
		t.Fatalf("unexpected low index")
	}
	if idx := percentileIndex(values, 2); idx != len(values)-1 {
		t.Fatalf("unexpected high index")
	}
	minVal, maxVal = percentileRange(nil, 0.2, 0.8)
	if minVal != 0 || maxVal != 1 {
		t.Fatalf("unexpected empty percentile range")
	}
	if percentileValue(nil, 0.5) != 1 {
		t.Fatalf("unexpected empty percentile value")
	}
}

func TestSampleValues(t *testing.T) {
	values := []float64{1, 2, 3}
	sample := sampleValues(values, 10)
	if len(sample) != 3 {
		t.Fatalf("unexpected sample size")
	}
	large := make([]float64, 100)
	for i := range large {
		large[i] = float64(i)
	}
	sample = sampleValues(large, 10)
	if len(sample) == 0 || len(sample) > 10 {
		t.Fatalf("unexpected sample size")
	}
	sample = sampleValues(nil, 10)
	if sample != nil {
		t.Fatalf("expected nil sample")
	}
}

func TestClampMax(t *testing.T) {
	values := []float64{1, 10, 3}
	clamped := clampMax(values, 5)
	if clamped[1] != 5 {
		t.Fatalf("expected clamp")
	}
	clamped = clampMax(values, -1)
	if clamped[1] == 0 {
		t.Fatalf("expected fallback clamp")
	}
}

func TestApplyGamma(t *testing.T) {
	m := dsp.NewFeatureMap(2, 1)
	m.Set(0, 0, 1)
	m.Set(1, 0, 4)
	applyGamma(&m, 2)
	if m.At(1, 0) <= m.At(0, 0) {
		t.Fatalf("expected gamma effect")
	}
	applyGamma(&m, 0)
	applyGamma(nil, 1)
	empty := dsp.FeatureMap{}
	applyGamma(&empty, 1)
}

func testSamples() []float64 {
	samples := make([]float64, 4096)
	for i := range samples {
		t := float64(i) / 4096
		samples[i] = 0.6*math.Sin(2*math.Pi*440*t) + 0.2*math.Sin(2*math.Pi*880*t)
	}
	return samples
}

func colorRGBA(t float64) color.RGBA {
	val := uint8(255 * t)
	return color.RGBA{R: val, G: val, B: val, A: 255}
}
