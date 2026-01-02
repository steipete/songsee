package render

import (
	"image/color"
	"testing"

	"github.com/steipete/songsee/internal/dsp"
)

func TestPaletteByName(t *testing.T) {
	names := []string{"classic", "magma", "inferno", "viridis", "gray", "grey"}
	for _, name := range names {
		if _, err := PaletteByName(name); err != nil {
			t.Fatalf("palette %s: %v", name, err)
		}
	}
	if _, err := PaletteByName("nope"); err == nil {
		t.Fatalf("expected error for unknown palette")
	}
}

func TestRenderSpectrogram(t *testing.T) {
	spec := dsp.Spectrogram{
		Frames: 2,
		Bins:   2,
		Values: []float64{-20, -5, -10, -1},
		Min:    -20,
		Max:    -1,
		BinHz:  100,
	}
	img, err := Spectrogram(&spec, Options{
		Width:   4,
		Height:  4,
		Palette: func(t float64) color.RGBA { return color.RGBA{R: uint8(255 * t), A: 255} },
	})
	if err != nil {
		t.Fatalf("RenderSpectrogram: %v", err)
	}
	if img.Bounds().Dx() != 4 || img.Bounds().Dy() != 4 {
		t.Fatalf("unexpected bounds")
	}
	c1 := img.RGBAAt(0, 0)
	c2 := img.RGBAAt(3, 3)
	if c1 == c2 {
		t.Fatalf("expected varying pixels")
	}
}

func TestRenderSpectrogramErrors(t *testing.T) {
	if _, err := Spectrogram(nil, Options{Width: 1, Height: 1, Palette: func(float64) color.RGBA { return color.RGBA{} }}); err == nil {
		t.Fatalf("expected spec error")
	}
	spec := dsp.Spectrogram{
		Frames: 1,
		Bins:   1,
		Values: []float64{0},
		Min:    0,
		Max:    1,
		BinHz:  100,
	}
	if _, err := Spectrogram(&spec, Options{Width: 0, Height: 1, Palette: func(float64) color.RGBA { return color.RGBA{} }}); err == nil {
		t.Fatalf("expected size error")
	}
	if _, err := Spectrogram(&spec, Options{Width: 1, Height: 1}); err == nil {
		t.Fatalf("expected palette error")
	}
}

func TestRenderSpectrogramClampAndRange(t *testing.T) {
	spec := dsp.Spectrogram{
		Frames: 3,
		Bins:   4,
		Values: []float64{-80, -40, -20, 0, -70, -35, -15, -2, -60, -30, -10, -1},
		Min:    -80,
		Max:    0,
		BinHz:  100,
	}
	img, err := Spectrogram(&spec, Options{
		Width:    3,
		Height:   2,
		MinFreq:  50,
		MaxFreq:  250,
		Palette:  func(t float64) color.RGBA { return color.RGBA{B: uint8(255 * t), A: 255} },
		MinDB:    -60,
		MaxDB:    -10,
		ClampDB:  true,
		FlipVert: true,
	})
	if err != nil {
		t.Fatalf("RenderSpectrogram: %v", err)
	}
	if img.Bounds().Dx() != 3 || img.Bounds().Dy() != 2 {
		t.Fatalf("unexpected bounds")
	}
}

func TestGradientEndpoints(t *testing.T) {
	p := gradient([]stop{{0, rgb(0, 0, 0)}, {1, rgb(255, 0, 0)}})
	if c := p(0); c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("start color mismatch")
	}
	if c := p(1); c.R != 255 || c.G != 0 || c.B != 0 {
		t.Fatalf("end color mismatch")
	}
	if c := p(0.5); c.R == 0 || c.R == 255 {
		t.Fatalf("mid color not interpolated")
	}
	if c := p(-1); c.R != 0 {
		t.Fatalf("clamp low")
	}
	if c := p(2); c.R != 255 {
		t.Fatalf("clamp high")
	}
}

func TestRenderSpectrogramSinglePixel(t *testing.T) {
	spec := dsp.Spectrogram{
		Frames: 1,
		Bins:   1,
		Values: []float64{-10},
		Min:    -10,
		Max:    -10,
		BinHz:  100,
	}
	img, err := Spectrogram(&spec, Options{
		Width:   1,
		Height:  1,
		Palette: func(_ float64) color.RGBA { return color.RGBA{G: 200, A: 255} },
	})
	if err != nil {
		t.Fatalf("RenderSpectrogram: %v", err)
	}
	if img.Bounds().Dx() != 1 || img.Bounds().Dy() != 1 {
		t.Fatalf("unexpected bounds")
	}
}

func TestRenderSpectrogramRangeReset(t *testing.T) {
	spec := dsp.Spectrogram{
		Frames: 2,
		Bins:   3,
		Values: []float64{-10, -5, -1, -10, -5, -1},
		Min:    -10,
		Max:    -1,
		BinHz:  100,
	}
	_, err := Spectrogram(&spec, Options{
		Width:   2,
		Height:  2,
		MinFreq: 1000,
		MaxFreq: 200,
		Palette: func(_ float64) color.RGBA { return color.RGBA{R: 50, A: 255} },
	})
	if err != nil {
		t.Fatalf("RenderSpectrogram: %v", err)
	}
}
