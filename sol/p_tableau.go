package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gosol/util"
)

type Tableau struct {
	Core
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Tableau {
	tableau := &Tableau{Core: NewCore("Tableau", slot, fanType, moveType)}
	TheBaize.AddPile(tableau)
	return tableau
}

func (self *Tableau) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	var tail []*Card = []*Card{card}
	return TheBaize.script.TailAppendError(self, tail)
}

func powerMoves(piles []Pile, pDraggingTo Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range piles {
		if p.Empty() {
			switch (p).(type) {
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

func (self *Tableau) CanAcceptTail(tail []*Card) (bool, error) {
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card")
	}
	if self.MoveType() == MOVE_ONE_PLUS {
		if ThePreferences.PowerMoves {
			moves := powerMoves(TheBaize.piles, self)
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
	return TheBaize.script.TailAppendError(self, tail)
}

// use Core.TailTapped

// use Core.Collect

func (self *Tableau) Conformant() bool {
	return TheBaize.script.UnsortedPairs(self) == 0
}

func (self *Tableau) Complete() bool {
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
	if self.Empty() {
		return true
	}
	if len(TheBaize.script.Discards()) > 0 {
		if self.Len() == len(CardLibrary)/len(TheBaize.script.Discards()) {
			// eg 13 == 52 / 4
			if TheBaize.script.UnsortedPairs(self) == 0 {
				return true
			}
		}
	}
	return false
}

func (self *Tableau) UnsortedPairs() int {
	return TheBaize.script.UnsortedPairs(self)
}
