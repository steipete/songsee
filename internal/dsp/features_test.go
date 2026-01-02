package dsp

import (
	"math"
	"testing"
)

func TestFeatureMapSet(t *testing.T) {
	m := NewFeatureMap(2, 2)
	m.Set(0, 0, 1)
	m.Set(1, 0, 2)
	m.Set(0, 1, -1)
	if m.Min >= m.Max {
		t.Fatalf("min/max not updated")
	}
	if m.At(1, 0) != 2 {
		t.Fatalf("at mismatch")
	}
}

func TestMelSpectrogram(t *testing.T) {
	spec := testSpectrogram()
	mel := MelSpectrogram(&spec, 20, 0, 0)
	if mel.Width != spec.Frames || mel.Height != 20 {
		t.Fatalf("mel size mismatch")
	}
	if mel.Min >= mel.Max {
		t.Fatalf("mel min/max not set")
	}
}

func TestChroma(t *testing.T) {
	spec := testSpectrogram()
	chroma := Chroma(&spec)
	if chroma.Width != spec.Frames || chroma.Height != 12 {
		t.Fatalf("chroma size mismatch")
	}
	if chroma.Min >= chroma.Max {
		t.Fatalf("chroma min/max not set")
	}
}

func TestMFCC(t *testing.T) {
	spec := testSpectrogram()
	mfcc := MFCC(&spec, 32, 13, 0, 0)
	if mfcc.Width != spec.Frames || mfcc.Height != 13 {
		t.Fatalf("mfcc size mismatch")
	}
	if mfcc.Min >= mfcc.Max {
		t.Fatalf("mfcc min/max not set")
	}
}

func TestHPSS(t *testing.T) {
	spec := testSpectrogram()
	harm, perc := HPSS(&spec, 5, 5)
	if harm.Width != spec.Frames || harm.Height != spec.Bins {
		t.Fatalf("hpss size mismatch")
	}
	if perc.Width != spec.Frames || perc.Height != spec.Bins {
		t.Fatalf("hpss size mismatch")
	}
	if harm.Min >= harm.Max || perc.Min >= perc.Max {
		t.Fatalf("hpss min/max not set")
	}
}

func TestHPSSDefaults(t *testing.T) {
	spec := testSpectrogram()
	_, _ = HPSS(&spec, 0, 0)
}

func TestSelfSimilarity(t *testing.T) {
	m := NewFeatureMap(3, 2)
	m.Set(0, 0, 1)
	m.Set(0, 1, 0)
	m.Set(1, 0, 1)
	m.Set(1, 1, 0)
	m.Set(2, 0, 0)
	m.Set(2, 1, 1)
	ss := SelfSimilarity(m, 0)
	if ss.Width != 3 || ss.Height != 3 {
		t.Fatalf("selfsim size mismatch")
	}
	if ss.At(0, 0) < 0.9 {
		t.Fatalf("selfsim diag expected high")
	}
}

func TestTempogram(t *testing.T) {
	spec := testSpectrogram()
	temp := Tempogram(&spec, 60, 120, 32)
	if temp.Width == 0 || temp.Height != 61 {
		t.Fatalf("tempogram size mismatch")
	}
	if math.IsInf(temp.Min, 0) || math.IsInf(temp.Max, 0) {
		t.Fatalf("tempogram min/max not set")
	}
}

func TestRMSFrames(t *testing.T) {
	samples := make([]float64, 2048)
	for i := range samples {
		samples[i] = 0.5
	}
	rms := RMSFrames(samples, 512, 256)
	if len(rms) == 0 {
		t.Fatalf("rms empty")
	}
	if rms[0] <= 0 {
		t.Fatalf("rms value not set")
	}
}

func TestRMSFramesDefaults(t *testing.T) {
	samples := make([]float64, 1024)
	rms := RMSFrames(samples, 0, 0)
	if len(rms) == 0 {
		t.Fatalf("rms default empty")
	}
}

func TestDownsampleFeatureMap(t *testing.T) {
	m := NewFeatureMap(10, 2)
	for x := 0; x < 10; x++ {
		m.Set(x, 0, float64(x))
		m.Set(x, 1, float64(x))
	}
	out := DownsampleFeatureMap(m, 5)
	if out.Width != 5 || out.Height != 2 {
		t.Fatalf("downsample size mismatch")
	}
}

func TestDownsampleFeatureMapNoop(t *testing.T) {
	m := NewFeatureMap(4, 1)
	m.Set(0, 0, 1)
	out := DownsampleFeatureMap(m, 0)
	if out.Width != 4 {
		t.Fatalf("expected noop downsample")
	}
}

func TestDownsampleFeatureMapNoopLarge(t *testing.T) {
	m := NewFeatureMap(3, 1)
	m.Set(0, 0, 1)
	out := DownsampleFeatureMap(m, 10)
	if out.Width != 3 {
		t.Fatalf("expected noop downsample")
	}
}

