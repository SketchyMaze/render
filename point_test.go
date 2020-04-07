package render_test

import (
	"strconv"
	"testing"

	"git.kirsle.net/go/render"
)

func TestPointInside(t *testing.T) {
	type testCase struct {
		rect       render.Rect
		p          render.Point
		shouldPass bool
	}
	tests := []testCase{
		testCase{
			rect: render.Rect{
				X: 0,
				Y: 0,
				W: 500,
				H: 500,
			},
			p:          render.NewPoint(128, 256),
			shouldPass: true,
		},
		testCase{
			rect: render.Rect{
				X: 100,
				Y: 80,
				W: 40,
				H: 60,
			},
			p:          render.NewPoint(128, 256),
			shouldPass: false,
		},
		testCase{
			// true values when debugging why Doodads weren't
			// considered inside the viewport.
			rect: render.Rect{
				X: 0,
				Y: -232,
				H: 874,
				W: 490,
			},
			p:          render.NewPoint(509, 260),
			shouldPass: false,
		},
	}

	for _, test := range tests {
		if test.p.Inside(test.rect) != test.shouldPass {
			t.Errorf("Failed: %s inside %s should be %s",
				test.p,
				test.rect,
				strconv.FormatBool(test.shouldPass),
			)
		}
	}
}

// Test the Compare function of Point.
func TestPointDelta(t *testing.T) {
	var tests = []struct {
		A render.Point // source
		B render.Point // comparator
		D render.Point // expected delta value
	}{
		{
			A: render.NewPoint(0, 0),
			B: render.NewPoint(10, 10),
			D: render.NewPoint(10, 10),
		},
		{
			A: render.NewPoint(128, 128),
			B: render.NewPoint(128, 128),
			D: render.NewPoint(0, 0),
		},
		{
			A: render.NewPoint(128, 128),
			B: render.NewPoint(127, 129),
			D: render.NewPoint(-1, 1),
		},
		{
			A: render.NewPoint(200, 500),
			B: render.NewPoint(180, 528),
			D: render.NewPoint(-20, 28),
		},
	}

	for _, test := range tests {
		actual := test.A.Compare(test.B)
		if actual != test.D {
			t.Errorf("Failed: (%s).Compare(%s) expected %s but got %s",
				test.A, test.B, test.D, actual,
			)
		}
	}
}
