package main

import (
	"image/color"
	"strings"

	"golang.org/x/image/colornames"
)

type Colors int

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

func GetEnum(c string) Colors {
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

func (c Colors) GetColors() []color.Color {
	switch c {
	case Grey:
		return []color.Color{colornames.Darkgrey, colornames.Lightgray}
	case White:
		return []color.Color{colornames.Whitesmoke, colornames.White}
	case Purple:
		return []color.Color{colornames.Darkmagenta, colornames.Lavender}
	case Blue:
		return []color.Color{colornames.Darkslateblue, colornames.Aliceblue}
	case Green:
		return []color.Color{colornames.Darkolivegreen, colornames.Mintcream}
	case Yellow:
		return []color.Color{colornames.Darkgoldenrod, colornames.Lemonchiffon}
	case Orange:
		return []color.Color{colornames.Darkorange, colornames.Peachpuff}
	case Red:
		return []color.Color{colornames.Darkred, colornames.Tomato}
	case Rainbow:
		return []color.Color{colornames.Purple, colornames.Blue, colornames.Green, colornames.Yellow, colornames.Orange, colornames.Red}
	default:
		return []color.Color{colornames.Black, colornames.Darkslategray}
	}
}
