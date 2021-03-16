package sol

/*
--[[
  "Microsoft FreeCell or FreeCell Pro only plays an available card to its
  homecell automatically when all of the lower-ranked cards of the opposite color
  are already on the homecells (except that a two is played if the corresponding
  ace is on its homecell); aces are always played when available. This is one
  version of what can be called safe autoplay"
]]
*/

func (b *Baize) safeCheck(c *Card, dst *Pile) bool {
	// we already know that dst can accept the card, so don't need to check if dst is empty
	localSuit := c.owner.buildRules / 10
	// this is only really for localSuit == 4 (alternate colors)
	if localSuit != 4 {
		return true
	}
	// TODO are all the lower ranked cards of the opposite color to c already on the Foundations?
	for _, p := range b.Piles {
		if p.Class == "Foundation" {
			if p == dst {
				continue
			}
			fc := p.Peek()
			if fc == nil {
				continue
			}
			if fc.Color() == c.Color() {
				continue
			}
			if fc.Ordinal() < c.Ordinal()-1 {
				return false
			}
		}
	}
	return true
}

func (b *Baize) collectFromPile(src *Pile, dst *Pile) int {
	var count int
	for { // collect as many as possible from this pile (think Limited)
		c := src.Peek()
		if c == nil {
			break
		}
		if dst.CanAcceptCard(c) && b.safeCheck(c, dst) {
			b.MoveCards(c, dst)
			count++
		} else {
			break
		}
	}
	return count
}

// Collect automatically moves cards to the Foundations
func (b *Baize) Collect() {

	var count, totalCount int
	for {
		count = 0
		for _, fp := range b.Piles {
			switch fp.Class {
			case "Foundation":
				for _, p := range b.Piles {
					if p.Class == "Tableau" || p.Class == "Cell" || p.Class == "Waste" {
						count += b.collectFromPile(p, fp)
					}
				}
			case "FoundationSpider":
				if fp.CardCount() == 0 {
					for _, p := range b.Piles {
						if p.Class == "Tableau" && p.CardCount() >= 13 {
							for i := 0; i < p.CardCount(); i++ {
								c := p.Cards[i]
								tail := p.makeTail(c)
								if len(tail) == 13 && isConformant(p.buildRules, p.buildFlags, tail) {
									b.MoveCards(c, fp)
									count += 13
									goto NextFoundationPile
								}
							}
						}
					}
				}
			}
		}
		if count == 0 {
			break
		}
	NextFoundationPile:
		totalCount += count
	}

	if totalCount != 0 {
		b.AfterUserMove()
	}
}

// TableauxComplete returns true if every tableau is complete
// func (b *Baize) TableauxComplete() bool {
// 	for _, p := range b.Piles {
// 		if strings.HasPrefix(p.Class, "Tableau") {
// 			if c0 := p.Peek(); c0 != nil {
// 				tail := p.makeTail(c0)
// 				if !isConformant(p.buildRules, p.buildFlags, tail) {
// 					return false
// 				}
// 			}
// 		}
// 	}
// 	return true
// }

// Complete returns true if this game is complete
func (b *Baize) Complete() bool {
	if b.State == Complete {
		println("testing a complete game for completeness")
		return true
	}
	complete := true
	for _, p := range b.Piles {
		if !p.Complete() {
			complete = false
			break
		}
	}
	return complete
}

// Conformant returns true if stock is empty and all piles are conformant
func (b *Baize) Conformant() bool {
	stock := b.findPilePrefix("Stock")
	if len(stock.Cards) > 0 {
		return false
	}
	for _, p := range b.Piles {
		// no need to exclude Foundation* piles as they are guaranteed to be conformant
		// Cells will always be conformant, Reserve is unlikely to be
		if !p.Conformant() {
			return false
		}
	}
	return true
}
