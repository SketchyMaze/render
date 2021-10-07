package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
)

var (
	// Regexps to parse hex color codes. Three formats are supported:
	// * reHexColor3 uses only 3 hex characters, like #F90
	// * reHexColor6 uses standard 6 characters, like #FF9900
	// * reHexColor8 is the standard 6 plus alpha channel, like #FF9900FF
	reHexColor3 = regexp.MustCompile(`^([A-Fa-f0-9])([A-Fa-f0-9])([A-Fa-f0-9])$`)
	reHexColor6 = regexp.MustCompile(`^([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})$`)
	reHexColor8 = regexp.MustCompile(`^([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})$`)
)

// Common color names.
var (
	Invisible  = Color{}
	White      = RGBA(255, 255, 255, 255)
	Grey       = RGBA(153, 153, 153, 255)
	DarkGrey   = RGBA(64, 64, 64, 255)
	Black      = RGBA(0, 0, 0, 255)
	SkyBlue    = RGBA(0, 153, 255, 255)
	Blue       = RGBA(0, 0, 255, 255)
	DarkBlue   = RGBA(0, 0, 153, 255)
	Red        = RGBA(255, 0, 0, 255)
	DarkRed    = RGBA(153, 0, 0, 255)
	Green      = RGBA(0, 255, 0, 255)
	DarkGreen  = RGBA(0, 153, 0, 255)
	Cyan       = RGBA(0, 255, 255, 255)
	DarkCyan   = RGBA(0, 153, 153, 255)
	Yellow     = RGBA(255, 255, 0, 255)
	Orange     = RGBA(255, 153, 0, 255)
	DarkYellow = RGBA(153, 153, 0, 255)
	Magenta    = RGBA(255, 0, 255, 255)
	Purple     = RGBA(153, 0, 153, 255)
	Pink       = RGBA(255, 153, 255, 255)
)

// Color holds an RGBA color value.
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// RGBA creates a new Color.
func RGBA(r, g, b, a uint8) Color {
	return Color{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a,
	}
}

// FromColor creates a render.Color from a Go color.Color
func FromColor(from color.Color) Color {
	// downscale a 16-bit color value to 8-bit. input range 0x0000..0xffff
	downscale := func(in uint32) uint8 {
		var scale = float64(in) / 0xffff
		return uint8(scale * 0xff)
	}
	r, g, b, a := from.RGBA()
	return RGBA(
		downscale(r),
		downscale(g),
		downscale(b),
		downscale(a),
	)
}

// ToRGBA converts to a standard Go color.Color
func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{
		R: c.Red,
		G: c.Green,
		B: c.Blue,
		A: c.Alpha,
	}
}

// MustHexColor parses a color from hex code or panics.
func MustHexColor(hex string) Color {
	color, err := HexColor(hex)
	if err != nil {
		panic(err)
	}
	return color
}

// HexColor parses a color from hexadecimal code.
func HexColor(hex string) (Color, error) {
	c := Black // default color

	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	var m []string
	if len(hex) == 3 {
		m = reHexColor3.FindStringSubmatch(hex)
		// Double up the hex characters.
		m[1] += m[1]
		m[2] += m[2]
		m[3] += m[3]
	} else if len(hex) == 6 {
		m = reHexColor6.FindStringSubmatch(hex)
	} else if len(hex) == 8 {
		m = reHexColor8.FindStringSubmatch(hex)
	} else {
		return c, errors.New("not a valid length for color code; only 3, 6 and 8 supported")
	}

	// Any luck?
	if m == nil {
		return c, errors.New("not a valid hex color code")
	}

	// Parse the color values. 16=base, 8=bit size
	red, _ := strconv.ParseUint(m[1], 16, 8)
	green, _ := strconv.ParseUint(m[2], 16, 8)
	blue, _ := strconv.ParseUint(m[3], 16, 8)

	// Alpha channel available?
	var alpha uint64 = 255
	if len(m) == 5 {
		alpha, _ = strconv.ParseUint(m[4], 16, 8)
	}

	c.Red = uint8(red)
	c.Green = uint8(green)
	c.Blue = uint8(blue)
	c.Alpha = uint8(alpha)
	return c, nil
}

func (c Color) String() string {
	return fmt.Sprintf(
		"Color<#%02x%02x%02x+%02x>",
		c.Red, c.Green, c.Blue, c.Alpha,
	)
}

// ToHex converts a render.Color to standard #RRGGBB hexadecimal format.
func (c Color) ToHex() string {
	return fmt.Sprintf(
		"#%02x%02x%02x",
		c.Red, c.Green, c.Blue,
	)
}

// ToColor converts a render.Color into a Go standard color.Color
func (c Color) ToColor() color.RGBA {
	return color.RGBA{
		R: c.Red,
		G: c.Green,
		B: c.Blue,
		A: c.Alpha,
	}
}

// Transparent returns whether the alpha channel is zeroed out and the pixel
// won't appear as anything when rendered.
func (c Color) Transparent() bool {
	return c.Alpha == 0x00
}

// MarshalJSON serializes the Color for JSON.
func (c Color) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`"#%02x%02x%02x"`,
		c.Red, c.Green, c.Blue,
	)), nil
}

// UnmarshalJSON reloads the Color from JSON.
func (c *Color) UnmarshalJSON(b []byte) error {
	var hex string
	err := json.Unmarshal(b, &hex)
	if err != nil {
		return err
	}

	parsed, err := HexColor(hex)
	if err != nil {
		return err
	}

	c.Red = parsed.Red
	c.Blue = parsed.Blue
	c.Green = parsed.Green
	c.Alpha = parsed.Alpha
	return nil
}

// IsZero returns if the color is all zeroes (invisible).
func (c Color) IsZero() bool {
	return c.Red+c.Green+c.Blue+c.Alpha == 0
}

// Add a relative color value to the color.
func (c Color) Add(r, g, b, a int) Color {
	var (
		R = int(c.Red) + r
		G = int(c.Green) + g
		B = int(c.Blue) + b
		A = int(c.Alpha) + a
	)

	cap8 := func(v int) uint8 {
		if v > 255 {
			v = 255
		} else if v < 0 {
			v = 0
		}
		return uint8(v)
	}

	return Color{
		Red:   cap8(R),
		Green: cap8(G),
		Blue:  cap8(B),
		Alpha: cap8(A),
	}
}

// AddColor adds another Color to your Color.
func (c Color) AddColor(other Color) Color {
	return c.Add(
		int(other.Red),
		int(other.Green),
		int(other.Blue),
		int(other.Alpha),
	)
}

// Lighten a color value.
func (c Color) Lighten(v int) Color {
	return c.Add(v, v, v, 0)
}

// Darken a color value.
func (c Color) Darken(v int) Color {
	return c.Add(-v, -v, -v, 0)
}

// Transparentize adjusts the alpha value.
func (c Color) Transparentize(v int) Color {
	return c.Add(0, 0, 0, v)
}

// SetAlpha sets the alpha value to a specific setting.
func (c Color) SetAlpha(v uint8) Color {
	c.Alpha = v
	return c
}
