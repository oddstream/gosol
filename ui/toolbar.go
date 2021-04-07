package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// Toolbar object (hamburger button, variant name, undo, help buttons)
type Toolbar struct {
	BarBase
}

// func (tb *Toolbar) createImg() *ebiten.Image {
// 	// override BarBase.createImg to draw title
// 	dc := gg.NewContext(tb.width, tb.height) // should always be 48,48
// 	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
// 	dc.DrawRectangle(0, 0, float64(tb.width), float64(tb.height))
// 	dc.Fill()
// 	if tb.title == "" {
// 		tb.title = "(unnamed)"
// 	}
// 	if tb.title != "" {
// 		dc.SetFontFace(schriftbank.RobotoMedium24)
// 		dc.SetColor(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
// 		dc.DrawStringAnchored(tb.title, float64(tb.width)/2, float64(tb.height)/2, 0.5, 0.5)
// 	}
// 	dc.Stroke()
// 	return ebiten.NewImageFromImage(dc.Image())
// }

// NewToolbar creates a new toolbar
func NewToolbar(input *input.Input) *Toolbar {
	// img will created first time it's drawn if width == 0
	tb := &Toolbar{BarBase: BarBase{input: input, x: 0, y: 0, width: 0, height: 48}}

	tb.widgets = []Widget{
		// button's x will be set by LayoutWidgets() (y will always be 0 in a toolbar)
		NewIconButton(tb, input, 0, 0, 48, 48, -1, "menu", ebiten.KeyMenu),
		NewLabel(tb, input, 0, "title", schriftbank.RobotoMedium24, ""),
		NewIconButton(tb, input, 0, 0, 48, 48, 1, "undo", ebiten.KeyU),
		NewIconButton(tb, input, 0, 0, 48, 48, 1, "done", ebiten.KeyC),
	}
	return tb
}

// SetTitle of the toolbar
func (u *UI) SetTitle(title string) {
	var l *Label = u.toolbar.widgets[1].(*Label)
	l.UpdateText(title)
	// u.toolbar.LayoutWidgets()
}

// Layout implements Ebiten's Layout
func (tb *Toolbar) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get screen height and position bar
	if tb.img == nil || outsideWidth != tb.width {
		tb.width = outsideWidth
		// tb.height is fixed (at 48)
		tb.img = tb.createImg()
		tb.LayoutWidgets()
	}
	// tb.x, tb.y = 0, 0
	return outsideWidth, outsideHeight
}

// Draw the toolbar; override to use our own createImg
// func (tb *Toolbar) Draw(screen *ebiten.Image) {
// 	w, _ := screen.Size()
// 	if tb.img == nil || w != tb.width {
// 		tb.width = w
// 		tb.img = tb.createImg()
// 		tb.LayoutWidgets()
// 	}
// 	op := &ebiten.DrawImageOptions{}
// 	screen.DrawImage(tb.img, op) // draw toolbar at 0,0

// 	for _, w := range tb.widgets {
// 		w.Draw(screen)
// 	}
// }
