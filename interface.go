package render

import (
	"image"

	"git.kirsle.net/go/render/event"
)

// Engine is the interface for the rendering engine, keeping SDL-specific stuff
// far away from the core of Doodle.
type Engine interface {
	Setup() error

	// Poll for events like keypresses and mouse clicks.
	Poll() (*event.State, error)
	GetTicks() uint32
	WindowSize() (w, h int)

	// Present presents the current state to the screen.
	Present() error

	// Clear the full canvas and set this color.
	Clear(Color)
	SetTitle(string)
	DrawPoint(Color, Point)
	DrawLine(Color, Point, Point)
	DrawRect(Color, Rect)
	DrawBox(Color, Rect)
	DrawText(Text, Point) error
	ComputeTextRect(Text) (Rect, error)

	// Texture caching.
	StoreTexture(name string, img image.Image) (Texturer, error)
	LoadTexture(name string) (Texturer, error)
	Copy(t Texturer, src, dst Rect)

	// Teardown and free memory for all textures, returning the number
	// of textures that were freed.
	FreeTextures() int

	// Delay for a moment using the render engine's delay method,
	// implemented by sdl.Delay(uint32)
	Delay(uint32)

	// Tasks that the Setup function should defer until tear-down.
	Teardown()

	Loop() error // maybe?
}

// Texturer is a stored image texture used by the rendering engine while
// abstracting away its inner workings.
type Texturer interface {
	Size() Rect
	Image() image.Image
	Free() error // teardown and free memory
}
