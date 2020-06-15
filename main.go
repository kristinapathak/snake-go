package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/viper"
	"golang.org/x/image/colornames"
)

type SnakeConfig struct {
	Board      BoardConfig
	SnakeSpeed time.Duration
}

type BoardConfig struct {
	SquareSize     float64
	NumSquaresWide float64
	NumSquaresHigh float64
	Buffer         float64
	BorderWidth    float64
}

func main() {
	pixelgl.Run(run)
}

func run() {

	// load configuration with viper
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("snake")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read in viper config: %v\n", err.Error())
		os.Exit(1)
	}
	config := new(SnakeConfig)
	err = v.Unmarshal(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal config: %v\n", err.Error())
		os.Exit(1)
	}

	boardWidth := config.Board.SquareSize * config.Board.NumSquaresWide
	boardHeight := config.Board.SquareSize * config.Board.NumSquaresHigh
	windowWidth := boardWidth + config.Board.Buffer*2
	windowHeight := boardHeight + config.Board.Buffer*2
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
	playingBoard := NewPlayingBoard(boardWidth, boardHeight, config.Board.Buffer, config.Board.BorderWidth)

	es := Edges{
		left:   0,
		right:  int(config.Board.NumSquaresWide),
		bottom: 0,
		top:    int(config.Board.NumSquaresHigh),
	}
	// TODO: set up items for the snake to eat

	// set up the snake itself
	snake := NewSnake(nil, es, config.SnakeSpeed, config.Board.SquareSize, config.Board.Buffer, colornames.Darkmagenta)

	// keep running and updating things until the window is closed.
	for !win.Closed() {
		win.Clear(colornames.Mediumaquamarine)
		playingBoard.Draw(win)

		if win.Pressed(pixelgl.KeyLeft) {
			snake.SetDirection(Left)
		}
		if win.Pressed(pixelgl.KeyRight) {
			snake.SetDirection(Right)
		}
		if win.Pressed(pixelgl.KeyDown) {
			snake.SetDirection(Down)
		}
		if win.Pressed(pixelgl.KeyUp) {
			snake.SetDirection(Up)
		}
		snake.Paint().Draw(win)

		win.Update()
	}
	snake.Stop()
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
