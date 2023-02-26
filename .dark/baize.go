package dark

import (
	"errors"
	"log"

	"oddstream.games/gosol/util"
)

// Baize holds the state of the baize, piles and cards therein.
// Baize is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Baize object as immutable, hence the unexported fields and getters.
type Baize struct {
	variant    string
	script     scripter
	pack       []Card // needed for undo/Pile.UpdateFromSavable
	powerMoves bool
	cardCount  int
	// undoStack
	// statistics (for all variants)
	piles     []*Pile // needed by LIGHT to display piles and cards
	recycles  int     // needed by LIGHT to determine Stock rune
	bookmark  int     // needed by LIGHT to grey out goto bookmark menu item
	moves     int     // count of all available card moves
	fmoves    int     // count of available moves to foundation
	undoStack []*savableBaize
}

func (d *dark) NewBaize(variant string) (*Baize, error) {
	var b *Baize = &Baize{variant: variant}
	var ok bool
	if b.script, ok = variants[variant]; !ok {
		return nil, errors.New("unknown variant " + variant)
	}
	return b, nil
}

// LoadBaize loads variant.saved.json
func (d *dark) LoadBaize(variant string) (*Baize, error) {
	// var b *Baize = &Baize{variant: variant}
	// if undoStack := b.loadUndoStack(); b.IsSavableStackOk(undoStack) {
	// 	b.SetUndoStack(undoStack)
	// }
	return nil, errors.New("not implemented")
}

// Baize public interface ////////////////////////////////////////////////////////////

func (b *Baize) Bookmark() int {
	return b.bookmark
}

func (b *Baize) Complete() bool {
	return false
}

func (b *Baize) Conformant() bool {
	return false
}

func (b *Baize) LoadPosition() (bool, error) {
	if b.bookmark == 0 || b.bookmark > len(b.undoStack) {
		return false, errors.New("No bookmark")
	}
	if b.Complete() {
		return false, errors.New("Cannot undo a completed game") // otherwise the stats can be cooked
	}
	var sav *savableBaize
	var ok bool
	for len(b.undoStack)+1 > b.bookmark {
		sav, ok = b.undoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.updateFromSavable(sav)
	b.undoPush() // replace current state
	b.findDestinations()
	return true, nil
}

func (b *Baize) PercentComplete() int {
	var pairs, unsorted, percent int
	for _, p := range b.piles {
		if p.Len() > 1 {
			pairs += p.Len() - 1
		}
		unsorted += p.vtable.UnsortedPairs()
	}
	percent = (int)(100.0 - util.MapValue(float64(unsorted), 0, float64(pairs), 0.0, 100.0))
	return percent
}

func (b *Baize) Piles() []*Pile {
	return b.piles
}

func (b *Baize) PileTapped(pile *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) Recycles() int {
	return b.recycles
}

func (b *Baize) NewDeal() (bool, error) {
	// a virgin game has one state on the undo stack
	if len(b.undoStack) > 1 && !b.Complete() {
		percent := b.PercentComplete()
		theDark.statistics.RecordLostGame(b.variant, percent)
	}

	b.reset()

	for _, p := range b.piles {
		p.reset()
	}

	// Stock.Fill() needs parameters
	b.cardCount = b.script.Stock().fill(b.script.Packs(), b.script.Suits())
	b.script.Stock().shuffle()
	b.script.StartGame()
	b.undoPush()
	b.findDestinations()
	return true, nil
}

func (b *Baize) RestartDeal() (bool, error) {
	if b.Complete() {
		return false, errors.New("Cannot restart a completed game") // otherwise the stats can be cooked
	}
	var sav *savableBaize
	var ok bool
	for len(b.undoStack) > 0 {
		sav, ok = b.undoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.updateFromSavable(sav)
	b.bookmark = 0 // do this AFTER UpdateFromSavable
	b.undoPush()   // replace current state
	b.findDestinations()
	return true, nil
}

func (b *Baize) SaveGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) SavePosition() (bool, error) {
	if b.Complete() {
		return false, errors.New("Cannot bookmark a completed game") // otherwise the stats can be cooked
	}
	b.bookmark = len(b.undoStack)
	sb := b.undoPeek()
	sb.Bookmark = b.bookmark
	sb.Recycles = b.recycles
	return true, nil
}

func (b *Baize) SetPowerMoves(value bool) {
	b.powerMoves = value
}

func (b *Baize) Statistics() []string {
	return []string{}
}

func (b *Baize) TailDragged(src *Pile, card *Card, dst *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) TailTapped(pile *Pile, card *Card) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) Undo() (bool, error) {
	if len(b.undoStack) < 2 {
		return false, errors.New("Nothing to undo")
	}
	if b.Complete() {
		return false, errors.New("Cannot undo a completed game") // otherwise the stats can be cooked
	}
	_, ok := b.undoPop() // removes current state
	if !ok {
		log.Panic("error popping current state from undo stack")
	}

	sav, ok := b.undoPop() // removes previous state for examination
	if !ok {
		log.Panic("error popping second state from undo stack")
	}
	b.updateFromSavable(sav)
	b.undoPush() // replace current state
	b.findDestinations()
	return true, nil
}

func (b *Baize) UndoStackSize() int {
	return len(b.undoStack)
}

func (b *Baize) Moves() (int, int) {
	return b.moves, b.fmoves
}

// Baize private functions

func (b *Baize) reset() {
	b.undoStack = nil
	b.bookmark = 0
}

func (b *Baize) addPile(pile *Pile) {
	b.piles = append(b.piles, pile)
}

func (b *Baize) calcPowerMoves(pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range b.piles {
		if p.Empty() {
			switch p.vtable.(type) {
			case *Cell:
				emptyCells++
			case *Tableau:
				if p.Label() == "" && p != pDraggingTo {
					// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
					emptyCols++
				}
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	// println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
}

func (b *Baize) findHomesForTail(tail []*Card) []*Pile {
	var homes []*Pile

	var card = tail[0]
	var src = card.owner()
	// can the tail be moved in general?
	if ok, _ := src.canMoveTail(tail); !ok {
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

// foreachCard applys a function to each card
func (b *Baize) foreachCard(fn func(*Card)) {
	for _, p := range b.piles {
		for _, c := range p.cards {
			fn(c)
		}
	}
}

// findAllMovableTails helper function for findDestinations.
func (b *Baize) findAllMovableTails() []*movableTail {
	var tails = []*movableTail{}
	for _, p := range b.piles {
		var t2 []*movableTail = p.vtable.MovableTails()
		if len(t2) > 0 {
			tails = append(tails, t2...)
		}
	}
	return tails
}

// FindDestinations sets Baize.moves, Baize.fmoves, Card.destinations
func (b *Baize) findDestinations() {
	b.moves, b.fmoves = 0, 0

	// Golang gotcha:
	// Go uses a copy of the value instead of the value itself within a range clause.
	// fine for pointers, be careful with objects
	// for _, c := range CardLibrary {
	// 	c.movable = false
	// }
	// https://medium.com/@betable/3-go-gotchas-590b8c014e0a
	b.foreachCard(func(c *Card) { c.tapDestination = nil; c.tapWeight = 0 })

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
		src := card.owner()
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
				} else if dst.peek().Suit() == card.Suit() {
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
}
