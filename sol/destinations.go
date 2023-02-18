package sol

func (b *Baize) FindHomesForTail(tail []*Card) []*Pile {
	var homes []*Pile

	var card = tail[0]
	var src = card.Owner()
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
	pilesToCheck = append(pilesToCheck, b.script.Discards()...)
	if b.script.Waste() != nil {
		// in Go 1.19, append will add a nil
		// in Go 1.17, nil was not appended?
		pilesToCheck = append(pilesToCheck, b.script.Waste())
	}

	for _, dst := range pilesToCheck {
		// if !dst.Valid() {
		// 	log.Println("Destination pile not valid", dst)
		// }
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

	// Golang gotcha:
	// Go uses a copy of the value instead of the value itself within a range clause.
	// fine for pointers, be careful with objects
	// for _, c := range CardLibrary {
	// 	c.movable = false
	// }
	// https://medium.com/@betable/3-go-gotchas-590b8c014e0a
	for i := 0; i < len(CardLibrary); i++ {
		CardLibrary[i].tapDestination = nil
		CardLibrary[i].tapWeight = 0
	}

	if !b.script.Stock().Hidden() {
		if b.script.Stock().Empty() {
			if b.Recycles() > 0 {
				b.moves++
			}
		} else {
			// games like Agnes B (with a Spider-like stock) need to report an available move
			// so we can't do this:
			// card := b.script.Stock().Peek()
			// card.destinations = b.FindHomesForTail([]*Card{card})
			// b.moves += len(card.destinations)
			b.moves += 1
		}
	}

	for _, mc := range b.findAllMovableTails() {
		movable := true
		card := mc.tail[0]
		src := card.Owner()
		dst := mc.dst
		// moving an full tail from one pile to another empty pile is pointless
		if dst.Len() == 0 && len(mc.tail) == len(src.cards) {
			if src.label == dst.label && src.category == dst.category {
				movable = false
			}
		}
		if movable {
			b.moves++
			if _, ok := dst.vtable.(*Foundation); ok {
				b.fmoves++
			}
			var weight int
			switch dst.vtable.(type) {
			case *Cell:
				weight = 0
			case *Tableau:
				if dst.Empty() {
					if dst.Label() != "" {
						weight = 1
					} else {
						weight = 0
					}
				} else if dst.Peek().Suit() == card.Suit() {
					// Simple Simon, Spider
					weight = 2
				} else {
					weight = 1
				}
			case *Foundation, *Discard:
				// moves to Foundation get priority when card is tapped
				weight = 3
			default:
				weight = 0
			}
			if card.tapDestination == nil || weight > card.tapWeight {
				card.tapDestination = dst
				card.tapWeight = weight
			}
		}
	}

	b.UpdateToolbar()
	b.UpdateDrawers()
	b.UpdateStatusbar()

	if !TheSettings.AlwaysShowMovableCards {
		TheSettings.ShowMovableCards = false
	}
}

/*
func (b *Baize) BestDestination(card *Card, destinations []*Pile) *Pile {
	if len(destinations) == 1 {
		return destinations[0]
	}
	var paw []*PileAndWeight
	for _, dst := range destinations {
		var tmp PileAndWeight = PileAndWeight{pile: dst, weight: len(dst.cards)}

		switch dst.vtable.(type) {
		case *Cell:
			tmp.weight -= 1
		case *Foundation:
			tmp.weight += 52 // magic number, sorry
		case *Tableau:
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
*/
