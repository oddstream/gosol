package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Discard struct {
	Core
}

func NewDiscard(slot image.Point, fanType FanType) *Discard {
	discard := &Discard{Core: NewCore("Discard", slot, FAN_NONE, MOVE_NONE)}
	TheBaize.AddPile(discard)
	return discard
}

func (*Discard) CanAcceptCard(card *Card) (bool, error) {
	return false, errors.New("Cannot move a single card to a Discard")
}

func (self *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	if len(tail) != len(CardLibrary)/len(TheBaize.script.Discards()) {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	return TheBaize.script.TailMoveError(tail) // check cards are conformant
}

func (*Discard) TailTapped([]*Card) {
	// do nothing
}

func (*Discard) Collect() {
	// over-ride Core collect to do nothing
}

func (*Discard) Conformant() bool {
	// no Baize that contains any discard piles should be Conformant,
	// because there is no use showing the collect all FAB
	// because that would do nothing
	// because cards are not collected to discard piles
	return false
}

func (self *Discard) Complete() bool {
	if self.Empty() {
		return true
	}
	if self.Len() == len(CardLibrary)/len(TheBaize.script.Discards()) {
		return true
	}
	return false
}

func (*Discard) UnsortedPairs() int {
	// you can only put a sorted sequence into a Discard, so this will always be zero
	return 0
}
