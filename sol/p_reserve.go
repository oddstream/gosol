package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Reserve struct {
	pile *Pile
}

func NewReserve(slot image.Point, fanType FanType) *Pile {
	pile := NewPile("Reserve", slot, fanType, MOVE_ONE)
	pile.vtable = &Reserve{pile: &pile}
	TheBaize.AddPile(&pile)
	return &pile
}

func (*Reserve) CanAcceptTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot add a card to a Reserve")
}

func (self *Reserve) TailTapped(tail []*Card) {
	self.pile.DefaultTailTapped(tail)
}

// Conformant when contains zero or one card(s), same as Waste
func (self *Reserve) Conformant() bool {
	return self.pile.Len() < 2
}

// UnsortedPairs - cards in a reserve pile are always considered to be unsorted
func (self *Reserve) UnsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Reserve) MovableTails() []*MovableTail {
	// nb same as Cell.MovableTails
	var tails []*MovableTail = []*MovableTail{}
	if self.pile.Len() > 0 {
		var card *Card = self.pile.Peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = TheBaize.FindHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &MovableTail{dst: home, tail: tail})
		}
	}
	return tails
}

func (*Reserve) Placeholder() *ebiten.Image {
	return nil
}
