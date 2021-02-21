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

// findPileAt finds the pile under the mouse click or touch
func (b *Baize) findPileAt(pt image.Point) CardOwner {
	for _, o := range b.owners {
		if util.InRect(pt, o.FannedRect) {
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
	return nil
}

// PileTapped is called when a pile has been tapped
func (b *Baize) PileTapped(o CardOwner) {
	println("pile", o.Class(), "tapped")
}

// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {
	type HasTapTarget interface {
		TapTarget() string
	}

	typ, ok := c.owner.(HasTapTarget)
	if !ok {
		return
	}
	targetClass := typ.TapTarget()
	if targetClass == "" {
		return
	}
	for _, o := range b.owners {
		if targetClass == o.Class() {
			// println("found a", o.Class())
			if o.CanAcceptCard(c) {
				// println(o.Class(), "can accept", c.id)
				moveCards0(c.owner, o, 1)
			}
		}
	}
}

func moveCards0(src, dst CardOwner, nCards int) int {
	var nMoved int = 0

	if nCards == 1 && len(src.Cards()) > 0 {
		c := src.Pop()
		if src.Class() == "Stock" {
			c.FlipUp()
		}
		dst.Push(c)
		nMoved = 1
	} else if nCards > 1 {
		var tmp []*Card
		for n := nCards; n > 0 && len(src.Cards()) > 0; n-- {
			c := src.Pop()
			if src.Class() == "Stock" {
				c.FlipUp()
			}
			tmp = append(tmp, c)
		}
		for len(tmp) > 0 {
			c := tmp[len(tmp)-1]
			tmp = tmp[:len(tmp)-1]
			dst.Push(c)
			nMoved++
		}
	}
	return nMoved
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
			sx, sy := s.Position()
			// maybe user is tapping or starting to drag a card
			c := b.findCardAt(image.Point{X: sx, Y: sy})
			if c != nil {
				b.stroke = s
				b.stroke.SetDraggingObject(c)
				c.owner.StartDrag(c)
				// Pile.StartDrag(Card*)
				// TODO this card and the rest in the pile are being dragged
			} else {
				// maybe user is tapping an empty pile (eg to recycle waste to stock)
				o := b.findPileAt(image.Point{X: sx, Y: sy})
				if o != nil {
					b.stroke = s
					b.stroke.SetDraggingObject(o)
				}
			}
		}
	} else {
		b.stroke.Update()
		switch v := b.stroke.DraggingObject().(type) {
		case *Card:
			c := v
			if b.stroke.IsReleased() {
				if b.stroke.IsTapped() {
					c.owner.StopDrag(c)
					b.CardTapped(c)
				} else {
					sx, sy := b.stroke.Position()
					o := b.findPileAt(image.Point{X: sx, Y: sy})
					if o == nil {
						println("no pile found")
						c.owner.CancelDrag(c)
					}
					if o != nil {
						println("found pile", o.Class())
						if o.CanAcceptCard(c) {
							c.owner.StopDrag(c)
							moveCards0(c.owner, o, 1)
						} else {
							c.owner.CancelDrag(c)
						}
					}
				}
				b.stroke = nil
			} else {
				// would have used https://golang.org/ref/spec#Method_expressions
				// but couldn't figure out the syntax
				// so using a standalone loop instead
				// or could have (*Pile) DragTailBy(c, int, int) method
				dx, dy := b.stroke.PositionDiff()
				cards := c.owner.Cards()
				marking := false
				for i := 0; i < len(cards); i++ {
					ci := cards[i]
					if !marking && ci == c {
						marking = true
					}
					if marking {
						ci.DragBy(dx, dy)
					}
				}
			}
		case CardOwner:
			o := v
			if b.stroke.IsReleased() {
				if b.stroke.IsTapped() {
					b.PileTapped(o)
				}
				b.stroke = nil
			}
		default:
			log.Fatal("unknown type of dragging object")
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
	for _, o := range b.owners {
		o.DrawCards(screen)
	}
	for _, o := range b.owners {
		o.DrawMovingCards(screen)
	}
}
