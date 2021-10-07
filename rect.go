package render

import "fmt"

// Rect has a coordinate and a width and height.
type Rect struct {
	X int
	Y int
	W int
	H int
}

// NewRect creates a rectangle of size `width` and `height`. The X,Y values
// are initialized to zero.
func NewRect(width, height int) Rect {
	return Rect{
		W: width,
		H: height,
	}
}

func (r Rect) String() string {
	return fmt.Sprintf("Rect<%d,%d,%d,%d>",
		r.X, r.Y, r.W, r.H,
	)
}

// Point returns the rectangle's X,Y values as a Point.
func (r Rect) Point() Point {
	return Point{
		X: r.X,
		Y: r.Y,
	}
}

// Bigger returns if the given rect is larger than the current one.
func (r Rect) Bigger(other Rect) bool {
	// TODO: don't know why this is !
	return !(other.X < r.X || // Lefter
		other.Y < r.Y || // Higher
		other.W > r.W || // Wider
		other.H > r.H) // Taller
}

// Intersects with the other rectangle in any way.
func (r Rect) Intersects(other Rect) bool {
	// Do a bidirectional compare.
	compare := func(a, b Rect) bool {
		var corners = []Point{
			NewPoint(b.X, b.Y),
			NewPoint(b.X, b.Y+b.H),
			NewPoint(b.X+b.W, b.Y),
			NewPoint(b.X+b.W, b.Y+b.H),
		}
		for _, pt := range corners {
			if pt.Inside(a) {
				return true
			}
		}
		return false
	}

	return compare(r, other) || compare(other, r) || false
}

// IsZero returns if the Rect is uninitialized.
func (r Rect) IsZero() bool {
	return r.X == 0 && r.Y == 0 && r.W == 0 && r.H == 0
}

// Add another rect.
func (r Rect) Add(other Rect) Rect {
	return Rect{
		X: r.X + other.X,
		Y: r.Y + other.Y,
		W: r.W + other.W,
		H: r.H + other.H,
	}
}

// Add a point to move the rect.
func (r Rect) AddPoint(other Point) Rect {
	return Rect{
		X: r.X + other.X,
		Y: r.Y + other.Y,
		W: r.W,
		H: r.H,
	}
}

// SubtractPoint is the inverse of AddPoint. Use this only if you need to invert
// the Point being added.
//
// This does r.X - other.X, r.Y - other.Y and keeps the width/height the same.
func (r Rect) SubtractPoint(other Point) Rect {
	return Rect{
		X: r.X - other.X,
		Y: r.Y - other.Y,
		W: r.W,
		H: r.H,
	}
}
