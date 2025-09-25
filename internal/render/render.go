package render

type Renderer interface {
	RenderFrame(rgb []byte, width, height int) string
}

func NewASCII() Renderer {
	return &asciiRenderer{}
}
