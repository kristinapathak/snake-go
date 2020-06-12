package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	// these should be made configurable but for now I'm hardcoding them
	squareSize := 15 // each square should be 15px by 15px
	numSquaresWide := 30
	numSquaresHigh := 20
	buffer := 20 // 20px buffer around the whole window

	windowWidth := float64(squareSize*numSquaresWide + buffer*2)
	windowHeight := float64(squareSize*numSquaresHigh + buffer*2)
	cfg := pixelgl.WindowConfig{
		Title:  "Snake!",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	// Start it up!
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// give us a nice background
	win.Clear(colornames.Mediumaquamarine)

	// keep running and updating things until the window is closed.
	for !win.Closed() {
		win.Update()
	}
}
