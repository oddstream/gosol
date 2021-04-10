package sol

import "oddstream.games/gosol/util"

// DraggableTail indicates if a tail from this card can be dragged or not without triggering any visible changes
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

func (b *Baize) NewHomesForCard(c *Card) []*Pile {
	homes := []*Pile{}
	for _, p := range b.Piles {
		if p == c.owner {
			continue
		}
		// if p.Class == "Cell" {
		// 	continue
		// }
		// if p.CardCount() == 0 && (p.localAccept == 0 || p.localAccept == c.Ordinal()) {
		// 	continue
		// }
		if p.CanAcceptCard(c) {
			homes = append(homes, p)
		}
	}
	return homes
}

func (b *Baize) NewHomesForTail(tail []*Card) []*Pile {
	homes := []*Pile{}
	for _, p := range b.Piles {
		c0 := tail[0]
		if p == c0.owner {
			continue
		}
		// if p.Class == "Cell" && len(tail) == 1 {
		// 	continue
		// }
		// if p.CardCount() == 0 && (p.localAccept == 0 || p.localAccept == c0.Ordinal()) {
		// 	continue
		// }
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
	case dst.CardCount() == 0 && len(tail) == c0.owner.CardCount():
		m = 1
	case dst.CardCount() == 0 && dst.localAccept == 0:
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
		case "Stock", "StockSpider", "StockScorpion":
			b.movableCards += p.CardCount()
		case "Waste", "Reserve":
			// just check top card
			if c := p.Peek(); c != nil && !c.Prone() {
				if tail := p.DraggableTail(c); tail != nil {
					if homes := b.NewHomesForCard(c); len(homes) > 0 {
						b.movableCards++
						for _, dst := range homes {
							setMovable(dst, tail)
						}
					}
				}
			}
		case "Tableau", "Cell":
			// check all cards upwards until finding an undraggable one
			// for i := len(p.Cards) - 1; i >= 0; i-- {
			// 	c := p.Cards[i]
			// 	if c.Prone() {
			// 		continue
			// 	}
			// 	if tail := p.DraggableTail(c); tail != nil {
			// 		if dst := b.IsNewHomeForTail(tail); dst != nil {
			// 			setMovable(dst, tail)
			// 			b.movableCards++
			// 		}
			// 	} else {
			// 		break
			// 	}
			// }
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
