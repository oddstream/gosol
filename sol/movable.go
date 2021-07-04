package sol

import (
	"strings"

	"oddstream.games/gosol/util"
)

// DraggableTail indicates if a tail from this card can be dragged or not without triggering any visible changes
func (p *Pile) DraggableTail(c *Card) []*Card {
	tail := p.makeTail(c)
	if p.Flags&DragFlagSingle == DragFlagSingle {
		if ThePreferences.PowerMoves && p.Class == "Tableau" {
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

func (b *Baize) NewHomesForTail(tail []*Card) []*Pile {
	var c0 *Card = tail[0]
	var homes []*Pile
	for _, p := range b.Piles {
		if p == c0.owner {
			continue
		}
		if ok, _ := p.driver.CanAcceptTail(tail); ok {
			homes = append(homes, p)
		}
	}
	return homes
}

func setMovable(dst *Pile, tail []*Card) {
	c0 := tail[0]
	var m int
	switch {
	case strings.HasPrefix(dst.Class, "Foundation"):
		m = 3
	case dst.Class == "Cell":
		m = 1
	case dst.Empty() && len(tail) == c0.owner.CardCount():
		// moving an entire pile to another empty pile
		m = 0
	case dst.Empty() && len(tail) == 1 && dst.localAccept == 0 && c0.owner.Class == dst.Class:
		// moving a single card to an empty pile of the same type
		m = 1
	default:
		m = 2
	}
	c0.movable = util.Max(c0.movable, m)
}

func (b *Baize) MarkMovable() {

	b.movableCards = 0

	for _, p := range b.Piles {
		for _, c := range p.Cards {
			c.movable = 0
		}
		b.movableCards += p.driver.Movable()
	}
}

func genericMovableTopCard(p *Pile) int {
	var count int
	// just check top card
	if c := p.Peek(); c != nil && !c.Prone() {
		if tail := p.DraggableTail(c); tail != nil {
			if homes := TheBaize.NewHomesForTail(tail); len(homes) > 0 {
				count++
				for _, dst := range homes {
					setMovable(dst, tail)
				}
			}
		}
	}
	return count
}

func genericMovableAllCards(p *Pile) int {
	var count int
	for _, c := range p.Cards {
		if c.Prone() {
			continue
		}
		if tail := p.DraggableTail(c); tail != nil {
			if homes := TheBaize.NewHomesForTail(tail); len(homes) > 0 {
				count++
				for _, dst := range homes {
					setMovable(dst, tail)
				}
			}
		}
	}
	return count
}

func (c *Cell) Movable() int             { return genericMovableTopCard(c.parent) }
func (f *Foundation) Movable() int       { return 0 }
func (f *FoundationSpider) Movable() int { return 0 }
func (g *Golf) Movable() int             { return genericMovableTopCard(g.parent) }
func (r *Reserve) Movable() int          { return genericMovableTopCard(r.parent) }
func (s *Stock) Movable() int {
	if s.parent.Empty() {
		return 0
	}
	return 1
}
func (s *StockCruel) Movable() int    { return 0 }
func (s *StockScorpion) Movable() int { return len(s.parent.Cards) }
func (s *StockSpider) Movable() int   { return len(s.parent.Cards) }
func (t *Tableau) Movable() int       { return genericMovableAllCards(t.parent) }
func (t *TableauSpider) Movable() int { return genericMovableAllCards(t.parent) }
func (w Waste) Movable() int          { return genericMovableTopCard(w.parent) }
