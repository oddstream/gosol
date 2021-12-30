package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Waste struct {
	Core
}

func NewWaste(slot image.Point, fanType FanType) *Waste {
	waste := &Waste{Core: NewCore("Waste", slot, fanType, MOVE_ONE)}
	TheBaize.AddPile(waste)
	return waste
}

func (*Waste) CanAcceptCard(card *Card) (bool, error) {
	if !card.Owner().IsStock() {
		return false, errors.New("Waste can only accept cards from the Stock")
	}
	return true, nil
}

func (*Waste) CanAcceptTail(tail []*Card) (bool, error) {
	if !tail[0].Owner().IsStock() {
		return false, errors.New("Waste can only accept cards from the Stock")
	}
	return true, nil
}

func (self *Waste) TailTapped(tail []*Card) {
	GenericTailTapped(self, tail)
}

func (self *Waste) Collect() {
	GenericCollect(self)
}

func (self *Waste) Conformant() bool {
	// zero or one cards would be conformant
	return self.Len() < 2
}

func (self *Waste) Complete() bool {
	return self.Empty()
}

func (self *Waste) UnsortedPairs() int {
	// Waste (like Stock, Reserve) is always considered unsorted
	if self.Empty() {
		return 0
	}
	return self.Len() - 1
}
