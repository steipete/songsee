package viz

import (
	"image/color"
	"math"
	"testing"
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
	kinds := []Kind{Spectrogram, Mel, Chroma, MFCC, HPSS, SelfSim, Loudness, Tempogram}
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
