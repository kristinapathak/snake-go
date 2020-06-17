package main

import (
	"container/list"
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

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
	randomGen *rand.Rand

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
	if float64(int(l.X())) != s.currLocation.X() || float64(int(l.Y())) != s.currLocation.Y() {
		s.lock.RUnlock()
		return false
	}
	s.lock.RUnlock()
	return true
}

func (s *singleTracker) Reset(l *list.List) {
	s1 := rand.NewSource(time.Now().UnixNano())
	s.randomGen = rand.New(s1)
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

func (s *singleTracker) findNewLocation(locations *list.List) location {
	gridX := (s.edges.right - s.edges.left)
	gridY := (s.edges.top - s.edges.bottom)
	// TODO: use the snake location to make sure we don't pick a spot where the
	// snake is.
	fmt.Printf("GridWdith Square X: %f Y: %f\n", gridX, gridY)
	if locations == nil || locations.Len() < 1 {
		return location{
			x: float64(s.randomGen.Intn(int(gridX)-1) + 1),
			y: float64(s.randomGen.Intn(int(gridY)-1) + 1),
		}
	}
	for {

		newLocation := location{
			x: float64(s.randomGen.Intn(int(gridX)-1) + 1),
			y: float64(s.randomGen.Intn(int(gridY)-1) + 1),
		}
		if !pointInList(newLocation, locations) {
			return newLocation
		}
	}
}

func pointInList(point location, locations *list.List) bool {
	root := locations.Front()
	for root != nil {
		if root.Value.(location).Equal(point) {
			return true
		}
		root = root.Next()
	}
	return false
}

func (s *singleTracker) updateDrawing() {
	s.lock.RLock()
	newDrawing := imdraw.New(nil)
	newDrawing.Color = s.colorr
	newDrawing.EndShape = imdraw.SharpEndShape

	floatX := float64(s.currLocation.X())
	floatY := float64(s.currLocation.Y())

	newDrawing.Push(pixel.Vec{X: s.buffer + floatX*s.squareSize, Y: s.buffer + floatY*s.squareSize}, pixel.Vec{X: s.buffer + (floatX)*s.squareSize + s.squareSize, Y: s.buffer + (floatY)*s.squareSize + s.squareSize})
	newDrawing.Rectangle(0)

	s.lock.RUnlock()
	s.lock.Lock()
	s.currDrawing = newDrawing
	s.lock.Unlock()
}
