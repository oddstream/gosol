package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/util"
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
func NewToolbar(observer input.Observer) *Toolbar {
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
		NewRuneButton(rune(9776), tb.symbolFace, func() { observer.NotifyCallback(ebiten.KeyMenu) }, -1),
		NewLabel("", tb.symbolFace, 0),
		NewRuneButton('?', tb.symbolFace, func() { observer.NotifyCallback(ebiten.KeyH) }, 1),
		NewRuneButton(rune(8592), tb.symbolFace, func() { observer.NotifyCallback(ebiten.KeyU) }, 1),
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

// Rect returns the area this toolbar covers
func (tb *Toolbar) Rect() (x0, y0, x1, y1 int) {
	x0 = 0
	y0 = 0
	x1 = tb.width
	y1 = 48
	return // using named parameters
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	wgts := u.toolbar.widgets
	wgts[1] = NewLabel(title, u.toastTextFace, 0)
	u.toolbar.width = 0 // force img to be recreated
}

// Tapped is called when a tap happens over the toolbar
func (tb *Toolbar) Tapped(x, y int) {
	for _, w := range tb.widgets {
		if util.InRect(x, y, w.Rect) {
			println("UI widget tapped")
			w.Action()
		}
	}
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
