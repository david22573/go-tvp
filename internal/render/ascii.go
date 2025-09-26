package render

import (
	"image"
	"image/color"
)

// Extended ASCII ramp (dark → light)
var asciiRamp = []rune("█▓▒@#WMB8&%$0QOC?!i;:,. ")

type ASCIIRenderer struct {
	width  int
	height int
}

func NewASCIIRenderer(width, height int) *ASCIIRenderer {
	return &ASCIIRenderer{width: width, height: height}
}

func (r *ASCIIRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.width)
	dy := float64(bounds.Dy()) / float64(r.height)

	out := make([]rune, 0, (r.width+1)*r.height)

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			px := int(float64(x) * dx)
			py := int(float64(y) * dy)

			c := color.GrayModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+py)).(color.Gray)
			lum := float64(c.Y) / 255.0

			idx := int(lum * float64(len(asciiRamp)-1))
			out = append(out, asciiRamp[idx])
		}
		out = append(out, '\n')
	}

	return string(out)
}
