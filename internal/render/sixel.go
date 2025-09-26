package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
)

// SixelRenderer renders images using sixel graphics protocol
type SixelRenderer struct {
	BaseRenderer
}

func NewSixelRenderer() *SixelRenderer {
	return &SixelRenderer{}
}

func (r *SixelRenderer) Initialize(videoWidth, videoHeight int) error {
	charAspect := 1.0 // Sixel pixels are square
	return r.calculateDimensions(videoWidth, videoHeight, charAspect)
}

// Render converts image.Image to Sixel string
func (r *SixelRenderer) Render(img image.Image) string {
	bounds := img.Bounds()
	dx := float64(bounds.Dx()) / float64(r.frameWidth)
	dy := float64(bounds.Dy()) / float64(r.frameHeight)

	var buf bytes.Buffer

	// Start sixel sequence
	buf.WriteString("\033Pq")

	// Simple approach: define colors as we encounter them
	colorMap := make(map[uint32]int)
	nextColorId := 0

	// First pass: collect unique colors and assign IDs
	for y := 0; y < r.frameHeight; y++ {
		for x := 0; x < r.frameWidth; x++ {
			px := int(float64(x) * dx)
			py := int(float64(y) * dy)
			c := color.RGBAModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+py)).(color.RGBA)

			// Create a color key
			colorKey := uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)

			if _, exists := colorMap[colorKey]; !exists && nextColorId < 256 {
				colorMap[colorKey] = nextColorId
				// Define the color in sixel format
				buf.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d",
					nextColorId,
					c.R*100/255, // sixel uses 0-100 range
					c.G*100/255,
					c.B*100/255))
				nextColorId++
			}
		}
	}

	// Second pass: render pixels
	for y := 0; y < r.frameHeight; y += 6 { // sixel processes 6 rows at a time
		for colorId := 0; colorId < nextColorId; colorId++ {
			buf.WriteString(fmt.Sprintf("#%d", colorId))

			for x := 0; x < r.frameWidth; x++ {
				sixelChar := 0

				// Check up to 6 pixels vertically
				for dyi := 0; dyi < 6 && y+dyi < r.frameHeight; dyi++ {
					px := int(float64(x) * dx)
					py := int(float64(y+dyi) * dy)
					c := color.RGBAModel.Convert(img.At(bounds.Min.X+px, bounds.Min.Y+py)).(color.RGBA)

					colorKey := uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
					if colorMap[colorKey] == colorId {
						sixelChar |= (1 << dyi)
					}
				}

				// Convert to sixel character (0x3F + value)
				if sixelChar > 0 {
					buf.WriteByte(byte(0x3F + sixelChar))
				} else {
					buf.WriteByte('?') // empty sixel
				}
			}
		}
		buf.WriteString("-") // carriage return + line feed for sixel
	}

	// End sixel sequence
	buf.WriteString("\033\\")
	return buf.String()
}
