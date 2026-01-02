// Package dsp provides spectral analysis utilities.
package dsp

import "math"

// FeatureMap is a 2D grid of values over time (Width) and features (Height).
// Values are stored in row-major order: idx = y*Width + x.
type FeatureMap struct {
	Width  int
	Height int
	Values []float64
	Min    float64
	Max    float64
}

// NewFeatureMap allocates a map with initialized min/max sentinels.
func NewFeatureMap(width, height int) FeatureMap {
	return FeatureMap{
		Width:  width,
		Height: height,
		Values: make([]float64, width*height),
		Min:    math.Inf(1),
		Max:    math.Inf(-1),
	}
}

// Set writes a value and updates the min/max bounds.
func (m *FeatureMap) Set(x, y int, v float64) {
	if m == nil || m.Width <= 0 || m.Height <= 0 {
		return
	}
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}
	idx := y*m.Width + x
	m.Values[idx] = v
	if v < m.Min {
		m.Min = v
	}
	if v > m.Max {
		m.Max = v
	}
}

// At reads a value without bounds checks.
func (m FeatureMap) At(x, y int) float64 {
	return m.Values[y*m.Width+x]
}

func dbToPower(db float64) float64 {
	return math.Pow(10, db/10)
}

func powerToDB(power float64) float64 {
	return 10 * math.Log10(power+1e-12)
}
