package player

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/david22573/go-tvp/internal/render"
)

type Player struct {
	videoPath string
	renderer  render.Renderer
}

func New(videoPath string, r render.Renderer) *Player {
	return &Player{videoPath: videoPath, renderer: r}
}

func (p *Player) Play() error {
	cmd := exec.Command("ffmpeg", "-i", p.videoPath, "-f", "image2pipe", "-pix_fmt", "rgb24", "-vcodec", "rawvideo", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// TODO: make width/height dynamic
	width, height := 80, 40
	frameSize := width * height * 3
	buf := make([]byte, frameSize)

	for {
		_, err := stdout.Read(buf)
		if err != nil {
			break
		}
		out := p.renderer.RenderFrame(buf, width, height)
		fmt.Print("\033[H") // move cursor top-left
		fmt.Print(out)

		time.Sleep(time.Second / 24) // naive FPS
	}

	return cmd.Wait()
}
