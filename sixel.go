package main

import (
	"fmt"
)

func main() {
	// Start Sixel
	fmt.Print("\033Pq")

	// Print 4x4 blocks with different colors
	colors := [4][4][3]int{
		{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}, {255, 255, 0}},
		{{0, 255, 255}, {255, 0, 255}, {128, 128, 128}, {255, 255, 255}},
		{{255, 128, 0}, {128, 0, 255}, {0, 128, 255}, {128, 255, 0}},
		{{0, 0, 0}, {64, 64, 64}, {128, 128, 128}, {192, 192, 192}},
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			r := colors[y][x][0]
			g := colors[y][x][1]
			b := colors[y][x][2]
			fmt.Printf("#0;2;%d;%d;%d~", r, g, b) // Sixel pixel
		}
		fmt.Print("-") // move to next row
	}

	// End Sixel
	fmt.Print("\033\\\n")
}
