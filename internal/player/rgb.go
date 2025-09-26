package player

import (
	"fmt"
	"image"
	"image/color"
	"os/exec"
)

// RGBToImage converts raw RGB bytes to *image.RGBA
func RGBToImage(buf []byte, w, h int) *image.RGBA {
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

// prepareFFmpegCommand returns ffmpeg command
func prepareFFmpegCommand(videoPath string, w, h int) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", videoPath,
		"-f", "image2pipe",
		"-pix_fmt", "rgb24",
		"-vcodec", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", w, h),
		"-",
	)
}
