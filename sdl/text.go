package sdl

import (
	"fmt"
	"sync"

	"git.kirsle.net/go/render"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TODO: font filenames
var DefaultFontFilename = "DejaVuSans.ttf"

// Font holds cached SDL_TTF structures for loaded fonts. They are created
// automatically when fonts are either preinstalled (InstallFont) or loaded for
// the first time as demanded by the DrawText method.
type Font struct {
	Filename string
	data     []byte // raw binary data of font
	ttf      *ttf.Font
}

var (
	fonts         = map[string]*ttf.Font{} // keys like "DejaVuSans@14" by font size
	installedFont = map[string][]byte{}    // installed font files' binary handles
	fontsMu       sync.RWMutex
)

// InstallFont preloads the font cache using TTF binary data in memory.
func InstallFont(filename string, binary []byte) {
	fontsMu.Lock()
	installedFont[filename] = binary
	fontsMu.Unlock()
}

// LoadFont loads and caches the font at a given size.
func LoadFont(filename string, size int) (*ttf.Font, error) {
	if filename == "" {
		filename = DefaultFontFilename
	}

	// Cached font available?
	keyName := fmt.Sprintf("%s@%d", filename, size)
	if font, ok := fonts[keyName]; ok {
		return font, nil
	}

	// Do we have this font in memory?
	var (
		font *ttf.Font
		err  error
	)

	fontsMu.Lock()
	defer fontsMu.Unlock()
	if binary, ok := installedFont[filename]; ok {
		var RWops *sdl.RWops
		RWops, err = sdl.RWFromMem(binary)
		if err != nil {
			return nil, fmt.Errorf("LoadFont(%s): RWFromMem: %s", filename, err)
		}

		font, err = ttf.OpenFontRW(RWops, 0, size)
	} else {
		font, err = ttf.OpenFont(filename, size)
	}

	// Error opening the font?
	if err != nil {
		return nil, fmt.Errorf("LoadFont(%s): %s", filename, err)
	}

	// Cache this font name and size.
	fonts[keyName] = font

	return font, nil
}

// ComputeTextRect computes and returns a Rect for how large the text would
// appear if rendered.
func (r *Renderer) ComputeTextRect(text render.Text) (render.Rect, error) {
	var (
		rect    render.Rect
		font    *ttf.Font
		surface *sdl.Surface
		color   = ColorToSDL(text.Color)
		err     error
	)

	if font, err = LoadFont(text.FontFilename, text.Size); err != nil {
		return rect, err
	}

	if surface, err = font.RenderUTF8Blended(text.Text, color); err != nil {
		return rect, err
	}
	defer surface.Free()

	rect.W = int(surface.W)
	rect.H = int(surface.H)
	return rect, err
}

// DrawText draws text on the canvas.
func (r *Renderer) DrawText(text render.Text, point render.Point) error {
	var (
		font    *ttf.Font
		surface *sdl.Surface
		tex     *sdl.Texture
		err     error
	)

	if font, err = LoadFont(text.FontFilename, text.Size); err != nil {
		return err
	}

	write := func(dx, dy int32, color sdl.Color) {
		if surface, err = font.RenderUTF8Blended(text.Text, color); err != nil {
			return
		}
		defer surface.Free()

		if tex, err = r.renderer.CreateTextureFromSurface(surface); err != nil {
			return
		}
		defer tex.Destroy()

		tmp := &sdl.Rect{
			X: int32(point.X) + dx,
			Y: int32(point.Y) + dy,
			W: surface.W,
			H: surface.H,
		}
		r.renderer.Copy(tex, nil, tmp)
	}

	// Does the text have a stroke around it?
	if text.Stroke != render.Invisible {
		color := ColorToSDL(text.Stroke)
		write(-1, -1, color)
		write(-1, 0, color)
		write(-1, 1, color)
		write(1, -1, color)
		write(1, 0, color)
		write(1, 1, color)
		write(0, -1, color)
		write(0, 1, color)
	}

	// Does it have a drop shadow?
	if text.Shadow != render.Invisible {
		write(1, 1, ColorToSDL(text.Shadow))
	}

	// Draw the text itself.
	write(0, 0, ColorToSDL(text.Color))

	return err
}

// shiftMap maps keys to their Shift versions.
var shiftMap = map[string]string{
	"`": "~",
	"1": "!",
	"2": "@",
	"3": "#",
	"4": "$",
	"5": "%",
	"6": "^",
	"7": "&",
	"8": "*",
	"9": "(",
	"0": ")",
	"-": "_",
	"=": "+",
	"[": "{",
	"]": "}",
	`\`: "|",
	";": ":",
	`'`: `"`,
	",": "<",
	".": ">",
	"/": "?",
}
