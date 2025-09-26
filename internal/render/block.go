package render

import (
	"image"
	"image/color"
)

// Block characters give ~2x vertical resolution
var blockChars = []rune(" █▀▄")

type BlockRenderer struct {
	width  int
	height int
}

func NewBlockRenderer(width, height int) *BlockRenderer {
	return &BlockRenderer{width: width, height: height}
}

func (r *BlockRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.width)
	dy := float64(bounds.Dy()) / float64(r.height*2)

	out := make([]rune, 0, (r.width+1)*r.height)

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			px := int(float64(x) * dx)
			pyTop := int(float64(y*2) * dy)
			pyBot := int(float64(y*2+1) * dy)

			cTop := color.GrayModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+pyTop)).(color.Gray)
			cBot := color.GrayModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+pyBot)).(color.Gray)

			lumTop := float64(cTop.Y) / 255.0
			lumBot := float64(cBot.Y) / 255.0

			// Threshold for block rendering
			topDark := lumTop < 0.5
			botDark := lumBot < 0.5

			switch {
			case topDark && botDark:
				out = append(out, '█')
			case topDark && !botDark:
				out = append(out, '▀')
			case !topDark && botDark:
				out = append(out, '▄')
			default:
				out = append(out, ' ')
			}
		}
		out = append(out, '\n')
	}

	return string(out)
}
