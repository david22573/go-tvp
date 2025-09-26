package main

import (
	"flag"
	"log"

	"github.com/david22573/go-tvp/internal/player"
	"github.com/david22573/go-tvp/internal/render"
)

func main() {
	videoPath := flag.String("f", "", "Video file to play")
	mode := flag.String("mode", "braille", "Render mode: ascii|braille|block|sixel")
	colorMode := flag.Bool("color", true, "Render video in color (true/false)")
	flag.Parse()

	if *videoPath == "" {
		log.Fatal("No video file provided. Use -f <file>")
	}

	var renderer render.Renderer
	switch *mode {
	case "ascii":
		renderer = render.NewASCIIRenderer()
	case "braille":
		renderer = render.NewBrailleRenderer()
	case "block":
		renderer = render.NewBlockRenderer()
	case "sixel":
		renderer = render.NewSixelRenderer()
	default:
		log.Fatalf("Unknown render mode: %s. Available: ascii, braille, block, sixel", *mode)
	}

	player := player.New(*videoPath, renderer, *colorMode)
	if err := player.Play(); err != nil {
		log.Fatalf("Playback failed: %v", err)
	}
}
