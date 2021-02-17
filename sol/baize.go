package sol

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/util"
)

// Baize object describes the baize
type Baize struct {
	stock       *Stock
	waste       *Waste
	foundations []*Foundation
	stroke      *Stroke
}

// NewBaize is the factory func for Baize object
func NewBaize() *Baize {
	b := &Baize{}
	b.stock = NewStock(1, 1, 1)
	b.stock.Shuffle()

	b.waste = NewWaste(2, 1)

	for i := 0; i < 4; i++ {
		b.foundations = append(b.foundations, NewFoundation(4+i, 1))
	}
	return b
}

// findTileAt finds the tile under the mouse click or touch
func (b *Baize) findCardAt(pt image.Point) *Card {
	for i := len(b.stock.cards) - 1; i >= 0; i-- {
		c := b.stock.cards[i]
		if util.InRect(pt, c.Rect) {
			return c
		}
	}
	// for _, c := range b.stock.cards {
	// 	if util.InRect(pt, c.Rect) {
	// 		return c
	// 	}
	// }
	return nil
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	var s *Stroke

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s = NewStroke(&MouseStrokeSource{})
	}
	ts := inpututil.JustPressedTouchIDs()
	if ts != nil && len(ts) == 1 {
		s = NewStroke(&TouchStrokeSource{ts[0]})
	}

	if s != nil {
		c := b.findCardAt(s.Position())
		if c != nil {
			b.stroke = s
			b.stroke.SetDraggingObject(c)
			// TODO move Card to front?
		}
	}

	if b.stroke != nil {
		b.stroke.Update()
		c := b.stroke.DraggingObject().(*Card)
		pt := b.stroke.PositionDiff()
		x, y := c.owner.Position()
		c.PositionTo(x+pt.X, y+pt.Y)

		if b.stroke.IsReleased() {
			c := b.stroke.DraggingObject().(*Card)
			c.TransitionBackToPile()

			if b.stroke.IsTapped() {
				println("tap detected on", c.id)
			}
			b.stroke = nil
		}
	}

	b.stock.Update()
	b.waste.Update()
	for _, f := range b.foundations {
		f.Update()
	}
	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	b.stock.Draw(screen)
	b.waste.Draw(screen)
	for _, f := range b.foundations {
		f.Draw(screen)
	}
}
