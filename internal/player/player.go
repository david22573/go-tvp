package player

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
	"time"

	"github.com/david22573/go-tvp/internal/render"
)

type Player struct {
	videoPath string
	renderer  render.Renderer
	color     bool
}

func New(videoPath string, renderer render.Renderer, colorMode bool) *Player {
	return &Player{videoPath: videoPath, renderer: renderer, color: colorMode}
}

func (p *Player) Play() error {
	// Detect video properties
	videoFPS, err := DetectFPS(p.videoPath)
	if err != nil {
		fmt.Printf("Warning: Could not detect FPS (%v), using 24fps\n", err)
		videoFPS = 24.0
	}
	delay := time.Duration(float64(time.Second) / videoFPS)

	videoW, videoH, err := DetectResolution(p.videoPath)
	if err != nil {
		fmt.Printf("Warning: Could not detect resolution (%v), defaulting 1920x1080\n", err)
		videoW, videoH = 1920, 1080
	}

	// Initialize renderer with video dimensions
	if err := p.renderer.Initialize(videoW, videoH); err != nil {
		return fmt.Errorf("failed to initialize renderer: %w", err)
	}

	frameW, frameH := p.renderer.GetDimensions()
	padX, padY := p.renderer.GetPadding()

	fmt.Printf("Playing at %.2f FPS\n", videoFPS)
	fmt.Printf("Frame: %dx%d, Padding: %dx%d\n", frameW, frameH, padX, padY)

	// Prepare FFmpeg command
	cmd := prepareFFmpegCommand(p.videoPath, frameW, frameH)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	defer cmd.Wait()

	frameSize := frameW * frameH * 3
	buf := make([]byte, frameSize)

	// Hide cursor during playback
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// Clear screen once
	fmt.Print("\033[2J")

	for {
		start := time.Now()

		// Read frame data
		_, err := io.ReadFull(stdout, buf)
		if err != nil {
			if err == io.EOF {
				break // Normal end of video
			}
			return fmt.Errorf("failed to read frame: %w", err)
		}

		// Convert raw RGB to image
		img := RGBToImage(buf, frameW, frameH)

		// Convert to grayscale if needed
		if !p.color {
			img = convertToGray(img)
		}

		// Render frame
		output := p.renderer.Render(img)

		// Center and display
		centeredOutput := p.centerOutput(output, padX, padY)
		fmt.Print("\033[H") // move cursor to top-left
		fmt.Print(centeredOutput)

		// Frame rate control
		elapsed := time.Since(start)
		if elapsed < delay {
			time.Sleep(delay - elapsed)
		}
	}

	return nil
}

// centerOutput centers the rendered frame on the terminal
func (p *Player) centerOutput(frame string, padX, padY int) string {
	lines := strings.Split(strings.TrimSuffix(frame, "\n"), "\n")

	// Create padding strings
	leftPad := strings.Repeat(" ", padX)
	topPad := strings.Repeat("\n", padY)

	// Apply padding to each line
	var centeredLines []string
	for _, line := range lines {
		centeredLines = append(centeredLines, leftPad+line)
	}

	return topPad + strings.Join(centeredLines, "\n") + "\n"
}

// convertToGray converts an image.RGBA to grayscale
func convertToGray(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			out.Set(x, y, c)
		}
	}

	return out
}
