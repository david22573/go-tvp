package render

import (
	"image"
	"image/color"
)

// Braille uses 2x4 dots per cell → very high resolution
type BrailleRenderer struct {
	width  int
	height int
}

func NewBrailleRenderer(width, height int) *BrailleRenderer {
	return &BrailleRenderer{width: width, height: height}
}

func (r *BrailleRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.width*2)
	dy := float64(bounds.Dy()) / float64(r.height*4)

	out := make([]rune, 0, (r.width+1)*r.height)

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			cell := 0
			for cy := 0; cy < 4; cy++ {
				for cx := 0; cx < 2; cx++ {
					px := int(float64(x*2+cx) * dx)
					py := int(float64(y*4+cy) * dy)
					c := color.GrayModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+py)).(color.Gray)
					lum := float64(c.Y) / 255.0
					if lum < 0.5 {
						bit := brailleBit(cx, cy)
						cell |= bit
					}
				}
			}
			out = append(out, rune(0x2800+cell))
		}
		out = append(out, '\n')
	}

	return string(out)
}

func brailleBit(x, y int) int {
	// Dot numbering (Unicode Braille):
	// 1 4
	// 2 5
	// 3 6
	// 7 8
	idx := y*2 + x
	switch idx {
	case 0:
		return 1 << 0 // dot1
	case 1:
		return 1 << 3 // dot4
	case 2:
		return 1 << 1 // dot2
	case 3:
		return 1 << 4 // dot5
	case 4:
		return 1 << 2 // dot3
	case 5:
		return 1 << 5 // dot6
	case 6:
		return 1 << 6 // dot7
	case 7:
		return 1 << 7 // dot8
	}
	return 0
}
