package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Cell struct {
	Base
}

func NewCell(slot image.Point, fanType FanType) *Cell {
	c := &Cell{}
	c.Ctor(c, "Cell", slot, FAN_NONE)
	return c
}

func (*Cell) CanMoveTail(tail []*Card) (bool, error) {
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	return true, nil
}

func (c *Cell) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if !c.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	return true, nil
}

func (c *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	if !c.Empty() {
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

func (c *Cell) Conformant() bool {
	return true
}

func (c *Cell) Complete() bool {
	return c.Empty()
}

func (c *Cell) UnsortedPairs() int {
	return 0
}
