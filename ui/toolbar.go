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
}

func (tb *Toolbar) createImg() {

	if tb.symbolFace == nil {
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
	}

	dc := gg.NewContext(tb.width, 48)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(tb.width), 48)
	dc.Fill()
	dc.Stroke()

	dc.SetFontFace(tb.symbolFace)
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(string(rune(9776)), 24, 24, 0.5, 0.5) // trigram for heaven (U+2630)
	dc.DrawStringAnchored(string(rune(8592)), float64(tb.width-24), 24, 0.5, 0.5)
	dc.Stroke()

	if tb.title != "" {
		// TODO use a Roboto Bold font
		// dc.SetFontFace(u.toastTextFace)
		dc.SetRGBA(1, 1, 1, 1)
		dc.DrawStringAnchored(tb.title, float64(tb.width/2), 48/2, 0.5, 0.5)
		dc.Stroke()
	}

	tb.img = ebiten.NewImageFromImage(dc.Image())
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	u.toolbar.title = title
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
