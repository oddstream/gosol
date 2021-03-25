package sol

// DraggableTail indicates if a tail from this card can be dragged or not with triggering any visible changes
func (p *Pile) DraggableTail(c *Card) []*Card {
	tail := p.makeTail(c)
	if p.dragFlags&1 == 1 && len(tail) > 1 {
		return nil
	}
	if !isTailConformant(p.dragRules, p.dragFlags, tail) {
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
		if p.CardCount() == 0 && p.localAccept == 0 {
			continue
		}
		if p.CanAcceptCard(c) {
			return p
		}
	}
	return nil
}

func (b *Baize) IsNewHomeForTail(tail []*Card) *Pile {
	for _, p := range b.Piles {
		if p == tail[0].owner {
			continue
		}
		if p.Class == "Cell" {
			continue
		}
		if p.CardCount() == 0 && p.localAccept == 0 {
			continue
		}
		if p.CanAcceptTail(b.Piles, tail, true) {
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

func (b *Baize) MarkMovable() {

	if !TheUserData.HighlightMovable {
		return
	}

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
					}
				} else {
					break
				}
			}
		}
	}
}
