package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Cell struct {
	parent *Pile
}

func NewCell(slot image.Point) *Pile {
	cell := NewPile("Cell", slot, FAN_NONE, MOVE_ONE)
	cell.vtable = &Cell{parent: &cell}
	TheBaize.AddPile(&cell)
	return &cell
}

func (self *Cell) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if !self.parent.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	return true, nil
}

func (self *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.parent.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Cell")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	return true, nil
}

func (self *Cell) TailTapped(tail []*Card) {
	self.parent.DefaultTailTapped(tail)
}

func (self *Cell) Collect() {
	self.parent.DefaultCollect()
}

func (*Cell) Conformant() bool {
	return true
}

func (self *Cell) Complete() bool {
	return self.parent.Empty()
}

func (*Cell) UnsortedPairs() int {
	return 0
}

func (self *Cell) MovableTails() []*MovableTail {
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
