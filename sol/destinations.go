package sol

import "sort"

func (b *Baize) FindHomesForTail(tail []*Card) []*Pile {
	var homes []*Pile

	var card = tail[0]
	var src = card.owner
	// can the tail be moved in general?
	if ok, _ := src.CanMoveTail(tail); !ok {
		return homes
	}

	// is the tail conformant enough to move?
	if ok, _ := b.script.TailMoveError(tail); !ok {
		return homes
	}

	var pilesToCheck []*Pile = []*Pile{}
	pilesToCheck = append(pilesToCheck, b.script.Foundations()...)
	pilesToCheck = append(pilesToCheck, b.script.Tableaux()...)
	pilesToCheck = append(pilesToCheck, b.script.Cells()...)
	if b.script.Waste() != nil {
		// in Go 1.19, append will add a nil
		// in Go 1.17, nil was not appended
		pilesToCheck = append(pilesToCheck, b.script.Waste())
	}

	for _, dst := range pilesToCheck {
		if !dst.Valid() {
			println("Destination pile not valid", dst)
		}
		if dst != src {
			if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
				homes = append(homes, dst)
			}
		}
	}

	return homes
}

func (b *Baize) findAllMovableTails() []*MovableTail {
	var tails = []*MovableTail{}
	for _, p := range b.piles {
		var t2 []*MovableTail = p.vtable.MovableTails()
		if len(t2) > 0 {
			tails = append(tails, t2...)
		}
	}
	return tails
}

// func isWeakMove(src *Pile, card *Card) bool {
// 	return false
// }

// FindDestinations sets Baize.moves, Baize.fmoves, Card.destinations
func (b *Baize) FindDestinations() {
	b.moves, b.fmoves = 0, 0

	MarkAllCardsImmovable()

	if !b.script.Stock().Hidden() {
		if b.script.Stock().Empty() {
			if b.Recycles() > 0 {
				b.moves++
			}
		} else {
			// games like Agnes B (with a Spiker-like stock) need to report an available move
			b.moves += 1
			// card := b.script.Stock().Peek()
			// card.destinations = b.FindHomesForTail([]*Card{card})
			// b.moves += len(card.destinations)
		}
	}

	for _, mc := range b.findAllMovableTails() {
		movable := true
		card := mc.tail[0]
		src := card.owner
		dst := mc.dst
		// moving an full tail from one pile to another empty pile is pointless
		if dst.Len() == 0 && len(mc.tail) == len(src.cards) {
			if src.label == dst.label && src.category == dst.category {
				movable = false
			}
		}
		if movable {
			b.moves++
			card.destinations = append(card.destinations, mc.dst)
			if dst.category == "Foundation" {
				b.fmoves++
			}
		}
	}
}

type PileAndWeight struct {
	pile   *Pile
	weight int
}

func (b *Baize) BestDestination(card *Card, destinations []*Pile) *Pile {
	var paw []*PileAndWeight
	for _, dst := range destinations {
		var tmp PileAndWeight = PileAndWeight{pile: dst, weight: len(dst.cards)}

		switch dst.category {
		case "Foundation":
			tmp.weight += 52 // magic number, sorry
		case "Tableau":
			if len(dst.cards) > 0 {
				if card.Suit() == dst.Peek().Suit() {
					tmp.weight += 26 // magic number, sorry
				}
			}
		}
		paw = append(paw, &tmp)
	}
	sort.Slice(paw, func(i, j int) bool { return paw[i].weight > paw[j].weight })
	return paw[0].pile
}
