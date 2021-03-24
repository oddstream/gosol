package sol

func (b *Baize) IsNewHomeForCard(c *Card) bool {
	for _, p := range b.Piles {
		if p == c.owner {
			continue
		}
		if p.CanAcceptCard(c) {
			return true
		}
	}
	return false
}

func (b *Baize) IsNewHomeForTail(tail []*Card) bool {
	for _, p := range b.Piles {
		if p == tail[0].owner {
			continue
		}
		if p.CanAcceptTail(b.Piles, tail) {
			return true
		}
	}
	return false
}

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
			if c := p.Peek(); c != nil {
				if tail := p.DraggableTail(c); tail != nil {
					if b.IsNewHomeForCard(c) {
						c.SetMovable(true)
					}
				}
			}
		case "Tableau", "Cell":
			// check all cards upwards until finding an undraggable one
			for i := len(p.Cards) - 1; i >= 0; i-- {
				c := p.Cards[i]
				if tail := p.DraggableTail(c); tail != nil {
					if b.IsNewHomeForTail(tail) {
						c.SetMovable(true)
					}
				} else {
					break
				}
			}
		}
	}
}
