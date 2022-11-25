package sol

func (b *Baize) FindHomeForTail(owner *Pile, tail []*Card) *Pile {
	if len(tail) == 1 {
		for _, dst := range b.script.Foundations() {
			if dst == owner {
				continue
			}
			if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
				return dst
			}
		}
		for _, dst := range b.script.Cells() {
			if dst == owner {
				continue
			}
			if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
				return dst
			}
		}
	}
	for _, dst := range b.script.Tableaux() {
		if dst == owner {
			continue
		}
		if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
			return dst
		}
	}
	return nil
}

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

	for _, dst := range pilesToCheck {
		if dst == src {
			continue
		}
		if ok, _ := dst.vtable.CanAcceptTail(tail); ok {
			homes = append(homes, dst)
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

// CountMoves sets Baize.moves, Baize.fmoves, Card.movable
// 	0 - can't move, or pointless move
//	1 - weak move (conformant with card above)
//	2 - move to cell or empty pile
//	3 - move
//	4 - move to foundation
func (b *Baize) CountMoves() {
	b.moves, b.fmoves = 0, 0
	MarkAllCardsImmovable()
	if b.script.Stock().Empty() {
		if b.Recycles() > 0 {
			b.moves++
		}
	} else {
		b.moves++
		b.script.Stock().Peek().movable = true
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
			card.movable = true
			if dst.category == "Foundation" {
				b.fmoves++
			}
			// that's it, unless Card.movable changes from bool to int (like lsol)
		}
	}
}
