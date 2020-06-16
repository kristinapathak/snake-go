package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kristinaspring/snake-go/gameloop"

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

	gridBoard := drawGrid(config.Board.NumSquaresWide, config.Board.NumSquaresHigh, config.Board.Buffer, config.Board.SquareSize)

	es := Edges{
		left:   0,
		right:  config.Board.NumSquaresWide,
		bottom: 0,
		top:    config.Board.NumSquaresHigh,
	}
	// set up items for the snake to eat
	tracker := NewSingleTracker(es, config.Board.SquareSize, config.Board.Buffer, colornames.Greenyellow)

	// set up the snake itself
	snake := NewSnake(tracker, es, config.SnakeSpeed, config.Board.SquareSize, config.Board.Buffer, colornames.Darkmagenta)

	g := Game{
		playingBoard: playingBoard,
		tracker:      tracker,
		gridBoard:    gridBoard,
		window:       win,
	}
	stopChan := gameloop.StartLoop(g, time.Second/60, snake)

	// keep running and updating things until the window is closed.
	for !win.Closed() {
	}
	stopChan <- struct{}{}
	snake.Stop()
}

type Game struct {
	playingBoard *imdraw.IMDraw
	gridBoard    *imdraw.IMDraw
	tracker      tracker
	window       *pixelgl.Window
}

func (g Game) Integrate(currentState interface{}, t float64, deltaT float64) interface{} {
	snake := currentState.(*Snake)
	if g.window.Pressed(pixelgl.KeyLeft) {
		snake.SetDirection(Left)
	}
	if g.window.Pressed(pixelgl.KeyRight) {
		snake.SetDirection(Right)
	}
	if g.window.Pressed(pixelgl.KeyDown) {
		snake.SetDirection(Down)
	}
	if g.window.Pressed(pixelgl.KeyUp) {
		snake.SetDirection(Up)
	}
	snake.Tick(t, deltaT)
	return snake
}

func (g Game) Render(state interface{}, t float64, alpha float64) {
	g.window.Clear(colornames.Mediumaquamarine)
	g.playingBoard.Draw(g.window)
	g.gridBoard.Draw(g.window)
	g.tracker.Paint().Draw(g.window)
	snake := state.(*Snake)
	snake.Paint().Draw(g.window)

	g.window.Update()
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

func drawGrid(squareWidthCount float64, squareHeightCount float64, buffer float64, squareSize float64) *imdraw.IMDraw {
	fmt.Println(squareWidthCount, squareHeightCount, buffer, squareSize)
	gridBoard := imdraw.New(nil)

	for i := 0; i <= int(squareWidthCount); i++ {
		gridBoard.Color = colornames.Red
		gridBoard.EndShape = imdraw.RoundEndShape
		gridBoard.Push(pixel.V(buffer+float64(i)*squareSize, buffer), pixel.V(buffer+float64(i)*squareSize, buffer+(squareHeightCount*squareSize)))
		width := 1.0
		if i%10 == 0 {
			width = 2.0
		}
		gridBoard.Line(width)
	}
	for j := 0; j <= int(squareHeightCount); j++ {
		gridBoard.Color = colornames.Red
		gridBoard.EndShape = imdraw.RoundEndShape
		Y := buffer + (float64(j) * squareSize)
		start := pixel.V(buffer, Y)
		end := pixel.V(buffer+(squareWidthCount*squareSize), Y)
		gridBoard.Push(start, end)
		width := 1.0
		if j%10 == 0 {
			width = 2.0
		}
		gridBoard.Line(width)
		// fmt.Printf("%#v / %#v\n", start, end)
	}

	gridBoard.Color = colornames.Red
	gridBoard.EndShape = imdraw.SharpEndShape
	gridBoard.Push(pixel.V(0, 0), pixel.V(squareSize, squareSize))
	gridBoard.Rectangle(0.0)
	return gridBoard
}
