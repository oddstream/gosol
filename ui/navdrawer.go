package ui

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

const (
	aniLeft  = -1
	aniRight = 1
	aniHide  = -2
	aniShow  = 2
	aniStop  = 0
)

// NavDrawer slide out modal menu
type NavDrawer struct {
	img           *ebiten.Image
	width, height int
	x, y          int
	aniState      int
	widgets       []Widget
}

func (n *NavDrawer) createImg() {

	dc := gg.NewContext(n.width, n.height)
	dc.SetColor(color.RGBA{R: 0x32, G: 0x32, B: 0x32, A: 0xff})
	dc.DrawRectangle(0, 0, float64(n.width), float64(n.height))
	dc.Fill()
	dc.Stroke()
	n.img = ebiten.NewImageFromImage(dc.Image())
}

// NewNavDrawer creates the NavDrawer object; it starts life off screen to the left
func NewNavDrawer() *NavDrawer {
	// according to https://material.io/components/navigation-drawer#specs, always 256 wide
	n := &NavDrawer{width: 256, x: -256, y: 0}
	return n
}

// Rect returns the area the NavDrawer currently covers (it may be off screen to the left)
func (n *NavDrawer) Rect() (x0, y0, x1, y1 int) {
	x0 = n.x
	y0 = n.y
	x1 = n.x + n.width
	y1 = n.y + n.height
	return // using named parameters
}

// Show starts to animate the drawer on screen from the left
func (n *NavDrawer) Show() {
	n.aniState = aniRight
}

// Hide starts to animate the drawer off screen to the left
func (n *NavDrawer) Hide() {
	n.aniState = aniLeft
}

// Visible returns true if the NavDrawer is showing
func (n *NavDrawer) Visible() bool {
	return n.x == 0
}

// Tapped is called when a tap happens over the toolbar
func (n *NavDrawer) Tapped(x, y int) {
	for _, w := range n.widgets {
		if util.InRect(x, y, w.Rect) {
			println("UI toolbar tapped")
			w.Action()
		}
	}
}

// Update the NavDrawer
func (n *NavDrawer) Update() {
	switch n.aniState {
	case aniLeft:
		if n.x <= -256 {
			n.x = -256
			n.aniState = aniStop
		} else {
			n.x -= 16
		}
	case aniRight:
		if n.x >= 0 {
			n.x = 0
			n.aniState = aniStop
		} else {
			n.x += 16
		}
	}
}

// Draw the NavDrawer
func (n *NavDrawer) Draw(screen *ebiten.Image) {
	h, _ := screen.Size()
	if n.img == nil || h != n.height {
		n.height = h
		n.createImg()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.x), 0) // y is always zero
	screen.DrawImage(n.img, op)
}
