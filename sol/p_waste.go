package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Waste struct {
	pile *Pile
}

func NewWaste(slot image.Point, fanType FanType) *Pile {
	pile := NewPile("Waste", slot, fanType, MOVE_ONE)
	pile.vtable = &Waste{pile: pile}
	return pile
}

func (*Waste) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Can only move a single card to Waste")
	}
	if !tail[0].Owner().IsStock() {
		return false, errors.New("Waste can only accept cards from the Stock")
	}
	// nb card can be - usually is - face down
	return true, nil
}

func (self *Waste) TailTapped(tail []*Card) {
	self.pile.DefaultTailTapped(tail)
}

// Conformant when contains zero or one card(s), same as Reserve
func (self *Waste) Conformant() bool {
	return self.pile.Len() < 2
}

// UnsortedPairs - cards in a waste pile are always considered to be unsorted
func (self *Waste) UnsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Waste) MovableTails() []*MovableTail {
	// nb same as Reserve.MovableTails
	var tails []*MovableTail = []*MovableTail{}
	if self.pile.Len() > 0 {
		var card *Card = self.pile.Peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = TheGame.Baize.FindHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &MovableTail{dst: home, tail: tail})
		}
	}
	return tails
}

func (*Waste) Placeholder() *ebiten.Image {
	return nil
}
