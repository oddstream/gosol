package dark

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Stock struct {
	pile *Pile
}

func NewStock(slot image.Point, fanType FanType, packs int, suits int, cardFilter *[14]bool, jokersPerPack int) *Pile {
	pile := newPile("Stock", slot, fanType, MOVE_ONE)
	pile.vtable = &Stock{pile: pile}
	theDark.baize.cardCount = pile.fill(packs, suits)
	pile.shuffle()
	return pile
}

func (*Stock) CanAcceptTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) TailTapped([]*Card) {
	// do nothing, handled by script, which had first dibs
}

func (self *Stock) Conformant() bool {
	return self.pile.Empty()
}

// UnsortedPairs - cards in a stock pile are always considered to be unsorted
func (self *Stock) UnsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Stock) MovableTails() []*movableTail {
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
