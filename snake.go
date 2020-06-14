package main

import (
	"container/list"
	"fmt"
	"image/color"
	"sync"
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
	X() int
	Y() int
}

type location struct {
	x int
	y int
}

func (l location) X() int {
	return l.x
}

func (l location) Y() int {
	return l.y
}

type Edges struct {
	left   int
	right  int
	top    int
	bottom int
}

type Snake struct {
	// these values are changed asynchronously and need a lock.
	direction   Direction
	currSpeed   time.Duration
	locations   *list.List
	currDrawing *imdraw.IMDraw
	lock        sync.RWMutex

	edges            Edges
	startingPosition location
	squareSize       float64
	buffer           float64
	colorr           color.Color
	shutdown         chan struct{}
}

func NewSnake(edges Edges, speed time.Duration, squareSize float64, buffer float64, c color.Color) *Snake {
	e := edges
	if edges.right < edges.left {
		e.right = edges.left
		e.left = edges.right
	}
	if edges.top < edges.bottom {
		e.top = edges.bottom
		e.bottom = edges.top
	}

	middleY := int((e.top-e.bottom)/2) + e.bottom
	middleX := int((e.right-e.left)/2) + e.left

	l := list.New()

	s := &Snake{
		currSpeed:        speed,
		locations:        l,
		edges:            e,
		startingPosition: location{x: middleX, y: middleY},
		squareSize:       squareSize,
		buffer:           buffer,
		colorr:           c,
		shutdown:         make(chan struct{}, 1),
	}
	fmt.Printf("snake made")
	s.reset()
	fmt.Printf("reset snake")
	s.updateDrawing()
	fmt.Printf("created drawing")
	go s.move()
	fmt.Printf("started ticker")
	return s
}

func (s *Snake) SetDirection(d Direction) {
	s.lock.Lock()
	s.direction = d
	s.lock.Unlock()
}

func (s *Snake) Paint() *imdraw.IMDraw {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.currDrawing
}

func (s *Snake) Stop() {
	close(s.shutdown)
}

func (s *Snake) move() {
	s.lock.RLock()
	ticker := time.NewTicker(s.currSpeed)
	s.lock.RUnlock()
	for {
		select {
		case <-ticker.C:
			moved := s.updateLocations()
			if moved {
				s.updateDrawing()
			}
		case <-s.shutdown:
			return
		}
	}

}

func (s *Snake) updateLocations() bool {
	s.lock.RLock()
	if s.direction == None {
		s.lock.RUnlock()
		return false
	}

	s.lock.RUnlock()
	s.lock.Lock()
	defer s.lock.Unlock()

	// get new spot on the board based on direction
	h := s.locations.Front().Value.(point)
	newX := h.X()
	newY := h.Y()
	switch s.direction {
	case Up:
		newY += 1
	case Down:
		newY -= 1
	case Left:
		newX -= 1
	case Right:
		newX += 1
	}

	// check that the new spot won't be outside of the game board
	if newY < s.edges.bottom || newY >= s.edges.top || newX < s.edges.left || newX >= s.edges.right {
		s.reset()
		return true
	}

	// TODO: check for collisions with itself

	// add new item to the list
	s.locations.PushFront(location{x: newX, y: newY})

	// TODO: check if we ate something and if so, don't remove the last item
	// from the list.

	// remove the last item from the list
	s.locations.Remove(s.locations.Back())

	return true
}

// game is lost, bring everything back to the beginning
func (s *Snake) reset() {
	s.direction = None
	s.locations.Init()
	s.locations.PushFront(s.startingPosition)
}

func (s *Snake) updateDrawing() {
	s.lock.Lock()
	newDrawing := imdraw.New(nil)
	newDrawing.Color = s.colorr
	newDrawing.EndShape = imdraw.SharpEndShape

	e := s.locations.Front()
	for e != nil {
		l := e.Value.(point)
		floatX := float64(l.X())
		floatY := float64(l.Y())
		newDrawing.Push(pixel.Vec{X: s.buffer + floatX*s.squareSize, Y: s.buffer + floatY*s.squareSize}, pixel.Vec{X: s.buffer + (floatX+1)*s.squareSize, Y: s.buffer + (floatY+1)*s.squareSize})
		newDrawing.Rectangle(0)
		e = e.Next()
	}
	s.currDrawing = newDrawing
	s.lock.Unlock()
}
