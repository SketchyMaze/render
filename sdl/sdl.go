// Package sdl provides an SDL2 renderer for Doodle.
package sdl

import (
	"time"

	"git.kirsle.net/apps/doodle/lib/events"
	"git.kirsle.net/apps/doodle/lib/render"
	"git.kirsle.net/apps/doodle/pkg/log"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Renderer manages the SDL state.
type Renderer struct {
	// Configurable fields.
	title     string
	width     int32
	height    int32
	startTime time.Time

	// Private fields.
	events   *events.State
	window   *sdl.Window
	renderer *sdl.Renderer
	running  bool
	ticks    uint64

	// Optimizations to minimize SDL calls.
	lastColor render.Color
}

// New creates the SDL renderer.
func New(title string, width, height int) *Renderer {
	return &Renderer{
		events: events.New(),
		title:  title,
		width:  int32(width),
		height: int32(height),
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
	log.Info("Initializing SDL")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	// Initialize SDL_TTF.
	log.Info("Initializing SDL_TTF")
	if err := ttf.Init(); err != nil {
		return err
	}

	// Create our window.
	log.Info("Creating the Main Window")
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
	log.Info("Creating the SDL Renderer")
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.renderer = renderer

	return nil
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