func TestSpectrogramPower(t *testing.T) {
	spec := Spectrogram{Values: []float64{0, -10}}
	power := SpectrogramPower(&spec)
	if len(power) != 2 {
		t.Fatalf("power size mismatch")
	}
	if power[0] <= power[1] {
		t.Fatalf("expected higher power for 0 dB")
	}
}

func TestSelfSimilarityDownsample(t *testing.T) {
	m := NewFeatureMap(10, 2)
	for x := 0; x < 10; x++ {
		m.Set(x, 0, float64(x))
		m.Set(x, 1, float64(x))
	}
	ss := SelfSimilarity(m, 5)
	if ss.Width != 5 || ss.Height != 5 {
		t.Fatalf("selfsim downsample mismatch")
	}
}

func TestMedian(t *testing.T) {
	if median(nil) != 0 {
		t.Fatalf("median empty")
	}
	if median([]float64{1, 3}) != 2 {
		t.Fatalf("median even")
	}
	if median([]float64{1, 2, 3}) != 2 {
		t.Fatalf("median odd")
	}
}

func TestMelFilterBins(t *testing.T) {
	points := melFilterBins(100, 10, 5, 0, 400)
	if len(points) != 7 {
		t.Fatalf("points size mismatch")
	}
	for i := 1; i < len(points); i++ {
		if points[i] < points[i-1] {
			t.Fatalf("points not monotonic")
		}
	}
}

func TestMelSpectrogramBounds(t *testing.T) {
	spec := testSpectrogram()
	mel := MelSpectrogram(&spec, 0, 1000, 10)
	if mel.Width != spec.Frames || mel.Height == 0 {
		t.Fatalf("mel bounds mismatch")
	}
}

func TestMelSpectrogramNegativeMin(t *testing.T) {
	spec := testSpectrogram()
	mel := MelSpectrogram(&spec, 8, -10, 500)
	if mel.Width != spec.Frames {
		t.Fatalf("mel size mismatch")
	}
}

func TestTempogramHopZero(t *testing.T) {
	spec := Spectrogram{
		Frames:     2,
		Bins:       2,
		Values:     []float64{0, 0, 0, 0},
		Min:        0,
		Max:        1,
		SampleRate: 1,
		HopSize:    0,
		BinHz:      1,
	}
	temp := Tempogram(&spec, 30, 60, 0)
	if temp.Width != 2 {
		t.Fatalf("tempogram width mismatch")
	}
}

func TestTempogramDefaults(t *testing.T) {
	spec := testSpectrogram()
	temp := Tempogram(&spec, 0, 0, 0)
	if temp.Height == 0 {
		t.Fatalf("tempogram defaults mismatch")
	}
}

func TestSpectralFluxShort(t *testing.T) {
	spec := Spectrogram{
		Frames: 1,
		Bins:   2,
		Values: []float64{0, 0},
	}
	flux := SpectralFlux(&spec)
	if len(flux) != 1 {
		t.Fatalf("flux size mismatch")
	}
}

func TestDownsampleSignalNoop(t *testing.T) {
	in := []float64{1, 2, 3}
	out := downsampleSignal(in, 0)
	if len(out) != len(in) {
		t.Fatalf("downsample noop")
	}
}

func TestDownsampleSignal(t *testing.T) {
	in := []float64{1, 2, 3, 4}
	out := downsampleSignal(in, 2)
	if len(out) != 2 {
		t.Fatalf("downsample size mismatch")
	}
}

func TestMelConversions(t *testing.T) {
	if hzToMel(0) != 0 {
		t.Fatalf("hzToMel 0")
	}
	mel := hzToMel(1000)
	if melToHz(mel) <= 0 {
		t.Fatalf("melToHz invalid")
	}
}

func TestPowerConversions(t *testing.T) {
	if powerToDB(1) == 0 {
		t.Fatalf("powerToDB expected non-zero")
	}
	if dbToPower(0) != 1 {
		t.Fatalf("dbToPower expected 1")
	}
}

func TestMFCCDefaults(t *testing.T) {
	spec := testSpectrogram()
	power := SpectrogramPower(&spec)
	mfcc := MFCCFromPower(&spec, power, 0, 0, 0, 0)
	if mfcc.Height == 0 {
		t.Fatalf("mfcc defaults mismatch")
	}
}

func TestSelfSimilarityEmpty(t *testing.T) {
	m := FeatureMap{}
	ss := SelfSimilarity(m, 0)
	if ss.Width != 0 || ss.Height != 0 {
		t.Fatalf("selfsim empty mismatch")
	}
}

func testSpectrogram() Spectrogram {
	samples := make([]float64, 4096)
	for i := range samples {
		t := float64(i) / 4096
		samples[i] = 0.7*math.Sin(2*math.Pi*440*t) + 0.2*math.Sin(2*math.Pi*880*t)
	}
	return ComputeSpectrogram(samples, 44100, 512, 128)
}
