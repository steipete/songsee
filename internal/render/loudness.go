// Package render turns spectrograms into images.
package render

import (
	"fmt"
	"image"
	"math"
)

// Loudness renders a loudness curve into an RGBA image.
func Loudness(values []float64, width, height int, palette Palette) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid output size")
	}
	if palette == nil {
		return nil, fmt.Errorf("palette required")
	}
	if len(values) == 0 {
		return image.NewRGBA(image.Rect(0, 0, width, height)), nil
	}

	maxVal := 0.0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal <= 0 {
		return image.NewRGBA(image.Rect(0, 0, width, height)), nil
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		srcX := 0
		if len(values) > 1 && width > 1 {
			srcX = int(math.Round(float64(x) * float64(len(values)-1) / float64(width-1)))
		}
		norm := values[srcX] / maxVal
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}
		level := int(math.Round(norm * float64(height-1)))
		col := palette(norm)
		for y := height - 1; y >= height-1-level; y-- {
			if y < 0 || y >= height {
				continue
			}
			img.SetRGBA(x, y, col)
		}
	}
	return img, nil
}
