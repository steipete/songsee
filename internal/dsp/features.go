// Package dsp provides spectral analysis utilities.
package dsp

import (
	"math"
	"sort"
)

const (
	defaultMelBands = 40
	defaultMFCC     = 13
)

// SpectrogramPower converts log-magnitude spectrogram values to linear power.
func SpectrogramPower(spec *Spectrogram) []float64 {
	power := make([]float64, len(spec.Values))
	for i, v := range spec.Values {
		power[i] = dbToPower(v)
	}
	return power
}

// MelSpectrogram computes a mel-scaled spectrogram from log-magnitude FFT data.
func MelSpectrogram(spec *Spectrogram, bands int, minFreq, maxFreq float64) FeatureMap {
	return MelSpectrogramFromPower(spec, SpectrogramPower(spec), bands, minFreq, maxFreq)
}

// MelSpectrogramFromPower computes a mel spectrogram from linear power.
func MelSpectrogramFromPower(spec *Spectrogram, power []float64, bands int, minFreq, maxFreq float64) FeatureMap {
	if bands <= 0 {
		bands = defaultMelBands
	}
	if maxFreq <= 0 {
		maxFreq = float64(spec.SampleRate) / 2
	}
	if minFreq < 0 {
		minFreq = 0
	}
	if maxFreq <= minFreq {
		maxFreq = minFreq + 1
	}

	bins := spec.Bins
	frames := spec.Frames
	points := melFilterBins(spec.BinHz, bins, bands, minFreq, maxFreq)
	out := NewFeatureMap(frames, bands)

	for f := 0; f < frames; f++ {
		base := f * bins
		for m := 0; m < bands; m++ {
			start := points[m]
			center := points[m+1]
			end := points[m+2]
			if end <= start {
				continue
			}
			if center <= start {
				center = (start + end) / 2
			}
			if center >= end {
				center = (start + end) / 2
			}
			energy := 0.0
			for b := start; b <= end; b++ {
				var weight float64
				if b < center {
					den := float64(center - start)
					if den > 0 {
						weight = float64(b-start) / den
					}
				} else {
					den := float64(end - center)
					if den > 0 {
						weight = float64(end-b) / den
					}
				}
				if weight < 0 {
					weight = 0
				}
				if weight > 1 {
					weight = 1
				}
				energy += power[base+b] * weight
			}
			out.Set(f, m, powerToDB(energy))
		}
	}
	return out
}

// Chroma computes a 12-bin chromagram from log-magnitude FFT data.
func Chroma(spec *Spectrogram) FeatureMap {
	return ChromaFromPower(spec, SpectrogramPower(spec))
}

// ChromaFromPower computes a 12-bin chromagram from linear power.
func ChromaFromPower(spec *Spectrogram, power []float64) FeatureMap {
	frames := spec.Frames
	bins := spec.Bins
	out := NewFeatureMap(frames, 12)

	for f := 0; f < frames; f++ {
		base := f * bins
		for b := 1; b < bins; b++ {
			freq := float64(b) * spec.BinHz
			if freq < 30 {
				continue
			}
			midi := 69 + 12*math.Log2(freq/440.0)
			class := int(math.Round(midi)) % 12
			if class < 0 {
				class += 12
			}
			out.Set(f, class, out.At(f, class)+power[base+b])
		}
	}

	out.Min = math.Inf(1)
	out.Max = math.Inf(-1)
	for f := 0; f < frames; f++ {
		for c := 0; c < 12; c++ {
			out.Set(f, c, powerToDB(out.At(f, c)))
		}
	}
	return out
}

// MFCC computes MFCC coefficients from log-magnitude FFT data.
func MFCC(spec *Spectrogram, bands, coeffs int, minFreq, maxFreq float64) FeatureMap {
	return MFCCFromPower(spec, SpectrogramPower(spec), bands, coeffs, minFreq, maxFreq)
}

// MFCCFromPower computes MFCC coefficients from linear power.
func MFCCFromPower(spec *Spectrogram, power []float64, bands, coeffs int, minFreq, maxFreq float64) FeatureMap {
	if bands <= 0 {
		bands = defaultMelBands
	}
	if coeffs <= 0 {
		coeffs = defaultMFCC
	}
	if coeffs > bands {
		coeffs = bands
	}
	mel := MelSpectrogramFromPower(spec, power, bands, minFreq, maxFreq)

	out := NewFeatureMap(mel.Width, coeffs)
	tmp := make([]float64, bands)
	for f := 0; f < mel.Width; f++ {
		for m := 0; m < bands; m++ {
			tmp[m] = mel.At(f, m) / 10 * math.Ln10
		}
		for k := 0; k < coeffs; k++ {
			sum := 0.0
			for n := 0; n < bands; n++ {
				sum += tmp[n] * math.Cos(math.Pi/float64(bands)*(float64(n)+0.5)*float64(k))
			}
			out.Set(f, k, sum)
		}
	}
	return out
}

