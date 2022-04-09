package sdl

import (
	"bytes"
	"fmt"
	"image"

	"git.kirsle.net/go/render"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/image/bmp"
)

// Copy a texture into the renderer.
func (r *Renderer) Copy(t render.Texturer, src, dst render.Rect) {
	if tex, ok := t.(*Texture); ok {
		var (
			a = RectToSDL(src)
			b = RectToSDL(dst)
		)
		r.renderer.Copy(tex.tex, &a, &b)
	}
}

// Texture can hold on to SDL textures for caching and optimization.
type Texture struct {
	render *Renderer // backref to free them up thoroughly
	tex    *sdl.Texture
	image  image.Image
	width  int32
	height int32
}

// StoreTexture caches an SDL texture from a bitmap.
func (r *Renderer) StoreTexture(name string, img image.Image) (render.Texturer, error) {
	var (
		fh = bytes.NewBuffer([]byte{})
	)

	err := bmp.Encode(fh, img)
	if err != nil {
		return nil, fmt.Errorf("NewTexture: bmp.Encode: %s", err)
	}

	// Create an SDL RWOps from the bitmap data in memory.
	sdlRW, err := sdl.RWFromMem(fh.Bytes())
	if err != nil {
		return nil, fmt.Errorf("NewTexture: sdl.RWFromMem: %s", err)
	}

	surface, err := sdl.LoadBMPRW(sdlRW, true)
	if err != nil {
		return nil, fmt.Errorf("NewTexture: sdl.LoadBMPRW: %s", err)
	}
	defer surface.Free()

	// TODO: chroma key color hardcoded to white here
	key := sdl.MapRGB(surface.Format, 255, 255, 255)
	surface.SetColorKey(true, key)

	texture, err := r.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("NewBitmap: create texture: %s", err)
	}

	tex := &Texture{
		render: r,
		width:  surface.W,
		height: surface.H,
		tex:    texture,
		image:  img,
	}
	r.textures[name] = tex

	return tex, nil
}

// CountTextures is a custom function for the SDL2 Engine only that returns the
// size of the engine texture cache.
//
// Returns the number of (Go) Texture objects (which hold cached image.Image
// and can garbage collect nicely) and the number of those which also have
// SDL2 Texture objects (which need to be freed manually).
func (r *Renderer) CountTextures() (bitmaps int, sdl2textures int) {
	var withTex int // ones with active SDL2 Texture objects
	for _, tex := range r.textures {
		if tex.tex != nil {
			withTex++
		}
	}

	return len(r.textures), withTex
}

// ListTextures is a custom function to peek into the SDL2 texture cache names.
func (r *Renderer) ListTextures() []string {
	var keys = []string{}
	for key := range r.textures {
		keys = append(keys, key)
	}
	return keys
}

// Size returns the dimensions of the texture.
func (t *Texture) Size() render.Rect {
	return render.NewRect(int(t.width), int(t.height))
}

// Image returns the underlying Go image.Image.
func (t *Texture) Image() image.Image {
	return t.image
}

// Free the SDL2 texture object.
func (t *Texture) Free() error {
	var err error

	if t.tex != nil {
		err = t.tex.Destroy()
		t.tex = nil
	}

	// Free up the cached texture too to garbage collect the image.Image cache etc.
	for name, tex := range t.render.textures {
		if tex == t {
			delete(t.render.textures, name)
			break
		}
	}

	return err
}

// LoadTexture initializes a texture from a bitmap image.
func (r *Renderer) LoadTexture(name string) (render.Texturer, error) {
	if tex, ok := r.textures[name]; ok {
		// If the SDL2 texture had been freed, recreate it.
		if tex.tex == nil {
			return r.StoreTexture(name, tex.image)
		}
		return tex, nil
	}
	return nil, fmt.Errorf("LoadTexture(%s): not found in texture cache", name)
}

// FreeTextures flushes the internal cache of SDL2 textures and frees their memory.
func (r *Renderer) FreeTextures() int {
	var num = len(r.textures)
	for name, tex := range r.textures {
		delete(r.textures, name)
		tex.tex.Destroy()
	}
	return num
}
