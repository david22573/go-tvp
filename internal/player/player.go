package player

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"
	"time"

	"github.com/david22573/go-tvp/internal/render"
	"github.com/david22573/go-tvp/internal/term"
)

type Player struct {
	videoPath string
	renderer  render.Renderer
}

func New(videoPath string, r render.Renderer) *Player {
	return &Player{videoPath: videoPath, renderer: r}
}

func (p *Player) Play() error {
	// Get actual video FPS
	videoFPS, err := DetectFPS(p.videoPath)
	if err != nil {
		fmt.Printf("Warning: Could not detect FPS (%v), using 24fps\n", err)
		videoFPS = 24.0
	}

	fmt.Printf("Playing at %.2f FPS\n", videoFPS)

	// Get terminal size - this will be our target ASCII dimensions
	termW, termH, err := term.Size()
	if err != nil {
		termW, termH = 80, 24 // fallback
	}

	// Scale terminal size for better aspect ratio (characters are taller than wide)
	// Typical character aspect ratio is about 1:2, so we adjust width
	aspectRatio := 2.0
	videoW := int(float64(termW) * aspectRatio)
	videoH := termH

	fmt.Printf("Terminal: %dx%d, Video output: %dx%d\n", termW, termH, videoW, videoH)

	cmd := exec.Command("ffmpeg",
		"-i", p.videoPath,
		"-f", "image2pipe",
		"-pix_fmt", "rgb24",
		"-vcodec", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", videoW, videoH), // Dynamic size based on terminal
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	frameSize := videoW * videoH * 3
	buf := make([]byte, frameSize)

	// Create renderer with terminal dimensions
	renderer := render.NewASCIIRenderer(termW, termH)

	// Hide cursor during playback
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	delay := time.Duration(float64(time.Second) / videoFPS)

	for {
		frameStart := time.Now()

		_, err := io.ReadFull(stdout, buf)
		if err != nil {
			break
		}

		// Convert buf → image.RGBA
		img := rgbToImage(buf, videoW, videoH)

		// Render with the chosen renderer
		out := renderer.Render(img)

		// Move cursor to top-left (don't clear entire screen each frame)
		fmt.Print("\033[H")
		fmt.Print(out)

		// Calculate time spent rendering and adjust delay
		renderTime := time.Since(frameStart)
		if renderTime < delay {
			time.Sleep(delay - renderTime)
		}
		// If rendering took longer than expected delay, skip sleep to catch up
	}

	return cmd.Wait()
}

// rgbToImage converts raw RGB bytes into an *image.RGBA
func rgbToImage(buf []byte, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	idx := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if idx+2 < len(buf) {
				r := buf[idx]
				g := buf[idx+1]
				b := buf[idx+2]
				img.Set(x, y, color.RGBA{r, g, b, 255})
				idx += 3
			}
		}
	}
	return img
}
