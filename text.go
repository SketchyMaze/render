package render

import "fmt"

// Text holds information for drawing text.
type Text struct {
	Text         string
	Size         int
	Color        Color
	Padding      int
	PadX         int
	PadY         int
	Stroke       Color  // Stroke color (if not zero)
	Shadow       Color  // Drop shadow color (if not zero)
	FontFilename string // Path to *.ttf file on disk
}

func (t Text) String() string {
	return fmt.Sprintf(`Text<"%s" %dpx %s>`, t.Text, t.Size, t.Color)
}

// IsZero returns if the Text is the zero value.
func (t Text) IsZero() bool {
	return t.Text == "" && t.Size == 0 && t.Color == Invisible && t.Padding == 0 && t.PadX == 0 && t.PadY == 0 && t.Stroke == Invisible && t.Shadow == Invisible
}
