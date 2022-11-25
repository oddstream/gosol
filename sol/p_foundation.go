package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Foundation struct {
	parent *Pile
}

func NewFoundation(slot image.Point) *Pile {
	foundation := NewPile("Foundation", slot, FAN_NONE, MOVE_NONE)
	foundation.vtable = &Foundation{parent: &foundation}
	TheBaize.AddPile(&foundation)
	return &foundation
}

func (self *Foundation) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if self.parent.Len() == len(CardLibrary)/len(TheBaize.script.Foundations()) {
		return false, errors.New("The Foundation is full")
	}
	var tail []*Card = []*Card{card}
	return TheBaize.script.TailAppendError(self.parent, tail)
}

func (self *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	return TheBaize.script.TailAppendError(self.parent, tail)
}

func (*Foundation) TailTapped([]*Card) {
	// do nothing
}

func (*Foundation) Collect() {
	// do nothing
}

func (*Foundation) Conformant() bool {
	return true
}

func (self *Foundation) Complete() bool {
	return self.parent.Len() == len(CardLibrary)/len(TheBaize.script.Foundations())
}

func (*Foundation) UnsortedPairs() int {
	// you can only put a sorted sequence into a Foundation, so this will always be zero
	return 0
}

func (self *Foundation) MovableTails() []*MovableTail {
	return nil
}
