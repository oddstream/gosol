package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Waste struct {
	parent *Pile
}

func NewWaste(slot image.Point, fanType FanType) *Pile {
	waste := NewPile("Waste", slot, fanType, MOVE_ONE)
	waste.vtable = &Waste{parent: &waste}
	TheBaize.AddPile(&waste)
	return &waste
}

func (self *Waste) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Can only move a single card to Waste")
	}
	if !tail[0].owner.IsStock() {
		return false, errors.New("Waste can only accept cards from the Stock")
	}
	// nb card can be - usually is - face down
	return true, nil
}

func (self *Waste) TailTapped(tail []*Card) {
	self.parent.DefaultTailTapped(tail)
}

func (self *Waste) Conformant() bool {
	// zero or one cards would be conformant
	return self.parent.Len() < 2
}

func (self *Waste) Complete() bool {
	return self.parent.Empty()
}

func (self *Waste) UnsortedPairs() int {
	// Waste (like Stock, Reserve) is always considered unsorted
	if self.parent.Empty() {
		return 0
	}
	return self.parent.Len() - 1
}

func (self *Waste) MovableTails() []*MovableTail {
	// nb same as Reserve.MovableTails
	var tails []*MovableTail = []*MovableTail{}
	if self.parent.Len() > 0 {
		var card *Card = self.parent.Peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = TheBaize.FindHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &MovableTail{dst: home, tail: tail})
		}
	}
	return tails
}

func (self *Waste) Placeholder() *ebiten.Image {
	return nil
}
