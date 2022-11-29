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

func (self *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card")
	}
	if self.parent.Len() == len(CardLibrary)/len(TheBaize.script.Foundations()) {
		return false, errors.New("The Foundation is full")
	}
	return TheBaize.script.TailAppendError(self.parent, tail)
}

func (*Foundation) TailTapped([]*Card) {}

func (*Foundation) Collect() {}

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
