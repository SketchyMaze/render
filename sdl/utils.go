package sdl

import (
	"git.kirsle.net/go/render"
	"github.com/veandco/go-sdl2/sdl"
)

// ColorToSDL converts Doodle's Color type to an sdl.Color.
func ColorToSDL(c render.Color) sdl.Color {
	return sdl.Color{
		R: c.Red,
		G: c.Green,
		B: c.Blue,
		A: c.Alpha,
	}
}

// RectToSDL converts Doodle's Rect type to an sdl.Rect.
func RectToSDL(r render.Rect) sdl.Rect {
	return sdl.Rect{
		X: int32(r.X),
		Y: int32(r.Y),
		W: int32(r.W),
		H: int32(r.H),
	}
}
