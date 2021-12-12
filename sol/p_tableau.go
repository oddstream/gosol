package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gomps5/util"
)

type Tableau struct {
	pile     *Pile
	moveType MoveType
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Pile {
	p := &Pile{}
	p.Ctor(&Tableau{pile: p, moveType: moveType}, "Tableau", slot, fanType)
	return p
}

func (t *Tableau) CanMoveTail(tail []*Card) (bool, error) {
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	switch t.moveType {
	case MOVE_ANY:
		// well, that was easy
	case MOVE_ONE:
		if len(tail) > 1 {
			return false, errors.New("You can only move one card")
		}
	case MOVE_ONE_PLUS:
		// don't know destination, so we allow this as MOVE_ANY
	case MOVE_ONE_OR_ALL:
		if len(tail) == 1 {
			// that's okay
		} else if len(tail) == t.pile.Len() {
			// that's okay too
		} else {
			return false, errors.New("Only move one card, or the whole pile")
		}
	}
	return TheBaize.script.TailMoveError(tail)
}

func (t *Tableau) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	var tail []*Card = []*Card{card}
	return TheBaize.script.TailAppendError(t.pile, tail)
}

func powerMoves(piles []*Pile, pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range piles {
		if p.Empty() {
			switch (p.subtype).(type) {
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

func (t *Tableau) CanAcceptTail(tail []*Card) (bool, error) {
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card")
	}
	if t.moveType == MOVE_ONE_PLUS {
		if ThePreferences.PowerMoves {
			moves := powerMoves(TheBaize.piles, t.pile)
			if len(tail) > moves {
				if moves == 1 {
					return false, fmt.Errorf("Space to move 1 card, not %d", len(tail))
				} else {
					return false, fmt.Errorf("Space to move %d cards, not %d", moves, len(tail))
				}
			}
		} else {
			if len(tail) > 1 {
				return false, errors.New("Cannot move more than one card")
			}
		}
	}
	return TheBaize.script.TailAppendError(t.pile, tail)
}

func (t *Tableau) TailTapped(tail []*Card) {
	t.pile.GenericTailTapped(tail)
}

func (t *Tableau) Collect() {
	t.pile.GenericCollect()
}

func (t *Tableau) Conformant() bool {
	if t.pile.Len() > 1 {
		return TheBaize.script.UnsortedPairs(t.pile) == 0
	}
	return true
}

func (t *Tableau) Complete() bool {
	/*
	   'complete' means
	       (a) empty
	       (b) if discard piles exist, then there are no unsorted card pairs
	           and pile len == baize->numberOfCardsInLibrary / ndiscards
	           there will always be one discard pile for each suit
	           (number of discard piles is packs * suits)
	           (Simple Simon has 10 tableau piles and 4 discard piles)
	           (Spider has 10 tableau piles and 8 discard piles)
	*/
	if t.pile.Empty() {
		return true
	}
	if len(TheBaize.discards) > 0 {
		if t.pile.Len() == len(TheBaize.cardLibrary)/len(TheBaize.discards) {
			if TheBaize.script.UnsortedPairs(t.pile) == 0 {
				return true
			}
		}
	}
	return false
}

func (t *Tableau) UnsortedPairs() int {
	if t.pile.Len() > 1 {
		return TheBaize.script.UnsortedPairs(t.pile)
	} else {
		return 0
	}
}

func (t *Tableau) Reset() {
	t.pile.GenericReset()
}
