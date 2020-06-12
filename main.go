package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	var (
		squareSize     float64
		numSquaresWide float64
		numSquaresHigh float64
		buffer         float64
		boardWidth     float64
		boardHeight    float64
		borderWidth    float64
	)

	// these should be made configurable but for now I'm hardcoding them
	squareSize = 15 // each square should be 15px by 15px
	numSquaresWide = 30
	numSquaresHigh = 20
	buffer = 20 // 20px buffer around the whole window

	boardWidth = squareSize * numSquaresWide
	boardHeight = squareSize * numSquaresHigh
	windowWidth := boardWidth + buffer*2
	windowHeight := boardHeight + buffer*2
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
	// not sure if we want this
	// win.SetSmooth(true)

	// give us a nice background
	borderWidth = 3 //3px border around the playing area
	playingBoard := NewPlayingBoard(boardWidth, boardHeight, buffer, borderWidth)

	// keep running and updating things until the window is closed.
	for !win.Closed() {
		win.Clear(colornames.Mediumaquamarine)
		playingBoard.Draw(win)
		win.Update()
	}
}

// NewPlayingBoard highlights the playing area with a background and border.
func NewPlayingBoard(boardWidth float64, boardHeight float64, buffer float64, borderWidth float64) *imdraw.IMDraw {
	playingBoard := imdraw.New(nil)

	playingBoard.Color = colornames.Black
	playingBoard.EndShape = imdraw.SharpEndShape
	playingBoard.Push(pixel.Vec{X: buffer, Y: buffer}, pixel.Vec{X: buffer + boardWidth, Y: buffer + boardHeight})
	playingBoard.Rectangle(borderWidth * 2) // half the border is inside the rectange and half is outside...very annoying

	playingBoard.Color = colornames.Cornsilk
	playingBoard.EndShape = imdraw.SharpEndShape
	playingBoard.Push(pixel.Vec{X: buffer, Y: buffer}, pixel.Vec{X: buffer + boardWidth, Y: buffer + boardHeight})
	playingBoard.Rectangle(0)

	return playingBoard
}
