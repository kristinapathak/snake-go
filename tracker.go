package main

import (
	"container/list"
	"image/color"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type tracker interface {
	At(location) bool
	Reset(*list.List)
	Paint() *imdraw.IMDraw
}

type defaultTracker struct{}

func (d defaultTracker) At(_ location) bool {
	return false
}

func (d defaultTracker) Reset(_ *list.List) {}

func (d defaultTracker) Paint() *imdraw.IMDraw {
	return imdraw.New(nil)
}

type singleTracker struct {
	currLocation location
	currDrawing  *imdraw.IMDraw
	lock         sync.RWMutex

	edges      Edges
	squareSize float64
	buffer     float64
	colorr     color.Color
}

func NewSingleTracker(edges Edges, squareSize float64, buffer float64, c color.Color) *singleTracker {
	e := edges
	if edges.right < edges.left {
		e.right = edges.left
		e.left = edges.right
	}
	if edges.top < edges.bottom {
		e.top = edges.bottom
		e.bottom = edges.top
	}

	s := singleTracker{
		edges:      edges,
		squareSize: squareSize,
		buffer:     buffer,
		colorr:     c,
	}
	s.Reset(nil)
	return &s
}

func (s *singleTracker) At(l location) bool {
	s.lock.RLock()
	if l.X() != s.currLocation.X() || l.Y() != s.currLocation.Y() {
		s.lock.RUnlock()
		return false
	}
	s.lock.RUnlock()
	return true
}

func (s *singleTracker) Reset(l *list.List) {
	loc := s.findNewLocation(l)
	s.lock.Lock()
	s.currLocation = loc
	s.lock.Unlock()
	s.updateDrawing()
}

func (s *singleTracker) Paint() *imdraw.IMDraw {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.currDrawing
}

func (s *singleTracker) findNewLocation(_ *list.List) location {
	// TODO: use the snake location to make sure we don't pick a spot where the
	// snake is.
	// TODO: Make this random.
	return location{
		x: 3,
		y: 7,
	}
}

func (s *singleTracker) updateDrawing() {
	s.lock.RLock()
	newDrawing := imdraw.New(nil)
	newDrawing.Color = s.colorr
	newDrawing.EndShape = imdraw.SharpEndShape

	floatX := float64(s.currLocation.X())
	floatY := float64(s.currLocation.Y())

	newDrawing.Push(pixel.Vec{X: s.buffer + floatX*s.squareSize, Y: s.buffer + floatY*s.squareSize}, pixel.Vec{X: s.buffer + (floatX+1)*s.squareSize, Y: s.buffer + (floatY+1)*s.squareSize})
	newDrawing.Rectangle(0)

	s.lock.RUnlock()
	s.lock.Lock()
	s.currDrawing = newDrawing
	s.lock.Unlock()
}
