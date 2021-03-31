package sol

import (
	"time"

	"oddstream.games/gosol/util"
)

// DraggableTail indicates if a tail from this card can be dragged or not with triggering any visible changes
func (p *Pile) DraggableTail(c *Card) []*Card {
	tail := p.makeTail(c)
	if p.Flags&DragFlagSingle == DragFlagSingle {
		if TheUserData.PowerMoves {
			pm := powerMoves(TheBaize.Piles, p)
			if len(tail) > pm {
				return nil
			}
		} else {
			if len(tail) > 1 {
				return nil
			}
		}
	}
	if !isTailConformant(p.Drag, p.Flags, tail) {
		return nil
	}
	return tail
}

func (b *Baize) IsNewHomeForCard(c *Card) *Pile {
	for _, p := range b.Piles {
		if p == c.owner {
			continue
		}
		if p.Class == "Cell" {
			continue
		}
		// if p.CardCount() == 0 && (p.localAccept == 0 || p.localAccept == c.Ordinal()) {
		// 	continue
		// }
		if p.CanAcceptCard(c) {
			return p
		}
	}
	return nil
}

func (b *Baize) IsNewHomeForTail(tail []*Card) *Pile {
	for _, p := range b.Piles {
		c0 := tail[0]
		if p == c0.owner {
			continue
		}
		if p.Class == "Cell" && len(tail) == 1 {
			continue
		}
		// if p.CardCount() == 0 && (p.localAccept == 0 || p.localAccept == c0.Ordinal()) {
		// 	continue
		// }
		if p.CanAcceptTail(tail, true) {
			return p
		}
	}
	return nil
}

// func PointlessTailMove(dst *Pile, tail []*Card) bool {
// 	c1 := tail[0]
// 	c2 := dst.Peek()
// 	if c1 != nil && c2 != nil {
// 		if c1.owner.Class == dst.Class {
// 			if c1.Ordinal() == c2.Ordinal() {
// 				if c1.Suit() == c2.Suit() {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

func (b *Baize) HighlightMovable() {

	b.movableCards = 0

	if !TheUserData.HighlightMovable {
		return
	}

	defer util.Duration(time.Now(), "HighlightMovable")

	for _, p := range b.Piles {
		for _, c := range p.Cards {
			c.SetMovable(false)
		}
		switch p.Class {
		case "Waste", "Reserve":
			// just check top card
			if c := p.Peek(); c != nil && !c.Prone() {
				if tail := p.DraggableTail(c); tail != nil {
					if dst := b.IsNewHomeForCard(c); dst != nil {
						c.SetMovable(true)
						b.movableCards++
					}
				}
			}
		case "Tableau", "Cell":
			// check all cards upwards until finding an undraggable one
			for i := len(p.Cards) - 1; i >= 0; i-- {
				c := p.Cards[i]
				if c.Prone() {
					continue
				}
				if tail := p.DraggableTail(c); tail != nil {
					if dst := b.IsNewHomeForTail(tail); dst != nil {
						c.SetMovable(true)
						b.movableCards++
					}
				} else {
					break
				}
			}
		}
	}
}
