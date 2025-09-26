package render

import (
	"fmt"
	"image"
	"image/color"
)

type BrailleRenderer struct {
	BaseRenderer
}

func NewBrailleRenderer() *BrailleRenderer {
	return &BrailleRenderer{}
}

func (r *BrailleRenderer) Initialize(videoWidth, videoHeight int) error {
	charAspect := 0.5 // Braille: 2x4 dots
	return r.calculateDimensions(videoWidth, videoHeight, charAspect)
}

func (r *BrailleRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.frameWidth*2)
	dy := float64(bounds.Dy()) / float64(r.frameHeight*4)

	str := ""
	for y := 0; y < r.frameHeight; y++ {
		for x := 0; x < r.frameWidth; x++ {
			cell := 0
			var avgR, avgG, avgB float64
			count := 0

			for cy := 0; cy < 4; cy++ {
				for cx := 0; cx < 2; cx++ {
					px := int(float64(x*2+cx) * dx)
					py := int(float64(y*4+cy) * dy)
					c := color.NRGBAModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+py)).(color.NRGBA)
					lum := float64(c.R+c.G+c.B) / 3.0
					if lum < 128 {
						cell |= brailleBit(cx, cy)
					}
					avgR += float64(c.R)
					avgG += float64(c.G)
					avgB += float64(c.B)
					count++
				}
			}

			avgR /= float64(count)
			avgG /= float64(count)
			avgB /= float64(count)

			str += fmt.Sprintf("\033[38;2;%d;%d;%dm%s",
				int(avgR), int(avgG), int(avgB),
				string(rune(0x2800+cell)))
		}
		str += "\033[0m\n"
	}

	return str
}

func brailleBit(x, y int) int {
	idx := y*2 + x
	switch idx {
	case 0:
		return 1 << 0
	case 1:
		return 1 << 3
	case 2:
		return 1 << 1
	case 3:
		return 1 << 4
	case 4:
		return 1 << 2
	case 5:
		return 1 << 5
	case 6:
		return 1 << 6
	case 7:
		return 1 << 7
	}
	return 0
}
