// +build js,wasm

// This is a version of the 'hello-world' example but build for WebAssembly.

package main

import (
	"image/png"
	"log"
	"net/http"
	"time"

	"git.kirsle.net/go/render"
	"git.kirsle.net/go/render/canvas"
)

var (
	// Cap animation speed to 60 FPS
	targetFPS = 1000 / 60

	// Gopher sprite variables
	gopher   render.Texturer
	texSize  render.Rect
	speed    = 4
	position = render.NewPoint(0, 0)
	velocity = render.NewPoint(speed, speed)

	// Decorative border variables
	borderColor = render.Red
	borderSize  = 12

	// Background color of the window.
	bgColor = render.RGBA(255, 255, 128, 255)
)

func main() {
	// Parameter to New() is the HTML Canvas ID.
	mw, err := canvas.New("canvas")
	if err != nil {
		panic(err)
	}
	setup(mw)

	for {
		start := time.Now()

		ev, err := mw.Poll()
		if err != nil {
			panic(err)
		}

		if ev.Escape {
			mw.Teardown()
			break
		}

		update(mw)
		draw(mw)
		mw.Present()

		// Delay to maintain constant 60 FPS.
		var delay uint32
		elapsed := time.Now().Sub(start)
		tmp := elapsed / time.Millisecond
		if targetFPS-int(tmp) > 0 {
			delay = uint32(targetFPS - int(tmp))
		}
		mw.Delay(delay)
	}
}

func setup(e *canvas.Engine) {
	if err := e.Setup(); err != nil {
		panic(err)
	}

	// Bind DOM event handlers.
	e.AddEventListeners()

	// Load gopher sprite via ajax request.
	resp, err := http.Get("gopher.png")
	if err != nil {
		log.Fatalf("GET gopher.png: %s", err)
	}
	defer resp.Body.Close()

	img, _ := png.Decode(resp.Body)

	gopher, _ = e.StoreTexture("gopher.png", img)
	texSize = gopher.Size()
}

func update(e render.Engine) {
	position.X += velocity.X
	position.Y += velocity.Y

	// Bounce off the walls.
	w, h := e.WindowSize()

	if velocity.X > 0 && position.X+texSize.W >= w-borderSize {
		velocity.X *= -1
	} else if velocity.X < 0 && position.X <= borderSize {
		velocity.X *= -1
	}

	if velocity.Y > 0 && position.Y+texSize.H >= h-borderSize {
		velocity.Y *= -1
	} else if velocity.Y < 0 && position.Y <= borderSize {
		velocity.Y *= -1
	}
}

func draw(e render.Engine) {
	w, h := e.WindowSize()

	drawBorder(e, w, h)

	// Draw some text centered along the top of the canvas.
	text := render.Text{
		Text:         "Hello, world!",
		Size:         24,
		Color:        render.SkyBlue,
		Shadow:       render.Blue,
		FontFilename: "DejaVuSans.ttf",
	}
	rect, _ := e.ComputeTextRect(text)
	e.DrawText(text, render.NewPoint(
		(w/2)-(rect.W/2),
		25,
	))

	e.Copy(gopher, texSize, render.Rect{
		X: position.X,
		Y: position.Y,
		W: texSize.W,
		H: texSize.H,
	})
}

func drawBorder(e render.Engine, w, h int) {
	// Draw the decorative border. We're going for a "ridged" border
	// style here. First draw the light and dark edges of the top/left
	// sides of the border.
	e.DrawBox(borderColor.Lighten(40), render.Rect{
		X: 0,
		Y: 0,
		W: w,
		H: h,
	})
	e.DrawBox(borderColor.Darken(40), render.Rect{
		X: borderSize / 2,
		Y: borderSize / 2,
		W: w - (borderSize / 2),
		H: h - (borderSize / 2),
	})

	// Now inset a bit and draw the light/dark edges of the bottom/right.
	e.DrawBox(borderColor.Darken(40), render.Rect{
		X: borderSize,
		Y: borderSize,
		W: w,
		H: h,
	})
	e.DrawBox(borderColor.Lighten(40), render.Rect{
		X: borderSize,
		Y: borderSize,
		W: w - borderSize - (borderSize / 2),
		H: h - borderSize - (borderSize / 2),
	})

	e.DrawBox(bgColor, render.Rect{
		X: borderSize,
		Y: borderSize,
		W: w - (borderSize * 2),
		H: h - (borderSize * 2),
	})
}
