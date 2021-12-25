package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Reserve struct {
	pile *Pile
}

func NewReserve(slot image.Point, fanType FanType) *Pile {
	p := &Pile{}
	p.Ctor(&Reserve{pile: p}, "Reserve", slot, fanType, MOVE_ONE)
	return p
}

func (r *Reserve) CanAcceptCard(card *Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (r *Reserve) CanAcceptTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (r *Reserve) TailTapped(tail []*Card) {
	r.pile.GenericTailTapped(tail)
}

func (r *Reserve) Collect() {
	r.pile.GenericCollect()
}

func (r *Reserve) Conformant() bool {
	return r.pile.Len() < 2
}

func (r *Reserve) Complete() bool {
	return r.pile.Empty()
}

func (r *Reserve) UnsortedPairs() int {
	// Reserve (like Stock) is always considered unsorted
	if r.pile.Empty() {
		return 0
	}
	return r.pile.Len() - 1
}
