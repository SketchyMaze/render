# Render: Go Graphics Library

[![GoDoc](https://godoc.org/git.kirsle.net/go/render?status.svg)](https://godoc.org/git.kirsle.net/go/render)

Render is a graphics rendering library for Go.

It supports SDL2 and HTML Canvas back-ends enabling its use for both desktop
applications (Linux, Mac and Windows) and WebAssembly modules for running in
the web browser.

![Screenshot](examples/hello-world/screenshot.png)

**Notice:** [github.com/SketchyMaze/render](https://github.com/SketchyMaze/render) is a
mirror to the upstream repository at [git.kirsle.net/go/render](https://git.kirsle.net/go/render).
Issues and pull requests will be accepted at GitHub.

## Example

```go
package main

import (
    "git.kirsle.net/go/render"
    "git.kirsle.net/go/render/sdl"
)

func main() {
    mw := sdl.New("Hello World", 320, 240)

    if err := mw.Setup(); err != nil {
        panic(err)
    }

    // Text that we're gonna draw in the window.
    text := render.Text{
        Text:         "Hello, world!",
        Size:         24,
        Color:        render.SkyBlue,
        Shadow:       render.Blue,
        FontFilename: "DejaVuSans.ttf",
    }

    // Compute the rendered size of the text.
    rect, _ := mw.ComputeTextRect(text)

    for {
        // Blank the window.
        mw.Clear(render.White)

        // Poll for events (mouse clicks, keyboard keys, etc.)
        ev, err := mw.Poll()
        if err != nil {
            panic(err)
        }

        // Escape key closes the window.
        if ev.Escape {
            mw.Teardown()
            break
        }

        // Get the window size.
        w, h := mw.WindowSize()

        // Draw the text centered in the window.
        mw.DrawText(text, render.NewPoint(
            (w/2) - (rect.W/2),
            (h/2) - (rect.H/2),
        ))

        mw.Present()
    }
}
```

See the `examples/` directory for examples. More will come eventually,
including some WebAssembly examples.

## Project Status: Alpha

This module was written as part of my drawing-based maze game, code named
[Project: Doodle](https://www.kirsle.net/doodle). It is currently in
**alpha status** and its API may change and be cleaned up in the future.

## Drawing Methods (Engine)

This package provides some _basic_ primitive drawing methods which are
implemented for SDL2 (desktops) and HTML Canvas (WebAssembly). See the
render.Engine interface. The drawing methods supported are:

* Clear(Color): blank the window and fill it with this color.
* DrawPoint(Color, Point): draw a single pixel at a coordinate.
* DrawLine(Color, A Point, B Point): draw a line between two points.
* DrawRect(Color, Rect): draw a rectangle outline between two points.
* DrawBox(Color, Rect): draw a filled rectangle between two points.
* DrawText(Text, Point): draw text at a location.
* StoreTexture(name string, image.Image): load a Go image.Image object into
  the engine as a "texture" that can be re-used and pasted on the canvas.
* LoadTexture(filename string): load an image from disk into a texture.
* Copy(Texturer, src Rect, dst Rect): copy a texture onto the canvas.

## Drawing Types

This package defines a handful of types useful for drawing operations.
See the godoc for full details.

* Color: an RGBA color holding uint8 values for each channel.
  * NewRGBA(red, green, blue, alpha uint8) to construct a new color.
* Point: holds an X,Y pair of coordinates.
* Rect: holds an X,Y and a W,H value.
* Text: holds text and configuration for rendering (color, stroke, shadow,
  size, etc.)

## Shape Generator Functions

The render package includes a few convenience functions for drawing
complex shapes.

The generator functions return a channel that yields all of the Points
that should be drawn to complete the shape. Example:

```go
var (
    A Point = render.NewPoint(10, 10)
    B Point = render.NewPoint(15, 20)
)

for pt := range render.IterLine(A, B) {
    engine.DrawPoint(render.Red, pt)
}
```

* IterLine(A Point, B Point): draw a line from A to B.
* IterRect(A Point, B Point): iterate all the points to draw a rectangle.
* IterEllipse(A Point, B Point): draw an elipse fitting inside the
  rectangle bounded by points A and B.

## Multitouch Gesture Notes

Support for SDL2's MultiGestureEvent is added on October 6 2021.
The event.State will have the property `Touching=true` while the engine
believes multitouch gestures are afoot. This begins when the user touches
the screen with two fingers, and _then_ motion is detected.

SDL2 spams us with gesture events for each tiny change detected, and then
just stops. The SDL driver in this repo doesn't set ev.State.Touching=false.
One heuristic you may use in your program to detect when multitouch has
ended is this:

SDL2 always emulates the mouse Button1 click for one of the fingers.
Record the position at the first Touching=true event, and monitor for
delta changes in position as the "mouse cursor" moves. When delta
movements become stale and don't update, you can set State.Touching=false
in your program.

## License

MIT.
