package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	BarBase
	title string
}

func (tb *Toolbar) createImg() *ebiten.Image {
	// override BarBase.createImg to draw title
	dc := gg.NewContext(tb.width, tb.height) // should always be 48,48
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(tb.width), float64(tb.height))
	dc.Fill()
	if tb.title != "" {
		dc.SetFontFace(schriftbank.RobotoMedium24)
		dc.SetColor(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
		dc.DrawStringAnchored(tb.title, float64(tb.width)/2, float64(tb.height)/2, 0.5, 0.5)
	}
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	// img will created first time it's drawn if width == 0
	tb := &Toolbar{BarBase: BarBase{input: input, x: 0, y: 0, width: 0, height: 48}}

	tb.widgets = []Widget{
		// button's x will be set by LayoutWidgets() (y will always be 0 in a toolbar)
		NewIconButton(tb, input, 0, 0, 48, 48, -1, "menu", ebiten.KeyMenu),
		NewIconButton(tb, input, 0, 0, 48, 48, 1, "info", ebiten.KeyI),
		NewIconButton(tb, input, 0, 0, 48, 48, 1, "done", ebiten.KeyC),
		NewIconButton(tb, input, 0, 0, 48, 48, 1, "undo", ebiten.KeyU),
	}
	return tb
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	// u.toolbar.ReplaceWidget(1, NewLabel(u.toolbar, u.input, 0, 0, 0, 48, 0, title, schriftbank.RobotoMedium24, ""))
	u.toolbar.title = title
	u.toolbar.width = 0 // force img to be recreated
}

// Draw the toolbar; override to use our own createImg
func (tb *Toolbar) Draw(screen *ebiten.Image) {
	w, _ := screen.Size()
	if tb.img == nil || w != tb.width {
		tb.width = w
		tb.img = tb.createImg()
		tb.LayoutWidgets()
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(tb.img, op) // draw toolbar at 0,0

	for _, w := range tb.widgets {
		w.Draw(screen)
	}
}
