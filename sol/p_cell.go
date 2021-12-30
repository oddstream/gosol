package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
)

type Cell struct {
	Core
}

func NewCell(slot image.Point) *Cell {
	cell := &Cell{Core: NewCore("Stock", slot, FAN_NONE, MOVE_ONE)}
	TheBaize.AddPile(cell)
	return cell
}

func (self *Cell) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if !self.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	return true, nil
}

func (self *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.Empty() {
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
	GenericTailTapped(self, tail)
}

func (self *Cell) Collect() {
	GenericCollect(self)
}

func (*Cell) Conformant() bool {
	return true
}

func (self *Cell) Complete() bool {
	return self.Empty()
}

func (*Cell) UnsortedPairs() int {
	return 0
}
