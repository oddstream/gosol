package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Reserve struct {
	pile *Pile
}

func NewReserve(slot image.Point, fanType FanType) *Pile {
	pile := newPile("Reserve", slot, fanType, MOVE_ONE)
	pile.vtable = &Reserve{pile: pile}
	return pile
}

func (*Reserve) CanAcceptTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (self *Reserve) TailTapped(tail []*Card) {
	self.pile.defaultTailTapped(tail)
}

// Conformant when contains zero or one card(s), same as Waste
func (self *Reserve) Conformant() bool {
	return self.pile.Len() < 2
}

// UnsortedPairs - cards in a reserve pile are always considered to be unsorted
func (self *Reserve) UnsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Reserve) MovableTails() []*movableTail {
	// nb same as Cell.MovableTails
	var tails []*movableTail = []*movableTail{}
	if self.pile.Len() > 0 {
		var card *Card = self.pile.peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = theDark.baize.findHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &movableTail{dst: home, tail: tail})
		}
	}
	return tails
}
