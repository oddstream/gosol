package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gomps5/util"
)

type Tableau struct {
	pile *Pile
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Pile {
	p := &Pile{}
	p.Ctor(&Tableau{pile: p}, "Tableau", slot, fanType, moveType)
	return p
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
	if t.pile.moveType == MOVE_ONE_PLUS {
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
	return TheBaize.script.UnsortedPairs(t.pile) == 0
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
	if len(TheBaize.script.Discards()) > 0 {
		if t.pile.Len() == len(CardLibrary)/len(TheBaize.script.Discards()) {
			// eg 13 == 52 / 4
			if TheBaize.script.UnsortedPairs(t.pile) == 0 {
				return true
			}
		}
	}
	return false
}

func (t *Tableau) UnsortedPairs() int {
	return TheBaize.script.UnsortedPairs(t.pile)
}