// HPSS separates harmonic and percussive content using median filters.
func HPSS(spec *Spectrogram, timeWidth, freqWidth int) (harm, perc FeatureMap) {
	if timeWidth <= 0 {
		timeWidth = 9
	}
	if freqWidth <= 0 {
		freqWidth = 9
	}
	frames := spec.Frames
	bins := spec.Bins
	harm = NewFeatureMap(frames, bins)
	perc = NewFeatureMap(frames, bins)

	timeRadius := timeWidth / 2
	freqRadius := freqWidth / 2

	timeBuf := make([]float64, 0, timeWidth)
	freqBuf := make([]float64, 0, freqWidth)
	for f := 0; f < frames; f++ {
		for b := 0; b < bins; b++ {
			timeBuf = timeBuf[:0]
			for tf := f - timeRadius; tf <= f+timeRadius; tf++ {
				if tf < 0 || tf >= frames {
					continue
				}
				timeBuf = append(timeBuf, spec.Values[tf*bins+b])
			}
			freqBuf = freqBuf[:0]
			for tb := b - freqRadius; tb <= b+freqRadius; tb++ {
				if tb < 0 || tb >= bins {
					continue
				}
				freqBuf = append(freqBuf, spec.Values[f*bins+tb])
			}
			hMed := median(timeBuf)
			pMed := median(freqBuf)
			hPow := dbToPower(hMed)
			pPow := dbToPower(pMed)
			src := dbToPower(spec.Values[f*bins+b])
			den := hPow + pPow + 1e-12
			hVal := src * hPow / den
			pVal := src * pPow / den
			harm.Set(f, b, powerToDB(hVal))
			perc.Set(f, b, powerToDB(pVal))
		}
	}
	return harm, perc
}

// SpectralFlux computes the spectral flux across frames.
func SpectralFlux(spec *Spectrogram) []float64 {
	frames := spec.Frames
	bins := spec.Bins
	flux := make([]float64, frames)
	for f := 1; f < frames; f++ {
		base := f * bins
		prev := (f - 1) * bins
		sum := 0.0
		for b := 0; b < bins; b++ {
			diff := spec.Values[base+b] - spec.Values[prev+b]
			if diff > 0 {
				sum += diff
			}
		}
		flux[f] = sum
	}
	return flux
}

// Tempogram computes a tempo map from spectral flux.
func Tempogram(spec *Spectrogram, minBPM, maxBPM, maxFrames int) FeatureMap {
	if minBPM <= 0 {
		minBPM = 30
	}
	if maxBPM <= minBPM {
		maxBPM = minBPM + 60
	}
	flux := SpectralFlux(spec)
	if maxFrames > 0 && len(flux) > maxFrames {
		flux = downsampleSignal(flux, maxFrames)
	}
	frames := len(flux)
	bpmBins := maxBPM - minBPM + 1
	out := NewFeatureMap(frames, bpmBins)

	fps := float64(spec.SampleRate) / float64(spec.HopSize)
	if fps <= 0 {
		fps = 1
	}
	window := int(math.Round(fps * 8))
	if window < 8 {
		window = 8
	}
	if window > frames {
		window = frames
	}

	for t := 0; t < frames; t++ {
		start := t - window/2
		end := t + window/2
		if start < 0 {
			start = 0
		}
		if end >= frames {
			end = frames - 1
		}
		for bpm := minBPM; bpm <= maxBPM; bpm++ {
			lag := int(math.Round(fps * 60 / float64(bpm)))
			if lag <= 0 {
				continue
			}
			sum := 0.0
			for i := start + lag; i <= end; i++ {
				sum += flux[i] * flux[i-lag]
			}
			out.Set(t, bpm-minBPM, sum)
		}
	}
	return out
}

