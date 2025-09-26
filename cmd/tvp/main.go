package main

import (
	"flag"
	"log"

	"github.com/david22573/go-tvp/internal/player"
	"github.com/david22573/go-tvp/internal/render"
	"github.com/david22573/go-tvp/internal/term"
)

func main() {
	videoPath := flag.String("f", "", "Video file to play")
	mode := flag.String("mode", "ascii", "Render mode: ascii|braille")
	flag.Parse()

	if *videoPath == "" {
		log.Fatal("No video file provided. Use -f <file>")
	}

	var renderer render.Renderer

	h, w, err := term.Size()
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case "ascii":
		renderer = render.NewASCIIRenderer(h, w)
	default:
		log.Fatalf("Unknown render mode: %s", *mode)
	}

	p := player.New(*videoPath, renderer)
	if err := p.Play(); err != nil {
		log.Fatal(err)
	}
}
