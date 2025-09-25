package render

import (
	"bytes"
)

type asciiRenderer struct{}

func (a *asciiRenderer) RenderFrame(rgb []byte, w, h int) string {
	chars := " .:-=+*#%@"
	var b bytes.Buffer

	for i := 0; i < len(rgb); i += 3 {
		r, g, b2 := rgb[i], rgb[i+1], rgb[i+2]
		lum := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b2)
		idx := int((lum / 255.0) * float64(len(chars)-1))
		b.WriteByte(chars[idx])
		if (i/3+1)%w == 0 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
