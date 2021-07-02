package sol

import (
	"fmt"
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

func (b *Baize) collectFromPile(src *Pile, dst *Pile) int {
	var count int
	for { // collect as many as possible from this pile (think Limited)
		c := src.Peek()
		if c == nil {
			break
		}
		if dst.CanAcceptCard(c) && b.safeCheck(c, dst) {
			dst.MoveCards(c)
			count++
		} else {
			break
		}
	}
	return count
}

// Collect automatically moves cards to the Foundations
func (b *Baize) Collect() {

	var foundations []*Pile
	for _, fp := range b.Piles {
		if strings.HasPrefix(fp.Class, "Foundation") {
			foundations = append(foundations, fp)
		}
	}
	var count, totalCount int
	for {
		count = 0
		// iterate over foundations and pull cards to them
		for _, fp := range foundations {
			if fp.Spider() && fp.Empty() {
				for _, p := range b.Piles {
					if p.Class == "Tableau" && p.CardCount() >= 13 {
						for i := 0; i < p.CardCount(); i++ {
							c := p.Cards[i]
							tail := p.makeTail(c)
							// pearl from the mudbank:
							// mistress mop may have a run of 13 cards, in numerical order (which are conformant in a Tableau)
							// but these are not conformant for the Foundation
							if len(tail) == 13 && isTailConformant(p.Build, p.Flags, tail) {
								if ok, _ := fp.driver.CanAcceptTail(tail); ok {
									fp.MoveCards(c)
									count += 13
									goto NextFoundationPile
								}
							}
						}
					}
				}
			} else {
				for _, p := range b.Piles {
					if p.Class == "Tableau" || p.Class == "Cell" || p.Class == "Waste" || p.Class == "Reserve" || p.Class == "Golf" {
						count += b.collectFromPile(p, fp)
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
