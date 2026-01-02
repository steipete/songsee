// Package render turns spectrograms into images.
package render

import (
	"errors"
	"image/color"
)

// Palette maps a normalized value to a color.
type Palette func(t float64) color.RGBA

type stop struct {
	pos float64
	c   color.RGBA
}

// PaletteByName returns a palette for a given name.
func PaletteByName(name string) (Palette, error) {
	switch name {
	case "classic":
		return gradient([]stop{
			{0.0, rgb(0, 0, 0)},
			{0.2, rgb(0, 32, 96)},
			{0.45, rgb(0, 160, 200)},
			{0.7, rgb(255, 180, 0)},
			{1.0, rgb(255, 255, 255)},
		}), nil
	case "magma":
		return gradient([]stop{
			{0.0, rgb(0, 0, 4)},
			{0.25, rgb(59, 12, 87)},
			{0.5, rgb(180, 54, 122)},
			{0.75, rgb(251, 140, 60)},
			{1.0, rgb(252, 253, 191)},
		}), nil
	case "inferno":
		return gradient([]stop{
			{0.0, rgb(0, 0, 4)},
			{0.25, rgb(61, 9, 101)},
			{0.5, rgb(187, 55, 84)},
			{0.75, rgb(249, 142, 8)},
			{1.0, rgb(252, 255, 164)},
		}), nil
	case "viridis":
		return gradient([]stop{
			{0.0, rgb(68, 1, 84)},
			{0.25, rgb(58, 82, 139)},
			{0.5, rgb(32, 144, 140)},
			{0.75, rgb(94, 201, 98)},
			{1.0, rgb(253, 231, 37)},
		}), nil
	case "clawd":
		// 🦞 Lobster from the deep! Ocean depths to coral brightness
		return gradient([]stop{
			{0.0, rgb(2, 4, 15)},      // Abyss black-blue
			{0.2, rgb(11, 38, 74)},    // Deep ocean navy
			{0.4, rgb(18, 97, 117)},   // Ocean teal
			{0.6, rgb(193, 98, 92)},   // Coral/salmon
			{0.8, rgb(205, 55, 40)},   // Lobster red! 🦞
			{1.0, rgb(255, 230, 210)}, // Foam/shell highlight
		}), nil
	case "gray", "grey":
		return gradient([]stop{{0, rgb(0, 0, 0)}, {1, rgb(255, 255, 255)}}), nil
	default:
		return nil, errors.New("unknown palette")
	}
}

func gradient(stops []stop) Palette {
	return func(t float64) color.RGBA {
		if t <= 0 {
			return stops[0].c
		}
		if t >= 1 {
			return stops[len(stops)-1].c
		}
		for i := 0; i < len(stops)-1; i++ {
			if t >= stops[i].pos && t <= stops[i+1].pos {
				span := stops[i+1].pos - stops[i].pos
				if span <= 0 {
					return stops[i+1].c
				}
				local := (t - stops[i].pos) / span
				return lerp(stops[i].c, stops[i+1].c, local)
			}
		}
		return stops[len(stops)-1].c
	}
}

func lerp(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: 255,
	}
}

func rgb(r, g, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
