package main

type tracker interface {
	At(location) bool
}

type defaultTracker struct{}

func (d defaultTracker) At(_ location) bool {
	return false
}
