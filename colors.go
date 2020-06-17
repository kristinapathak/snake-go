package main

import (
	"image/color"
	"strings"

	"golang.org/x/image/colornames"
)

type Colors int
type Style int

const (
	Black Colors = iota
	Grey
	White
	Purple
	Blue
	Green
	Yellow
	Orange
	Red
	Rainbow
)

const (
	Solid Style = iota
	Striped
)

func GetColor(c string) Colors {
	lc := strings.ToLower(c)
	switch lc {
	case "grey":
		return Grey
	case "white":
		return White
	case "purple":
		return Purple
	case "blue":
		return Blue
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "orange":
		return Orange
	case "red":
		return Red
	case "rainbow":
		return Rainbow
	default:
		return Black
	}
}

func GetStyle(s string) Style {
	ls := strings.ToLower(s)
	switch ls {
	case "striped":
		return Striped
	default:
		return Solid
	}
}

func (c Colors) GetColors(s Style) []color.Color {
	switch c {
	case Grey:
		if s == Solid {
			return []color.Color{colornames.Grey}
		}
		return []color.Color{colornames.Darkgrey, colornames.Lightgray}
	case White:
		if s == Solid {
			return []color.Color{colornames.White}
		}
		return []color.Color{colornames.Navajowhite, colornames.White}
	case Purple:
		if s == Solid {
			return []color.Color{colornames.Darkmagenta}
		}
		return []color.Color{colornames.Darkmagenta, colornames.Lavender}
	case Blue:
		if s == Solid {
			return []color.Color{colornames.Blue}
		}
		return []color.Color{colornames.Darkblue, colornames.Deepskyblue}
	case Green:
		if s == Solid {
			return []color.Color{colornames.Green}
		}
		return []color.Color{colornames.Darkolivegreen, colornames.Palegreen}
	case Yellow:
		if s == Solid {
			return []color.Color{colornames.Yellow}
		}
		return []color.Color{colornames.Darkgoldenrod, colornames.Gold}
	case Orange:
		if s == Solid {
			return []color.Color{colornames.Orange}
		}
		return []color.Color{colornames.Darkorange, colornames.Peachpuff}
	case Red:
		if s == Solid {
			return []color.Color{colornames.Red}
		}
		return []color.Color{colornames.Darkred, colornames.Tomato}
	case Rainbow:
		return []color.Color{colornames.Purple, colornames.Blue, colornames.Green, colornames.Yellow, colornames.Orange, colornames.Red}
	default:
		return []color.Color{colornames.Black, colornames.Darkslategray}
	}
}
