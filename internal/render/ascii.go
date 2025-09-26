package render

import (
	"image"
	"image/color"
)

var asciiRamp = []rune("█▓▒@#WMB8&%$0QOC?!i;:,. ")

type ASCIIRenderer struct {
	BaseRenderer
}

func NewASCIIRenderer() *ASCIIRenderer {
	return &ASCIIRenderer{}
}

func (r *ASCIIRenderer) Initialize(videoWidth, videoHeight int) error {
	charAspect := 0.5 // ASCII character aspect ratio
	return r.calculateDimensions(videoWidth, videoHeight, charAspect)
}

func (r *ASCIIRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.frameWidth)
	dy := float64(bounds.Dy()) / float64(r.frameHeight)

	out := make([]rune, 0, (r.frameWidth+1)*r.frameHeight)

	for y := 0; y < r.frameHeight; y++ {
		for x := 0; x < r.frameWidth; x++ {
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