// SelfSimilarity computes a self-similarity matrix from a feature map.
func SelfSimilarity(mapIn FeatureMap, maxFrames int) FeatureMap {
	features := mapIn
	if maxFrames > 0 && mapIn.Width > maxFrames {
		features = DownsampleFeatureMap(mapIn, maxFrames)
	}
	frames := features.Width
	out := NewFeatureMap(frames, frames)
	if frames == 0 {
		return out
	}
	norms := make([]float64, frames)
	for f := 0; f < frames; f++ {
		sum := 0.0
		for k := 0; k < features.Height; k++ {
			v := features.At(f, k)
			sum += v * v
		}
		norms[f] = math.Sqrt(sum)
	}
	for i := 0; i < frames; i++ {
		for j := 0; j < frames; j++ {
			dot := 0.0
			for k := 0; k < features.Height; k++ {
				dot += features.At(i, k) * features.At(j, k)
			}
			den := norms[i] * norms[j]
			sim := 0.0
			if den > 0 {
				sim = dot / den
			}
			out.Set(i, j, sim)
		}
	}
	return out
}

// DownsampleFeatureMap reduces the time axis by averaging frame windows.
func DownsampleFeatureMap(mapIn FeatureMap, maxFrames int) FeatureMap {
	if mapIn.Width <= maxFrames || maxFrames <= 0 {
		return mapIn
	}
	out := NewFeatureMap(maxFrames, mapIn.Height)
	ratio := float64(mapIn.Width) / float64(maxFrames)
	for x := 0; x < maxFrames; x++ {
		start := int(math.Floor(float64(x) * ratio))
		end := int(math.Floor(float64(x+1) * ratio))
		if end <= start {
			end = start + 1
		}
		if end > mapIn.Width {
			end = mapIn.Width
		}
		count := float64(end - start)
		for y := 0; y < mapIn.Height; y++ {
			sum := 0.0
			for i := start; i < end; i++ {
				sum += mapIn.At(i, y)
			}
			out.Set(x, y, sum/count)
		}
	}
	return out
}

// RMSFrames computes RMS per frame for the given samples.
func RMSFrames(samples []float64, windowSize, hopSize int) []float64 {
	if windowSize <= 0 {
		windowSize = 2048
	}
	if hopSize <= 0 {
		hopSize = windowSize / 4
	}
	frames := 1
	if len(samples) > windowSize {
		frames = 1 + (len(samples)-windowSize+hopSize-1)/hopSize
	}
	out := make([]float64, frames)
	for f := 0; f < frames; f++ {
		start := f * hopSize
		sum := 0.0
		count := 0.0
		for i := 0; i < windowSize; i++ {
			idx := start + i
			if idx >= len(samples) {
				break
			}
			v := samples[idx]
			sum += v * v
			count++
		}
		if count > 0 {
			out[f] = math.Sqrt(sum / count)
		}
	}
	return out
}

func melFilterBins(binHz float64, bins, bands int, minFreq, maxFreq float64) []int {
	minMel := hzToMel(minFreq)
	maxMel := hzToMel(maxFreq)
	points := make([]int, bands+2)
	for i := 0; i < bands+2; i++ {
		mel := minMel + (maxMel-minMel)*float64(i)/float64(bands+1)
		hz := melToHz(mel)
		bin := int(math.Round(hz / binHz))
		if bin < 0 {
			bin = 0
		}
		if bin > bins-1 {
			bin = bins - 1
		}
		points[i] = bin
	}
	for i := 1; i < len(points); i++ {
		if points[i] < points[i-1] {
			points[i] = points[i-1]
		}
	}
	return points
}

func hzToMel(hz float64) float64 {
	return 2595 * math.Log10(1+hz/700)
}

func melToHz(mel float64) float64 {
	return 700 * (math.Pow(10, mel/2595) - 1)
}

func downsampleSignal(in []float64, maxFrames int) []float64 {
	if len(in) <= maxFrames || maxFrames <= 0 {
		return in
	}
	out := make([]float64, maxFrames)
	ratio := float64(len(in)) / float64(maxFrames)
	for x := 0; x < maxFrames; x++ {
		start := int(math.Floor(float64(x) * ratio))
		end := int(math.Floor(float64(x+1) * ratio))
		if end <= start {
			end = start + 1
		}
		if end > len(in) {
			end = len(in)
		}
		sum := 0.0
		for i := start; i < end; i++ {
			sum += in[i]
		}
		out[x] = sum / float64(end-start)
	}
	return out
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	tmp := append([]float64(nil), values...)
	sort.Float64s(tmp)
	mid := len(tmp) / 2
	if len(tmp)%2 == 0 {
		return 0.5 * (tmp[mid-1] + tmp[mid])
	}
	return tmp[mid]
}
