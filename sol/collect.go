package sol

import (
	"fmt"
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

/*
func (b *Baize) safeCheck(c *Card, dst *Pile) bool {
	// we already know that dst can accept the card, so don't need to check if dst is empty
	localSuit := c.owner.Build / 10
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
*/

// Collect is like Tapping on the top card of each pile (except Stock), or on a K in a Spider pile
// have a special CardTapped func that only targets b.foundations

func genericCollect(p *Pile) int {

	card := p.Peek()
	if card == nil {
		return 0
	}

	var cardsMoved int
	for _, fp := range TheBaize.foundations {
		if ok, _ := fp.driver.CanAcceptTail([]*Card{card}); ok {
			fp.MoveCards(card)
			cardsMoved++
		}
	}
	return cardsMoved
}

func (c *Cell) Collect() int {
	return genericCollect(c.parent)
}
func (f *Foundation) Collect() int {
	return 0
}
func (f *FoundationSpider) Collect() int {
	return 0
}
func (g *Golf) Collect() int {
	return genericCollect(g.parent)
}
func (r *Reserve) Collect() int {
	return genericCollect(r.parent)
}
func (s *Stock) Collect() int {
	return 0
}
func (s *StockCruel) Collect() int {
	return 0
}
func (s *StockScorpion) Collect() int {
	return 0
}
func (s *StockSpider) Collect() int {
	return 0
}
func (t *Tableau) Collect() int {
	return genericCollect(t.parent)
}
func (t *TableauSpider) Collect() int {

	p := t.parent

	for _, card := range p.Cards {
		if card.Ordinal() == 13 {
			tail := p.makeTail(card)
			if len(tail) == 13 && isTailConformant(p.Build, p.Flags, tail) {
				for _, fp := range TheBaize.foundations {
					// pearl from the mudbank:
					// mistress mop may have a run of 13 cards, in numerical order (which are conformant in a Tableau)
					// but these are not conformant for the Foundation (because the suits differ)
					if ok, _ := fp.driver.CanAcceptTail(tail); ok {
						fp.MoveCards(card)
						return 13
					}
				}
			}
		}
	}
	return 0
}
func (w *Waste) Collect() int {
	return genericCollect(w.parent)
}

// Collect automatically moves cards to the Foundations
func (b *Baize) Collect() {

	var count, totalCount int
	for {
		count = 0
		for _, p := range b.Piles {
			count += p.driver.Collect()
		}
		if count == 0 {
			break
		}
		totalCount += count
	}
	if totalCount > 0 {
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
		fmt.Println("testing a complete game for completeness")
		return true
	}
	for _, p := range b.Piles {
		if !p.driver.Complete() {
			return false
		}
	}
	return true
}

// Conformant returns true if all piles are either empty or all their cards are conformant
func (b *Baize) Conformant() bool {
	for _, p := range b.Piles {
		if !p.driver.Conformant() {
			return false
		}
	}
	return true
}
