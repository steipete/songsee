// Package viz builds visualization panels from audio features.
package viz

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sort"
	"strings"

	"github.com/steipete/songsee/internal/dsp"
	"github.com/steipete/songsee/internal/render"
)

// Kind names a visualization.
type Kind string

const (
	Spectrogram Kind = "spectrogram"
	Mel         Kind = "mel"
	Chroma      Kind = "chroma"
	HPSS        Kind = "hpss"
	SelfSim     Kind = "selfsim"
	Loudness    Kind = "loudness"
	Tempogram   Kind = "tempogram"
	MFCC        Kind = "mfcc"
)

var validKinds = map[Kind]struct{}{
	Spectrogram: {},
	Mel:         {},
	Chroma:      {},
	HPSS:        {},
	SelfSim:     {},
	Loudness:    {},
	Tempogram:   {},
	MFCC:        {},
}

// ParseList normalizes a list of viz names, allowing comma-separated values.
func ParseList(raw []string) ([]Kind, error) {
	if len(raw) == 0 {
		return []Kind{Spectrogram}, nil
	}
	seen := map[Kind]bool{}
	out := make([]Kind, 0, len(raw))
	for _, entry := range raw {
		for _, part := range strings.Split(entry, ",") {
			name := strings.ToLower(strings.TrimSpace(part))
			if name == "" {
				continue
			}
			kind := Kind(name)
			if _, ok := validKinds[kind]; !ok {
				return nil, fmt.Errorf("unknown viz: %s", name)
			}
			if seen[kind] {
				continue
			}
			seen[kind] = true
			out = append(out, kind)
		}
	}
	if len(out) == 0 {
		return []Kind{Spectrogram}, nil
	}
	return out, nil
}

// KindsHelp returns the supported viz list in deterministic order.
func KindsHelp() string {
	names := make([]string, 0, len(validKinds))
	for kind := range validKinds {
		names = append(names, string(kind))
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

// Context holds shared analysis data for multiple visualizations.
type Context struct {
	Samples    []float64
	SampleRate int
	WindowSize int
	HopSize    int
	Spec       dsp.Spectrogram
	power      []float64
}

// NewContext analyzes the samples and prepares the base spectrogram.
func NewContext(samples []float64, sampleRate, windowSize, hopSize int) *Context {
	spec := dsp.ComputeSpectrogram(samples, sampleRate, windowSize, hopSize)
	return &Context{
		Samples:    samples,
		SampleRate: sampleRate,
		WindowSize: windowSize,
		HopSize:    hopSize,
		Spec:       spec,
	}
}

// Power returns cached linear power for the spectrogram.
func (c *Context) Power() []float64 {
	if c.power == nil {
		c.power = dsp.SpectrogramPower(&c.Spec)
	}
	return c.power
}

// RenderOptions configures a visualization render.
type RenderOptions struct {
	Width   int
	Height  int
	Palette render.Palette
	MinFreq float64
	MaxFreq float64
}

// Render builds a visualization panel image for the given kind.
func Render(kind Kind, ctx *Context, opts RenderOptions) (*image.RGBA, error) {
	switch kind {
	case Spectrogram:
		return render.Spectrogram(&ctx.Spec, render.Options{
			Width:   opts.Width,
			Height:  opts.Height,
			MinFreq: opts.MinFreq,
			MaxFreq: opts.MaxFreq,
			Palette: opts.Palette,
		})
	case Mel:
		mel := dsp.MelSpectrogramFromPower(&ctx.Spec, ctx.Power(), 0, opts.MinFreq, opts.MaxFreq)
		return render.Heatmap(&mel, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			FlipVert: true,
		})
	case Chroma:
		chroma := dsp.ChromaFromPower(&ctx.Spec, ctx.Power())
		return render.Heatmap(&chroma, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			FlipVert: true,
		})
	case MFCC:
		mfcc := dsp.MFCCFromPower(&ctx.Spec, ctx.Power(), 0, 0, opts.MinFreq, opts.MaxFreq)
		return render.Heatmap(&mfcc, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			FlipVert: true,
		})
	case HPSS:
		return renderHPSS(ctx, opts)
	case SelfSim:
		chroma := dsp.ChromaFromPower(&ctx.Spec, ctx.Power())
		self := dsp.SelfSimilarity(chroma, 200)
		return render.Heatmap(&self, render.HeatmapOptions{
			Width:   opts.Width,
			Height:  opts.Height,
			Palette: opts.Palette,
		})
	case Loudness:
		rms := dsp.RMSFrames(ctx.Samples, ctx.WindowSize, ctx.HopSize)
		return render.Loudness(rms, opts.Width, opts.Height, opts.Palette)
	case Tempogram:
		temp := dsp.Tempogram(&ctx.Spec, 30, 240, 256)
		return render.Heatmap(&temp, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			FlipVert: true,
		})
	default:
		return nil, fmt.Errorf("unknown viz: %s", kind)
	}
}

func renderHPSS(ctx *Context, opts RenderOptions) (*image.RGBA, error) {
	gap := 4
	half := (opts.Height - gap) / 2
	if half <= 0 {
		return nil, fmt.Errorf("invalid output size")
	}
	harm, perc := dsp.HPSS(&ctx.Spec, 9, 9)
	top, err := render.Heatmap(&harm, render.HeatmapOptions{
		Width:    opts.Width,
		Height:   half,
		Palette:  opts.Palette,
		FlipVert: true,
	})
	if err != nil {
		return nil, err
	}
	bottom, err := render.Heatmap(&perc, render.HeatmapOptions{
		Width:    opts.Width,
		Height:   opts.Height - gap - half,
		Palette:  opts.Palette,
		FlipVert: true,
	})
	if err != nil {
		return nil, err
	}
	canvas := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(0, 0, opts.Width, half), top, image.Point{}, draw.Over)
	draw.Draw(canvas, image.Rect(0, half+gap, opts.Width, half+gap+bottom.Bounds().Dy()), bottom, image.Point{}, draw.Over)
	return canvas, nil
}
