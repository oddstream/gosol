package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Foundation struct {
	pile *Pile
}

func NewFoundation(slot image.Point, fanType FanType) *Pile {
	p := &Pile{}
	p.Ctor(&Foundation{pile: p}, "Foundation", slot, fanType)
	return p
}

func (*Foundation) CanMoveTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot move cards from a Foundation")
}

func (f *Foundation) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if f.pile.Len() == len(CardLibrary)/len(TheBaize.script.Foundations()) {
		return false, errors.New("The Foundation is full")
	}
	var tail []*Card = []*Card{card}
	// pearl from the mudbank cannot pass a *Foundation to script functions, only a *Pile
	return TheBaize.script.TailAppendError(f.pile, tail)
}

func (f *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card")
	}
	return TheBaize.script.TailAppendError(f.pile, tail)
}

func (*Foundation) TailTapped([]*Card) {
	// do nothing
}

func (f *Foundation) Collect() {
	// over-ride base collect to do nothing
}

func (f *Foundation) Conformant() bool {
	return true
}

func (f *Foundation) Complete() bool {
	return f.pile.Len() == len(CardLibrary)/len(TheBaize.script.Foundations())
}

func (f *Foundation) UnsortedPairs() int {
	// you can only put a sorted sequence into a Foundation, so this will always be zero
	return 0
}
