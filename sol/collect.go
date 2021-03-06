package sol

import (
	"strings"
)

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
			if fc.red == c.red {
				continue
			}
			if fc.ordinal < c.ordinal-1 {
				return false
			}
		}
	}
	return true
}

func (b *Baize) collectFromPile(src *Pile, dst *Pile) int {
	var count int
	c := src.Peek()
	if c != nil {
		if dst.CanAcceptCard(c) && b.safeCheck(c, dst) {
			b.MoveCards(c, dst)
			count++
		}
	}
	return count
}

// Collect automatically moves cards to the Foundations
func (b *Baize) Collect() {

	var count int
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
	}

}

// TableauxComplete returns true if every tableau is complete
func (b *Baize) TableauxComplete() bool {
	for _, p := range b.Piles {
		if strings.HasPrefix(p.Class, "Tableau") {
			if c0 := p.Peek(); c0 != nil {
				tail := p.makeTail(c0)
				if !isConformant(p.buildRules, p.buildFlags, tail) {
					return false
				}
			}
		}
	}
	return true
}

// Complete returns true if this game is complete
func (b *Baize) Complete() bool {
	complete := true
	for _, p := range b.Piles {
		if !p.IsComplete() {
			complete = false
			break
		}
	}
	return complete
}
