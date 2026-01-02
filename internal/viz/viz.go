// Package viz builds visualization panels from audio features.
package viz

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
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
	Flux        Kind = "flux"
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
	Flux:        {},
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
		minDB, maxDB := percentileRange(ctx.Spec.Values, 0.05, 0.98)
		return render.Spectrogram(&ctx.Spec, render.Options{
			Width:   opts.Width,
			Height:  opts.Height,
			MinFreq: opts.MinFreq,
			MaxFreq: opts.MaxFreq,
			Palette: opts.Palette,
			MinDB:   minDB,
			MaxDB:   maxDB,
			ClampDB: true,
		})
	case Mel:
		mel := dsp.MelSpectrogramFromPower(&ctx.Spec, ctx.Power(), 0, opts.MinFreq, opts.MaxFreq)
		minVal, maxVal := percentileRange(mel.Values, 0.05, 0.98)
		return render.Heatmap(&mel, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			Min:      minVal,
			Max:      maxVal,
			Clamp:    true,
			FlipVert: true,
		})
	case Chroma:
		chroma := dsp.ChromaFromPower(&ctx.Spec, ctx.Power())
		minVal, maxVal := percentileRange(chroma.Values, 0.1, 0.98)
		return render.Heatmap(&chroma, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			Min:      minVal,
			Max:      maxVal,
			Clamp:    true,
			FlipVert: true,
		})
	case MFCC:
		mfcc := dsp.MFCCFromPower(&ctx.Spec, ctx.Power(), 0, 0, opts.MinFreq, opts.MaxFreq)
		minVal, maxVal := percentileRange(mfcc.Values, 0.05, 0.98)
		return render.Heatmap(&mfcc, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			Min:      minVal,
			Max:      maxVal,
			Clamp:    true,
			FlipVert: true,
		})
	case HPSS:
		return renderHPSS(ctx, opts)
	case SelfSim:
		chroma := dsp.ChromaFromPower(&ctx.Spec, ctx.Power())
		self := dsp.SelfSimilarity(chroma, 200)
		applyGamma(&self, 1.4)
		minVal, maxVal := percentileRange(self.Values, 0.1, 0.98)
		return render.Heatmap(&self, render.HeatmapOptions{
			Width:   opts.Width,
			Height:  opts.Height,
			Palette: opts.Palette,
			Min:     minVal,
			Max:     maxVal,
			Clamp:   true,
		})
	case Loudness:
		rms := dsp.RMSFrames(ctx.Samples, ctx.WindowSize, ctx.HopSize)
		clamped := clampMax(rms, percentileValue(rms, 0.95))
		return render.Loudness(clamped, opts.Width, opts.Height, opts.Palette)
	case Tempogram:
		temp := dsp.Tempogram(&ctx.Spec, 30, 240, 256)
		minVal, maxVal := percentileRange(temp.Values, 0.05, 0.98)
		return render.Heatmap(&temp, render.HeatmapOptions{
			Width:    opts.Width,
			Height:   opts.Height,
			Palette:  opts.Palette,
			Min:      minVal,
			Max:      maxVal,
			Clamp:    true,
			FlipVert: true,
		})
	case Flux:
		flux := dsp.SpectralFlux(&ctx.Spec)
		clamped := clampMax(flux, percentileValue(flux, 0.95))
		return render.Loudness(clamped, opts.Width, opts.Height, opts.Palette)
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
	hMin, hMax := percentileRange(harm.Values, 0.05, 0.98)
	top, err := render.Heatmap(&harm, render.HeatmapOptions{
		Width:    opts.Width,
		Height:   half,
		Palette:  opts.Palette,
		Min:      hMin,
		Max:      hMax,
		Clamp:    true,
		FlipVert: true,
	})
	if err != nil {
		return nil, err
	}
	pMin, pMax := percentileRange(perc.Values, 0.05, 0.98)
	bottom, err := render.Heatmap(&perc, render.HeatmapOptions{
		Width:    opts.Width,
		Height:   opts.Height - gap - half,
		Palette:  opts.Palette,
		Min:      pMin,
		Max:      pMax,
		Clamp:    true,
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

func percentileRange(values []float64, low, high float64) (minVal, maxVal float64) {
	sample := sampleValues(values, 20000)
	if len(sample) == 0 {
		return 0, 1
	}
	sort.Float64s(sample)
	lo := percentileIndex(sample, low)
	hi := percentileIndex(sample, high)
	minVal = sample[lo]
	maxVal = sample[hi]
	if maxVal <= minVal {
		maxVal = minVal + 1e-6
	}
	return minVal, maxVal
}

func percentileValue(values []float64, pct float64) float64 {
	sample := sampleValues(values, 20000)
	if len(sample) == 0 {
		return 1
	}
	sort.Float64s(sample)
	return sample[percentileIndex(sample, pct)]
}

func percentileIndex(values []float64, pct float64) int {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	return int(math.Round(pct * float64(len(values)-1)))
}

func sampleValues(values []float64, maxSamples int) []float64 {
	if len(values) == 0 {
		return nil
	}
	if len(values) <= maxSamples {
		out := make([]float64, len(values))
		copy(out, values)
		return out
	}
	stride := len(values) / maxSamples
	if stride < 1 {
		stride = 1
	}
	out := make([]float64, 0, maxSamples)
	for i := 0; i < len(values); i += stride {
		out = append(out, values[i])
	}
	return out
}

func clampMax(values []float64, maxVal float64) []float64 {
	if len(values) == 0 {
		return values
	}
	if maxVal <= 0 {
		maxVal = 1
	}
	out := make([]float64, len(values))
	for i, v := range values {
		if v > maxVal {
			v = maxVal
		}
		out[i] = v
	}
	return out
}

func applyGamma(mapIn *dsp.FeatureMap, gamma float64) {
	if mapIn == nil || len(mapIn.Values) == 0 {
		return
	}
	if gamma <= 0 {
		gamma = 1
	}
	mapIn.Min = math.Inf(1)
	mapIn.Max = math.Inf(-1)
	for i, v := range mapIn.Values {
		if v < 0 {
			mapIn.Values[i] = v
			continue
		}
		out := math.Pow(v, gamma)
		mapIn.Values[i] = out
		if out < mapIn.Min {
			mapIn.Min = out
		}
		if out > mapIn.Max {
			mapIn.Max = out
		}
	}
}
