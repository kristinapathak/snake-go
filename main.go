package main

import (
	"fmt"
	"github.com/faiface/pixel/text"
	"github.com/kristinaspring/snake-go/gameloop"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"os"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/viper"
	"golang.org/x/image/colornames"
)

type ViperConfig struct {
	Board BoardConfig
	Snake SnakeViperConfig
}

type SnakeViperConfig struct {
	Color          string
	Speed          float64
	StartingFrames int
	FramesToGrow   int
	Threshold      float64
}

type BoardConfig struct {
	SquareSize     float64
	NumSquaresWide float64
	NumSquaresHigh float64
	Buffer         float64
	BorderWidth    float64
	ShowGrid       bool
	ShowCounters   bool
	TickRate       int
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
	config := new(ViperConfig)
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
		VSync:  false,
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
	if config.Board.ShowGrid {
		drawGrid(playingBoard, config.Board.NumSquaresWide, config.Board.NumSquaresHigh, config.Board.Buffer, config.Board.SquareSize)
	}
	es := Edges{
		left:   0,
		right:  config.Board.NumSquaresWide,
		bottom: 0,
		top:    config.Board.NumSquaresHigh,
	}
	// set up items for the snake to eat
	tracker := NewSingleTracker(es, config.Board.SquareSize, config.Board.Buffer, colornames.Indianred)

	// set up the snake itself
	c := SnakeConfig{
		Edges:          es,
		SquareSize:     config.Board.SquareSize,
		Buffer:         config.Board.Buffer,
		Colors:         GetEnum(config.Snake.Color).GetColors(),
		PixelsPerSec:   config.Snake.Speed,
		StartingFrames: config.Snake.StartingFrames,
		FramesToGrow:   config.Snake.FramesToGrow,
		Threshold:      config.Snake.Threshold,
	}

	snake := NewSnake(tracker, c)

	g := &Game{
		playingBoard: playingBoard,
		tracker:      tracker,
		window:       win,
		frameCount:   NewCounter(100),
		updateCount:  NewCounter(100),
	}
	if config.Board.ShowCounters {
		g.txt = text.New(pixel.V(1, 1), text.NewAtlas(
			ttfFromBytesMust(goregular.TTF, config.Board.Buffer),
			text.ASCII, text.RangeTable(unicode.Latin),
		))
	}
	stopChan := gameloop.StartLoop(g, time.Second/time.Duration(config.Board.TickRate), snake)

	// keep running and updating things until the window is closed.
	for !win.Closed() {
	}
	stopChan <- struct{}{}
	snake.Stop()
}

type Game struct {
	playingBoard *imdraw.IMDraw
	tracker      tracker
	window       *pixelgl.Window
	measurement  float64
	frameCount   *Counter
	updateCount  *Counter
	txt          *text.Text
}

func (g *Game) Integrate(currentState interface{}, t float64, deltaT float64) interface{} {
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
	if g.txt != nil {
		g.updateCount.Tick(t)
	}
	snake.Tick(t, deltaT)
	return snake
}
func ttfFromBytesMust(b []byte, size float64) font.Face {
	ttf, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	return truetype.NewFace(ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}

func (g *Game) Render(state interface{}, t float64, alpha float64) {
	g.window.Clear(colornames.Mediumaquamarine)

	g.playingBoard.Draw(g.window)

	g.tracker.Paint().Draw(g.window)
	snake := state.(*Snake)
	snake.Paint().Draw(g.window)
	if g.txt != nil {
		g.frameCount.Tick(t)
		g.txt.Clear()
		g.txt.WriteString(fmt.Sprintf("FPS :%4.2f, UPS: %4.2f", g.frameCount.GetRate(), g.updateCount.GetRate()))
		g.txt.Draw(g.window, pixel.IM)
	}
	g.window.Update()
}

type Counter struct {
	maxSamples int
	tickIndex  int
	tickSum    float64
	tickList   []float64
	lastTime   float64

	rate float64
}

func NewCounter(maxSamples int) *Counter {
	return &Counter{
		maxSamples: maxSamples,
		tickList:   make([]float64, maxSamples),
	}
}

func (c *Counter) Tick(currentT float64) float64 {
	delta := currentT - c.lastTime
	c.lastTime = currentT

	c.tickSum -= c.tickList[c.tickIndex]
	c.tickSum += delta
	c.tickList[c.tickIndex] = delta
	// c.tickIndex++
	// if c.tickIndex >= c.maxSamples {
	// 	c.tickIndex = 0
	// }
	c.tickIndex = (c.tickIndex + 1) % c.maxSamples
	c.rate = float64(c.maxSamples) / c.tickSum
	return c.rate
}

func (c *Counter) GetRate() float64 {
	return c.rate
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

func drawGrid(board *imdraw.IMDraw, squareWidthCount float64, squareHeightCount float64, buffer float64, squareSize float64) {
	fmt.Println(squareWidthCount, squareHeightCount, buffer, squareSize)

	for i := 0; i <= int(squareWidthCount); i++ {
		board.Color = colornames.Red
		board.EndShape = imdraw.RoundEndShape
		board.Push(pixel.V(buffer+float64(i)*squareSize, buffer), pixel.V(buffer+float64(i)*squareSize, buffer+(squareHeightCount*squareSize)))
		width := 2.0
		if i%int(squareSize) == 0 {
			width = 3.0
		}
		board.Line(width)
	}
	for j := 0; j <= int(squareHeightCount); j++ {
		board.Color = colornames.Red
		board.EndShape = imdraw.RoundEndShape
		Y := buffer + (float64(j) * squareSize)
		start := pixel.V(buffer, Y)
		end := pixel.V(buffer+(squareWidthCount*squareSize), Y)
		board.Push(start, end)
		width := 2.0
		if j%int(squareSize) == 0 {
			width = 3.0
		}
		board.Line(width)
		// fmt.Printf("%#v / %#v\n", start, end)
	}

}
