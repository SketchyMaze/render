# Hello WASM

This is a port of the "hello-world" example for a WebAssembly target.

Differences compared to the SDL2 version include:

* Import render/canvas instead of render/sdl
* Load the gopher.png via ajax HTTP request, as the syscall to open a file from
  disk is not supported in a WASM environment.
* A call to the `canvas.AddEventListeners()` function to set up DOM event
  bindings, which enables key and mouse events to be handled in the
  canvas.Engine.

```diff
--- ../hello-world/main.go	2019-12-22 15:17:42.186838967 -0800
+++ main_wasm.go	2019-12-22 16:21:58.449699021 -0800
@@ -1,13 +1,17 @@
+// +build js,wasm
+
+// This is a version of the 'hello-world' example but build for WebAssembly.
+
 package main

 import (
 	"image/png"
 	"log"
-	"os"
+	"net/http"
 	"time"

 	"git.kirsle.net/go/render"
-	"git.kirsle.net/go/render/sdl"
+	"git.kirsle.net/go/render/canvas"
 )

 var (
@@ -30,7 +34,11 @@
 )

 func main() {
-	mw := sdl.New("Hello World", 800, 600)
+	mw, _ := canvas.New("canvas")
 	setup(mw)

 	for {
@@ -61,17 +69,22 @@
 	}
 }

-func setup(e render.Engine) {
+func setup(e *canvas.Engine) {
 	if err := e.Setup(); err != nil {
 		panic(err)
 	}

-	// Load gopher sprite.
-	fh, err := os.Open("gopher.png")
+	// Bind DOM event handlers.
+	e.AddEventListeners()
+
+	// Load gopher sprite via ajax request.
+	resp, err := http.Get("gopher.png")
 	if err != nil {
 		log.Fatalf("read gopher.png: %s", err)
 	}
+	defer resp.Body.Close()

-	img, _ := png.Decode(fh)
+	img, _ := png.Decode(resp.Body)

 	gopher, _ = e.StoreTexture("gopher.png", img)
 	texSize = gopher.Size()
```

## Credits

Gopher image was created by Takuya Ueda (https://twitter.com/tenntenn)
from https://github.com/golang-samples/gopher-vector

This example comes with the DejaVu Sans font. License information
at https://dejavu-fonts.github.io/License.html
