package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Discard struct {
	pile *Pile
}

func NewDiscard(slot image.Point, fanType FanType) *Pile {
	p := &Pile{}
	p.Ctor(&Discard{pile: p}, "Discard", slot, FAN_NONE)
	return p
}

func (*Discard) CanMoveTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot move cards from a Discard")
}

func (*Discard) CanAcceptCard(card *Card) (bool, error) {
	return false, errors.New("Cannot move a single card to a Discard")
}

func (d *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !d.pile.Empty() {
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
	// do nothing
}

func (*Discard) Conformant() bool {
	// no Baize that contains any discard piles should be Conformant,
	// because there is no use showing the collect all FAB
	// because that would do nothing
	// because cards are not collected to discard piles
	return false
}

func (d *Discard) Complete() bool {
	if d.pile.Empty() {
		return true
	}
	if d.pile.Len() == len(CardLibrary)/len(TheBaize.script.Discards()) {
		return true
	}
	return false
}

func (d *Discard) UnsortedPairs() int {
	if d.pile.Len() > 1 {
		return TheBaize.script.UnsortedPairs(d.pile)
	} else {
		return 0
	}
}

func (d *Discard) Reset() {
	d.pile.GenericReset()
}
