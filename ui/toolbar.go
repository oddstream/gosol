package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

//go:embed assets/DejaVuSans-Bold.ttf
var symbolFontBytes []byte

// Toolbar object (hamburger button, variant name, undo button)
type Toolbar struct {
	img        *ebiten.Image
	symbolFace font.Face
	title      string
	width      int
	widgets    []Widget
}

// NewToolbar creates a new toolbar
func NewToolbar() *Toolbar {
	tb := &Toolbar{}

	tt, err := truetype.Parse(symbolFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	tb.symbolFace = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	symbolFontBytes = nil

	tb.widgets = []Widget{
		NewRuneButton(rune(9776), tb.symbolFace, func() {}, -1),
		NewLabel("", tb.symbolFace, 0),
		NewRuneButton('?', tb.symbolFace, func() {}, 1),
		NewRuneButton(rune(8592), tb.symbolFace, func() {}, 1),
	}

	return tb
}

func (tb *Toolbar) createImg() {

	dc := gg.NewContext(tb.width, 48)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(tb.width), 48)
	dc.Fill()
	dc.Stroke()

	width := dc.Width()
	height := dc.Height()
	nextLeft := 48 / 2
	nextRight := width - 48/2
	for _, w := range tb.widgets {
		switch w.Align() {
		case -1:
			w.Draw(dc, nextLeft, height/2)
			nextLeft += 48
		case 0:
			w.Draw(dc, width/2, height/2)
		case 1:
			w.Draw(dc, nextRight, height/2)
			nextRight -= 48
		}
	}
	tb.img = ebiten.NewImageFromImage(dc.Image())
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	wgts := u.toolbar.widgets
	wgts[1] = NewLabel(title, u.toastTextFace, 0)
	u.toolbar.width = 0 // force img to be recreated
}

// Update the toolbar
func (tb *Toolbar) Update() {
}

// Draw the toolbar
func (tb *Toolbar) Draw(screen *ebiten.Image) {
	w, _ := screen.Size()
	if tb.img == nil || w != tb.width {
		tb.width = w
		tb.createImg()
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tb.img, op)
}
