package render

import (
	"image"

	"github.com/david22573/go-tvp/internal/term"
)

type Renderer interface {
	// Initialize sets up the renderer with video and terminal info
	Initialize(videoWidth, videoHeight int) error

	// Render converts an image to ASCII art
	Render(img image.Image) string

	// GetDimensions returns the calculated frame dimensions
	GetDimensions() (width, height int)

	// GetPadding returns centering offsets
	GetPadding() (padX, padY int)
}

// BaseRenderer contains common functionality
type BaseRenderer struct {
	frameWidth  int
	frameHeight int
	padX        int
	padY        int
}

func (b *BaseRenderer) GetDimensions() (int, int) {
	return b.frameWidth, b.frameHeight
}

func (b *BaseRenderer) GetPadding() (int, int) {
	return b.padX, b.padY
}

// calculateDimensions is a helper that all renderers can use
func (b *BaseRenderer) calculateDimensions(videoW, videoH int, charAspect float64) error {
	termW, termH, err := term.Size()
	if err != nil {
		termW, termH = 80, 24
	}

	// Scale video to fit terminal (portrait TikTok style)
	frameH := termH
	frameW := int(float64(frameH) * float64(videoW) / float64(videoH) / charAspect)

	if frameW > termW {
		frameW = termW
		frameH = int(float64(frameW) * charAspect * float64(videoH) / float64(videoW))
	}

	// Compute centering offsets
	b.frameWidth = frameW
	b.frameHeight = frameH
	b.padX = (termW - frameW) / 2
	b.padY = (termH - frameH) / 2

	return nil
}
