// Package sdl provides an SDL2 renderer for Doodle.
package sdl

import (
	"git.kirsle.net/go/render"
	"github.com/veandco/go-sdl2/sdl"
)

// Clear the canvas and set this color.
func (r *Renderer) Clear(color render.Color) {
	if color != r.lastColor {
		r.renderer.SetDrawColor(color.Red, color.Green, color.Blue, color.Alpha)
	}
	r.renderer.Clear()
}

// DrawPoint puts a color at a pixel.
func (r *Renderer) DrawPoint(color render.Color, point render.Point) {
	if color != r.lastColor {
		r.renderer.SetDrawColor(color.Red, color.Green, color.Blue, color.Alpha)
	}
	r.renderer.DrawPoint(int32(point.X), int32(point.Y))
}

// DrawLine draws a line between two points.
func (r *Renderer) DrawLine(color render.Color, a, b render.Point) {
	if color != r.lastColor {
		r.renderer.SetDrawColor(color.Red, color.Green, color.Blue, color.Alpha)
	}
	r.renderer.DrawLine(int32(a.X), int32(a.Y), int32(b.X), int32(b.Y))
}

// DrawRect draws a rectangle.
func (r *Renderer) DrawRect(color render.Color, rect render.Rect) {
	if color != r.lastColor {
		r.renderer.SetDrawColor(color.Red, color.Green, color.Blue, color.Alpha)
	}
	r.renderer.DrawRect(&sdl.Rect{
		X: int32(rect.X),
		Y: int32(rect.Y),
		W: int32(rect.W),
		H: int32(rect.H),
	})
}

// DrawBox draws a filled rectangle.
func (r *Renderer) DrawBox(color render.Color, rect render.Rect) {
	if color != r.lastColor {
		r.renderer.SetDrawColor(color.Red, color.Green, color.Blue, color.Alpha)
	}
	r.renderer.FillRect(&sdl.Rect{
		X: int32(rect.X),
		Y: int32(rect.Y),
		W: int32(rect.W),
		H: int32(rect.H),
	})
}
