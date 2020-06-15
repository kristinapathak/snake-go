package main

import (
	"container/list"
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
	lastDirection Direction
	currDirection Direction
	currSpeed     time.Duration
	locations     *list.List
	currDrawing   *imdraw.IMDraw
	score         int
	lock          sync.RWMutex

	item             tracker
	edges            Edges
	startingPosition location
	squareSize       float64
	buffer           float64
	colorr           color.Color
	shutdown         chan struct{}
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

	middleY := int((e.top-e.bottom)/2) + e.bottom
	middleX := int((e.right-e.left)/2) + e.left

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
	}
	s.reset()
	s.updateDrawing()
	go s.move()
	return s
}

func (s *Snake) SetDirection(d Direction) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// don't let the snake do a 180 turn
	if s.lastDirection == Up && d == Down ||
		s.lastDirection == Down && d == Up ||
		s.lastDirection == Left && d == Right ||
		s.lastDirection == Right && d == Left {
		return
	}
	s.currDirection = d
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
	if s.currDirection == None {
		s.lock.RUnlock()
		return false
	}

	// get new spot on the board based on direction
	h := s.locations.Front().Value.(point)
	newX := h.X()
	newY := h.Y()
	switch s.currDirection {
	case Up:
		newY++
	case Down:
		newY--
	case Left:
		newX--
	case Right:
		newX++
	}

	s.lock.RUnlock()
	s.lock.Lock()
	defer s.lock.Unlock()

	// check that the new spot won't be outside of the game board
	if newY < s.edges.bottom || newY >= s.edges.top || newX < s.edges.left || newX >= s.edges.right {
		s.reset()
		return true
	}

	// check for collisions with itself
	e := s.locations.Front()
	for e != nil {
		l := e.Value.(point)
		if l.X() == newX && l.Y() == newY {
			s.reset()
			return true
		}
		e = e.Next()
	}

	// set current direction to last direction
	s.lastDirection = s.currDirection

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

	return true
}

// game is lost, bring everything back to the beginning
func (s *Snake) reset() {
	s.lastDirection = None
	s.currDirection = None
	s.locations.Init()
	s.locations.PushFront(s.startingPosition)
	s.score = 0
}

func (s *Snake) updateDrawing() {
	s.lock.RLock()
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
	s.lock.RUnlock()
	s.lock.Lock()
	s.currDrawing = newDrawing
	s.lock.Unlock()
}
