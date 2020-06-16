package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Direction int

const (
	None Direction = iota
	Up
	Down
	Right
	Left
)

type point interface {
	X() float64
	Y() float64
}

type location struct {
	x float64
	y float64
}

func (l location) X() float64 {
	return l.x
}

func (l location) Y() float64 {
	return l.y
}

type Edges struct {
	left   float64
	right  float64
	top    float64
	bottom float64
}

type Snake struct {
	// these values are changed asynchronously and need a lock.
	lastDirection Direction
	currDirection Direction
	nextDirection Direction
	currSpeed     time.Duration
	locations     *list.List
	currDrawing   *imdraw.IMDraw
	score         int

	item             tracker
	edges            Edges
	startingPosition location
	squareSize       float64
	buffer           float64
	colorr           color.Color
	shutdown         chan struct{}

	pixelsPerSec float64
}

func NewSnake(itemTracker tracker, edges Edges, speed time.Duration, squareSize float64, buffer float64, c color.Color) *Snake {
	e := edges
	if edges.right < edges.left {
		e.right = edges.left
		e.left = edges.right
	}
	if edges.top < edges.bottom {
		e.top = edges.bottom
		e.bottom = edges.top
	}

	middleY := (e.top-e.bottom)/2.0 + e.bottom
	middleX := (e.right-e.left)/2.0 + e.left

	l := list.New()

	item := itemTracker
	if item == nil {
		item = defaultTracker{}
	}

	s := &Snake{
		currSpeed:        speed,
		locations:        l,
		item:             item,
		edges:            e,
		startingPosition: location{x: middleX, y: middleY},
		squareSize:       squareSize,
		buffer:           buffer,
		colorr:           c,
		shutdown:         make(chan struct{}, 1),
		pixelsPerSec:     20,
	}
	s.reset()
	// s.updateDrawing()
	s.currDirection = None
	// go s.move()
	return s
}

func (s *Snake) SetDirection(d Direction) {
	// don't let the snake do a 180 turn
	if s.currDirection == Up && d == Down ||
		s.currDirection == Down && d == Up ||
		s.currDirection == Left && d == Right ||
		s.currDirection == Right && d == Left {
		fmt.Println("Can't do")
		return
	}
	s.nextDirection = d
}

func (s *Snake) Paint() *imdraw.IMDraw {
	newDrawing := imdraw.New(nil)
	newDrawing.Color = s.colorr
	newDrawing.EndShape = imdraw.SharpEndShape

	e := s.locations.Front()
	for e != nil {
		l := e.Value.(point)

		newDrawing.Push(pixel.Vec{X: s.buffer + l.X()*s.squareSize, Y: s.buffer + l.Y()*s.squareSize}, pixel.Vec{X: s.buffer + (l.X() * s.squareSize) + s.squareSize, Y: s.buffer + (l.Y() * s.squareSize) + s.squareSize})
		newDrawing.Rectangle(0)
		e = e.Next()
	}
	return newDrawing
}

func (s *Snake) Stop() {
	close(s.shutdown)
}

// func (s *Snake) move() {
// 	s.lock.RLock()
// 	ticker := time.NewTicker(s.currSpeed)
// 	s.lock.RUnlock()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			moved := s.updateLocations()
// 			if moved {
// 				s.updateDrawing()
// 			}
// 		case <-s.shutdown:
// 			return
// 		}
// 	}
// }

func (s *Snake) Tick(t float64, deltaT float64) {
	h := s.locations.Front().Value.(point)
	newX := h.X()
	newY := h.Y()

	switch s.currDirection {
	case Up:
		newY = h.Y() + s.pixelsPerSec*deltaT
	case Down:
		newY = h.Y() - s.pixelsPerSec*deltaT
	case Left:
		newX = h.X() - s.pixelsPerSec*deltaT
	case Right:
		newX = h.X() + s.pixelsPerSec*deltaT
	}

	// TODO:// Check if movemend direction should change
	// roundX := math.Round(newX)
	// roundY := math.Round(newY)
	xCheck := math.Mod(newX, s.squareSize)
	yCheck := math.Mod(newY, s.squareSize)
	// fmt.Printf("Update Direction Tick X: %f Y: %f -> DeltaT %f | RoundX: %f RoundY: %f SquareSize: %f \n", h.X(), h.Y(), deltaT, math.Abs(s.squareSize-xCheck), math.Abs(s.squareSize-yCheck), s.squareSize)
	if xCheck < 0.01 {
		if s.nextDirection != None {
			s.currDirection = s.nextDirection
			s.nextDirection = None
			newX = math.Round(newX)
			newY = math.Round(newY)
		}

	} else if yCheck < 0.01 {
		// fmt.Printf("Update Direction Tick X: %f Y: %f -> DeltaT %f | RoundX: %f RoundY: %f SquareSize: %f \n", h.X(), h.Y(), deltaT, xCheck, yCheck, s.squareSize)
		if s.nextDirection != None {
			s.currDirection = s.nextDirection
			s.nextDirection = None
			newX = math.Round(newX)
			newY = math.Round(newY)
		}

	} else {
		// fmt.Printf("Tick X: %f Y: %f -> DeltaT %f | RoundX: %f RoundY: %f SquareSize: %f \n", h.X(), h.Y(), deltaT, xCheck, yCheck, s.squareSize)
	}

	// check that the new spot won't be outside of the game board
	if newY < s.edges.bottom || newY >= s.edges.top || newX < s.edges.left || newX >= s.edges.right {
		s.reset()
		fmt.Println("died")
		return
	}

	// TODO:// check for collisions with itself
	// e := s.locations.Front()
	// for e != nil {
	// 	l := e.Value.(point)
	// 	if l.X() == newX && l.Y() == newY {
	// 		s.reset()
	// 		fmt.Println("killed self")
	// 		return
	// 	}
	// 	e = e.Next()
	// }

	// add new item to the list
	newSquare := location{x: newX, y: newY}
	s.locations.PushFront(newSquare)

	// check if we ate something and if so, don't remove the last item from the
	// list.
	if s.item.At(newSquare) {
		s.item.Reset(s.locations)
		s.score++
	} else {
		// remove the last item from the list
		s.locations.Remove(s.locations.Back())
	}
}

// game is lost, bring everything back to the beginning
func (s *Snake) reset() {
	s.lastDirection = None
	s.currDirection = None
	s.nextDirection = None
	s.locations.Init()
	s.locations.PushFront(s.startingPosition)
	s.score = 0
}

// func (s *Snake) updateDrawing() {
// 	s.lock.RLock()
// 	newDrawing := imdraw.New(nil)
// 	newDrawing.Color = s.colorr
// 	newDrawing.EndShape = imdraw.SharpEndShape
//
// 	e := s.locations.Front()
// 	for e != nil {
// 		l := e.Value.(point)
// 		floatX := float64(l.X())
// 		floatY := float64(l.Y())
// 		newDrawing.Push(pixel.Vec{X: s.buffer + floatX*s.squareSize, Y: s.buffer + floatY*s.squareSize}, pixel.Vec{X: s.buffer + (floatX+1)*s.squareSize, Y: s.buffer + (floatY+1)*s.squareSize})
// 		newDrawing.Rectangle(0)
// 		e = e.Next()
// 	}
// 	s.lock.RUnlock()
// 	s.lock.Lock()
// 	s.currDrawing = newDrawing
// 	s.lock.Unlock()
// }
