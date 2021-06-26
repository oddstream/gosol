package sol

import "oddstream.games/gosol/util"

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
		if p.CanAcceptTail(tail, false) {
			homes = append(homes, p)
		}
	}
	return homes
}

func setMovable(dst *Pile, tail []*Card) {
	c0 := tail[0]
	var m int
	switch {
	case dst.Class == "Foundation":
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
		switch p.Class {
		case "Stock":
			// only the top card is movable
			if !p.Empty() {
				b.movableCards++
			}
		case "StockSpider", "StockScorpion":
			if !p.Empty() {
				numTabs, _ := b.countPiles("Tableau")
				b.movableCards += util.Min(numTabs, p.CardCount())
			}
		case "Waste", "Reserve":
			// just check top card
			if c := p.Peek(); c != nil && !c.Prone() {
				if tail := p.DraggableTail(c); tail != nil {
					if homes := b.NewHomesForTail(tail); len(homes) > 0 {
						b.movableCards++
						for _, dst := range homes {
							setMovable(dst, tail)
						}
					}
				}
			}
		case "Tableau", "Cell":
			for _, c := range p.Cards {
				if c.Prone() {
					continue
				}
				if tail := p.DraggableTail(c); tail != nil {
					if homes := b.NewHomesForTail(tail); len(homes) > 0 {
						b.movableCards++
						for _, dst := range homes {
							setMovable(dst, tail)
						}
					}
				}
			}
		}
	}
}
