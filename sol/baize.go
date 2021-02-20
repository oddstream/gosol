package sol

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/util"
)

// Baize object describes the baize
type Baize struct {
	owners []CardOwner
	stroke *Stroke
}

// NewBaize is the factory func for Baize object
func NewBaize() *Baize {
	b := &Baize{}

	o2, ok := buildVariant("Klondike")
	if !ok {
		log.Fatal("Klondike" + " not found")
	}
	b.owners = o2

	b.dealCards()

	return b
}

func (b *Baize) dealCards() {
	stock := b.findPile("Stock")
	for _, o := range b.owners {
		deal := o.Deal()
		if deal == "" {
			continue
		}
		for _, d := range deal {
			switch d {
			case 'u':
				c := stock.Pop()
				c.FlipUp()
				o.Push(c)
			case 'd':
				c := stock.Pop()
				c.FlipDown()
				o.Push(c)
			}
		}
	}
}

func (b *Baize) findPile(cls string) CardOwner {
	for _, o := range b.owners {
		if o.Class() == cls {
			return o
		}
	}
	return nil
}

// findTileAt finds the tile under the mouse click or touch
func (b *Baize) findCardAt(pt image.Point) *Card {
	for _, o := range b.owners {
		cards := o.Cards()
		for i := len(cards) - 1; i >= 0; i-- {
			c := cards[i]
			if util.InRect(pt, c.Rect) {
				return c
			}
		}
	}
	// for _, c := range b.stock.cards {
	// 	if util.InRect(pt, c.Rect) {
	// 		return c
	// 	}
	// }
	return nil
}

// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {
	c.Flip()
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	for _, o := range b.owners {
		o.Layout(outsideWidth, outsideHeight)
	}

	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	if b.stroke == nil {
		var s *Stroke

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s = NewStroke(&MouseStrokeSource{})
		}
		ts := inpututil.JustPressedTouchIDs()
		if ts != nil && len(ts) == 1 {
			s = NewStroke(&TouchStrokeSource{ts[0]})
		}

		if s != nil {
			cx, cy := s.Position()
			c := b.findCardAt(image.Point{X: cx, Y: cy})
			if c != nil {
				b.stroke = s
				b.stroke.SetDraggingObject(c)
				c.StartDrag()
			}
		}
	} else {
		b.stroke.Update()
		c := b.stroke.DraggingObject().(*Card)
		if b.stroke.IsReleased() {
			if b.stroke.IsTapped() {
				c.StopDrag()
				b.CardTapped(c)
			} else {
				c.CancelDrag()
			}
			b.stroke = nil
		} else {
			x, y := c.DragStartPosition()
			dx, dy := b.stroke.PositionDiff()
			c.SetPosition(x+dx, y+dy)
		}
	}

	for _, o := range b.owners {
		o.Update()
	}

	CTQ.Update()

	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	for _, o := range b.owners {
		o.Draw(screen)
	}
	// b.stock.Draw(screen)
	// b.waste.Draw(screen)
	// for _, f := range b.foundations {
	// 	f.Draw(screen)
	// }

	for _, o := range b.owners {
		o.DrawCards(screen)
	}
	// b.stock.DrawCards(screen)
	// b.waste.DrawCards(screen)
	// for _, f := range b.foundations {
	// 	f.DrawCards(screen)
	// }
}
