package render

import (
	"fmt"
	"image"
	"image/color"
)

type BlockRenderer struct {
	BaseRenderer
}

func NewBlockRenderer() *BlockRenderer {
	return &BlockRenderer{}
}

func (r *BlockRenderer) Initialize(videoWidth, videoHeight int) error {
	charAspect := 0.5 // Block: 2x1 vertical (each char represents 2 vertical pixels)
	return r.calculateDimensions(videoWidth, videoHeight, charAspect)
}

func (r *BlockRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.frameWidth)
	dy := float64(bounds.Dy()) / float64(r.frameHeight*2)

	str := ""
	for y := 0; y < r.frameHeight; y++ {
		for x := 0; x < r.frameWidth; x++ {
			px := int(float64(x) * dx)
			pyTop := int(float64(y*2) * dy)
			pyBot := int(float64(y*2+1) * dy)

			cTop := color.NRGBAModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+pyTop)).(color.NRGBA)
			cBot := color.NRGBAModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+pyBot)).(color.NRGBA)

			// Calculate luminance for both pixels
			lumTop := (float64(cTop.R+cTop.G+cTop.B) / 3.0)
			lumBot := (float64(cBot.R+cBot.G+cBot.B) / 3.0)

			topDark := lumTop < 128
			botDark := lumBot < 128

			// Choose block character based on which pixels are dark
			var ch rune
			var colorToUse color.NRGBA

			switch {
			case topDark && botDark:
				ch = '█' // full block
				colorToUse = cTop
			case topDark && !botDark:
				ch = '▀' // upper half block
				colorToUse = cTop
			case !topDark && botDark:
				ch = '▄' // lower half block
				colorToUse = cBot
			default:
				ch = ' ' // empty space
				colorToUse = cTop
			}

			str += fmt.Sprintf("\033[38;2;%d;%d;%dm%c",
				colorToUse.R, colorToUse.G, colorToUse.B, ch)
		}
		str += "\033[0m\n" // reset color at end of line
	}

	return str
}
