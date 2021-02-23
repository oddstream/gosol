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
	Piles     []*Pile
	stroke    *Stroke
	UndoStack []SaveableBaize
}

// NewBaize is the factory func for Baize object
func NewBaize() *Baize {
	b := &Baize{}

	piles, ok := buildVariantPiles(TheUserData.Variant)
	if !ok {
		log.Fatal(TheUserData.Variant + " not found")
	}
	b.Piles = piles

	{
		maxX := 0
		for _, p := range b.Piles {
			if p.X > maxX {
				maxX = p.X
			}
		}
		ebiten.SetWindowSize((maxX+2)*(71+10), WindowHeight)
	}

	stock := b.findPile("Stock")
	createCards(stock)
	shuffleCards(stock)

	b.dealCards()

	return b
}

func (b *Baize) dealCards() {
	stock := b.findPile("Stock")
	for _, p := range b.Piles {
		deal := p.GetStringAttribute("Deal")
		if deal == "" {
			continue
		}
		for _, d := range deal {
			switch d {
			case 'u':
				c := stock.Pop()
				c.FlipUp()
				p.Push(c)
			case 'd':
				c := stock.Pop()
				c.FlipDown()
				p.Push(c)
			}
		}
	}
}

func (b *Baize) findPile(cls string) *Pile {
	for _, p := range b.Piles {
		if p.Class == cls {
			return p
		}
	}
	return nil
}

// findPileAt finds the pile under the mouse click or touch
func (b *Baize) findPileAt(pt image.Point) *Pile {
	for _, o := range b.Piles {
		if util.InRect(pt, o.FannedRect) {
			return o
		}
	}
	return nil
}

// findTileAt finds the tile under the mouse click or touch
func (b *Baize) findCardAt(pt image.Point) *Card {
	for _, p := range b.Piles {
		for i := len(p.Cards) - 1; i >= 0; i-- {
			c := p.Cards[i]
			if util.InRect(pt, c.Rect) {
				return c
			}
		}
	}
	return nil
}

// PileTapped is called when a pile has been tapped
func (b *Baize) PileTapped(p *Pile) {

	// this method is in Baize because it needs access to Baize.findPile()
	if p.Class != "Stock" {
		return
	}

	if p.localRecycles > 0 {
		waste := b.findPile("Waste")
		if waste == nil || len(waste.Cards) == 0 {
			return
		}
		stock := b.findPile("Stock")
		for len(waste.Cards) > 0 {
			c := waste.Pop()
			c.FlipDown()
			stock.Push(c)
		}
		p.localRecycles--
	}
	// println("pile", p.Class, "tapped")
}

// CardTapped is called when a card has been tapped
func (b *Baize) CardTapped(c *Card) {

	// println("card",c.id,"tapped")

	// can only tap top card
	if c != c.owner.Peek() {
		c.Shake()
		return
	}

	moved := false

	// Tap on a Stock card to send it to Waste
	// TODO fudge to send three cards
	targetClass := c.owner.GetStringAttribute("Target")
	if targetClass != "" {
		for _, p := range b.Piles {
			if targetClass == p.Class {
				// println("found a", p.Class)
				if p.CanAcceptCard(c) {
					// println(p.Class, "can accept", c.id)
					b.MoveCards(c, p)
					moved = true
				}
			}
		}
	}

	if !moved {
		c.Shake()
	}

	// TODO else test other piles to see if this card is accepted?
}

// func moveCards0(src, dst CardOwner, nCards int) int {
// 	var nMoved int = 0

// 	if nCards == 1 && len(src.Cards()) > 0 {
// 		c := src.Pop()
// 		if src.Class() == "Stock" {
// 			c.FlipUp()
// 		}
// 		dst.Push(c)
// 		nMoved = 1
// 	} else if nCards > 1 {
// 		var tmp []*Card
// 		for n := nCards; n > 0 && len(src.Cards()) > 0; n-- {
// 			c := src.Pop()
// 			if src.Class() == "Stock" {
// 				c.FlipUp()
// 			}
// 			tmp = append(tmp, c)
// 		}
// 		for len(tmp) > 0 {
// 			c := tmp[len(tmp)-1]
// 			tmp = tmp[:len(tmp)-1]
// 			dst.Push(c)
// 			nMoved++
// 		}
// 	}

// 	{
// 		cards := src.Cards()
// 		if len(cards) > 0 {
// 			cards[len(cards)-1].FlipUp()
// 		}
// 	}

