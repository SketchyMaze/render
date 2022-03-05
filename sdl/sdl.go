// Package sdl provides an SDL2 renderer for Doodle.
package sdl

import (
	"bytes"
	"fmt"
	"image"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/event"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"golang.org/x/image/bmp"
)

// Renderer manages the SDL state.
type Renderer struct {
	// Configurable fields.
	title  string
	width  int32
	height int32

	// Private fields.
	events   *event.State
	window   *sdl.Window
	renderer *sdl.Renderer
	textures map[string]*Texture // cached textures

	// Optimizations to minimize SDL calls.
	lastColor render.Color
}

// New creates the SDL renderer.
func New(title string, width, height int) *Renderer {
	return &Renderer{
		events:   event.NewState(),
		title:    title,
		width:    int32(width),
		height:   int32(height),
		textures: map[string]*Texture{},
	}
}

// Teardown tasks when exiting the program.
func (r *Renderer) Teardown() {
	r.renderer.Destroy()
	r.window.Destroy()
	sdl.Quit()
}

// Setup the renderer.
func (r *Renderer) Setup() error {
	// Initialize SDL.
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("sdl.Init: %s", err)
	}

	// Initialize SDL_TTF.
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("ttf.Init: %s", err)
	}

	// Create our window.
	window, err := sdl.CreateWindow(
		r.title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		r.width,
		r.height,
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return err
	}
	r.window = window

	// Blank out the window in white.
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.renderer = renderer

	return nil
}

// SetWindowIcon sets an icon for the SDL2 window.
func (r *Renderer) SetWindowIcon(icon image.Image) error {
	var (
		fh = bytes.NewBuffer([]byte{}) // bitmap icon buffer
	)

	// Encode icon image to bitmap
	err := bmp.Encode(fh, icon)
	if err != nil {
		return err
	}

	// Create an SDL2 RWOps from the bitmap data in memory.
	rw, err := sdl.RWFromMem(fh.Bytes())
	if err != nil {
		return err
	}

	// Surface from RWOps
	surface, err := sdl.LoadBMPRW(rw, true)
	if err != nil {
		return err
	}
	defer surface.Free()

	// Set the app window icon.
	r.window.SetIcon(surface)

	return nil
}

// SetTitle sets the SDL window title.
func (r *Renderer) SetTitle(title string) {
	r.title = title
	r.window.SetTitle(title)
}

// GetTicks gets SDL's current tick count.
func (r *Renderer) GetTicks() uint32 {
	return sdl.GetTicks()
}

// WindowSize returns the SDL window size.
func (r *Renderer) WindowSize() (int, int) {
	w, h := r.window.GetSize()
	return int(w), int(h)
}

// Maximize the window.
func (r *Renderer) Maximize() {
	r.window.Maximize()
}

// Present the current frame.
func (r *Renderer) Present() error {
	r.renderer.Present()
	return nil
}

// Delay using sdl.Delay
func (r *Renderer) Delay(time uint32) {
	sdl.Delay(time)
}

// Loop is the main loop.
func (r *Renderer) Loop() error {
	return nil
}

// Private access methods
func (r *Renderer) GetSDL2Renderer() *sdl.Renderer {
	return r.renderer
}
