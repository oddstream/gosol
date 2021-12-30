package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Reserve struct {
	Core
}

func NewReserve(slot image.Point, fanType FanType) *Reserve {
	reserve := &Reserve{Core: NewCore("Reserve", slot, fanType, MOVE_ONE)}
	TheBaize.AddPile(reserve)
	return reserve
}

func (*Reserve) CanAcceptCard(card *Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (*Reserve) CanAcceptTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (self *Reserve) TailTapped(tail []*Card) {
	GenericTailTapped(self, tail)
}

func (self *Reserve) Collect() {
	GenericCollect(self)
}

func (self *Reserve) Conformant() bool {
	return self.Len() < 2
}

func (self *Reserve) Complete() bool {
	return self.Empty()
}

func (self *Reserve) UnsortedPairs() int {
	// Reserve (like Stock) is always considered unsorted
	if self.Empty() {
		return 0
	}
	return self.Len() - 1
}
