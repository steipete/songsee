// Package render turns spectrograms into images.
package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Panel places an image at a coordinate in the final canvas.
type Panel struct {
	Image image.Image
	X     int
	Y     int
}

// Compose composites panels into a single RGBA canvas.
func Compose(width, height int, panels []Panel, bg color.RGBA) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid output size")
	}
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)
	for _, panel := range panels {
		if panel.Image == nil {
			continue
		}
		bounds := panel.Image.Bounds()
		target := image.Rect(panel.X, panel.Y, panel.X+bounds.Dx(), panel.Y+bounds.Dy())
		draw.Draw(canvas, target, panel.Image, bounds.Min, draw.Over)
	}
	return canvas, nil
}
