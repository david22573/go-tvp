package render

import "image"

type Renderer interface {
	Render(img image.Image) string
}