// 	return nMoved
// }

// MoveCards from one pile to another, always from card downwards (inclusive)
func (b *Baize) MoveCards(c *Card, dst *Pile) {

	src := c.owner
	moveFrom := len(src.Cards)
	tmp := make([]*Card, 0)

	// find the index of the first card we will move
	for i, sc := range src.Cards {
		if sc == c {
			moveFrom = i
			break
		}
	}

	if moveFrom == len(src.Cards) {
		log.Fatal("MoveCards could not find card in source")
	}

	b.UndoPush()

	oldSrcLen := len(src.Cards)

	// pop the tail off the source and push onto temp stack
	for i := len(src.Cards) - 1; i >= moveFrom; i-- {
		sc := src.Pop()
		if src.Class == "Stock" {
			sc.FlipUp()
		}
		tmp = append(tmp, sc)
	}

	// pop cards off the temp stack and onto the destination
	for len(tmp) > 0 {
		dc := tmp[len(tmp)-1]
		tmp = tmp[:len(tmp)-1]
		dst.Push(dc)
	}

	if oldSrcLen == len(src.Cards) {
		b.UndoPop() // discard return values
	} else {
		// flip up an exposed source card
		if src.Class != "Stock" {
			tc := src.Peek()
			if tc != nil {
				tc.FlipUp()
			}
		}
	}
}

// UndoPush pushes the current state onto the undo stack
func (b *Baize) UndoPush() {
	b.UndoStack = append(b.UndoStack, b.Saveable())
}

// UndoPop pops a state off the undo stack
func (b *Baize) UndoPop() (SaveableBaize, bool) {
	if len(b.UndoStack) > 0 {
		sav := b.UndoStack[len(b.UndoStack)-1]
		b.UndoStack = b.UndoStack[:len(b.UndoStack)-1]
		return sav, true
	}
	return SaveableBaize{}, false
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.UndoStack) == 0 {
		println("Nothing to undo")
		return
	}
	sav, ok := b.UndoPop()
	if ok {
		b.UpdateFromSaveable(sav)
	}
}

func (b *Baize) largestIntersection(c *Card) *Pile {
	var largest int = 0
	var pile *Pile = nil
	cx0, cy0, cx1, cy1 := c.Rect()
	for _, p := range b.Piles {
		px0, py0, px1, py1 := p.FannedRect()
		i := util.OverlapArea(cx0, cy0, cx1, cy1, px0, py0, px1, py1)
		if i > largest {
			largest = i
			pile = p
		}
	}
	return pile
}

// Layout implements ebiten.Game's Layout.
func (b *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {

	for _, p := range b.Piles {
		p.Layout(outsideWidth, outsideHeight)
	}

	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (b *Baize) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyU) {
		b.Undo()
		return nil
	}

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
				if c.owner.StartDrag(c) {
					b.stroke = s
					b.stroke.SetDraggingObject(c)
				} else {
					println("Cannot drag those cards")
				}
			} else {
				// maybe user is tapping an empty pile (eg to recycle waste to stock)
				p := b.findPileAt(image.Point{X: sx, Y: sy})
				if p != nil {
					b.stroke = s
					b.stroke.SetDraggingObject(p)
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
					// sx, sy := b.stroke.Position()
					// p := b.findPileAt(image.Point{X: sx, Y: sy})
					p := b.largestIntersection(c)
					if p == nil {
						c.owner.CancelDrag(c)
					} else {
						// println("found pile", o.Class())
						if p.CanAcceptTail(c.owner.Tail) {
							c.owner.StopDrag(c)
							b.MoveCards(c, p)
						} else {
							c.owner.CancelDrag(c)
						}
					}
				}
				b.stroke = nil
			} else {
				dx, dy := b.stroke.PositionDiff()
				c.owner.DragTailBy(dx, dy)
			}
		case *Pile:
			p := v
			if b.stroke.IsReleased() {
				if b.stroke.IsTapped() {
					b.PileTapped(p)
				}
				b.stroke = nil
			}
		default:
			log.Fatal("unknown type of dragging object")
		}
	}

	for _, p := range b.Piles {
		p.Update()
	}

	CTQ.Update()

	return nil
}

// Draw renders the baize into the screen
func (b *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	for _, p := range b.Piles {
		p.Draw(screen)
	}
	for _, p := range b.Piles {
		p.DrawCards(screen)
	}
	for _, p := range b.Piles {
		p.DrawAnimatingCards(screen)
	}
}
