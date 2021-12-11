package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Foundation struct {
	Base
}

func NewFoundation(slot image.Point, fanType FanType) *Foundation {
	f := &Foundation{}
	f.Ctor(f, "Foundation", slot, fanType)
	return f
}

func (*Foundation) CanMoveTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot move cards from a Foundation")
}

func (f *Foundation) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if len(f.cards) == len(TheBaize.cardLibrary)/len(TheBaize.foundations) {
		return false, errors.New("The Foundation is full")
	}
	var tail []*Card = []*Card{card}
	// pearl from the mudbank cannot pass a *Foundation to script functions, only a *Pile
	return TheBaize.script.TailAppendError(f.iface, tail)
}

func (f *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card")
	}
	return TheBaize.script.TailAppendError(f.iface, tail)
}

func (f *Foundation) TailTapped(tail []*Card) {
	// over-ride base to do nothing
}

func (f *Foundation) Collect() {
	// over-ride base collect to do nothing
}

func (f *Foundation) Conformant() bool {
	return true
}

func (f *Foundation) Complete() bool {
	return len(f.cards) == len(TheBaize.cardLibrary)/len(TheBaize.foundations)
}

func (f *Foundation) UnsortedPairs() int {
	return 0
}
