package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"
)

type Tableau struct {
	pile *Pile
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Pile {
	pile := newPile("Tableau", slot, fanType, moveType)
	pile.vtable = &Tableau{pile: pile}
	return pile
}

func (self *Tableau) CanAcceptTail(tail []*Card) (bool, error) {
	// AnyCardsProne check done by pile.CanMoveTail
	// checking at this level probably isn't needed
	// if AnyCardsProne(tail) {
	// 	return false, errors.New("Cannot add a face down card")
	// }

	// kludge
	// we couldn't check MOVE_PLUS_ONE in pile.CanMoveTail
	// because we didn't then know the destination pile
	// which we need to know to calculate power moves
	if self.pile.moveType == MOVE_ONE_PLUS {
		if theDark.baize.powerMoves {
			moves := theDark.baize.calcPowerMoves(self.pile)
			if len(tail) > moves {
				if moves == 1 {
					return false, fmt.Errorf("Space to move 1 card, not %d", len(tail))
				} else {
					return false, fmt.Errorf("Space to move %d cards, not %d", moves, len(tail))
				}
			}
		} else {
			if len(tail) > 1 {
				return false, errors.New("Cannot add more than one card")
			}
		}
	}
	return theDark.baize.script.TailAppendError(self.pile, tail)
}

func (self *Tableau) TailTapped(tail []*Card) {
	self.pile.defaultTailTapped(tail)
}

func (self *Tableau) Conformant() bool {
	// return TheGame.Baize.script.UnsortedPairs(self.pile) == 0
	return self.UnsortedPairs() == 0
}

func (self *Tableau) UnsortedPairs() int {
	return theDark.baize.script.UnsortedPairs(self.pile)
}

func (self *Tableau) MovableTails() []*movableTail {
	var tails []*movableTail = []*movableTail{}
	if self.pile.Len() > 0 {
		for _, card := range self.pile.cards {
			var tail = self.pile.makeTail(card)
			if ok, _ := self.pile.canMoveTail(tail); ok {
				if ok, _ := theDark.baize.script.TailMoveError(tail); ok {
					var homes []*Pile = theDark.baize.findHomesForTail(tail)
					for _, home := range homes {
						tails = append(tails, &movableTail{dst: home, tail: tail})
					}
				}
			}
		}
	}
	return tails
}
