package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

/*
	tried declaring Cell like this:
		type Cell struct {
			*Pile
		}
	and turning c.pile.Empty() into c.Empty()
	which compiled
	but panic when Empty() was called, because receiver was nil

	also tried
		type Cell struct {
			Pile
		}
	but that gave odd results when calling Empty() (pile.cards was empty when it wasn't)

	so for the moment, stuck with the subtype (Cell) containing a pointer to the outer type (Pile)
	and the Pile containing an interface (SubtypeAPI) to the subtype

	this smells all wrong

	just want to be able to have a list of piles with identical API
	that includes operations common to all piles types (like Empty)
	and operations that have functionality specific to the subtype (like CanMoveTail)
	tried having a large Pile interface
	but that got messy when accessing Pile's members; everything had to be done through functions

	if we are supposed to remove the base type (Pile) and just have a list of types
	that satisfy a Pile interface, you end up with a lot of duplicated code
	or calls to 'GenericPush' and 'GenericLen'

	all subtypes are now the same

	https://www.toptal.com/go/golang-oop-tutorial

	type PILER interface {
		Pop()
		Push()
		...
	}

	type Pile struct {
		cards []*Card
		...
	}

	type Cell Pile

	Pile will not satisfy PILER interface, because it's missing the 'subtype' functions
	Cell will satisfy PILER interface
	can cast: var myCell Cell = Cell(myPile)
	becase Cell and Pile have identical underlying types

	*do* now have a 'clean idiomatic go' design for this, with Pile as an interface
	and a 'struct Core', but it still needs the Pile interface injected into the Core

*/

type Cell struct {
	pile *Pile
}

func NewCell(slot image.Point) *Pile {
	p := &Pile{}
	p.Ctor(&Cell{pile: p}, "Cell", slot, FAN_NONE, MOVE_ONE)
	return p
}

func (c *Cell) CanAcceptCard(card *Card) (bool, error) {
	if card.Prone() {
		return false, errors.New("Cannot add a face down card")
	}
	if !c.pile.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	return true, nil
}

func (c *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	if !c.pile.Empty() {
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

func (c *Cell) TailTapped(tail []*Card) {
	c.pile.GenericTailTapped(tail)
}

func (c *Cell) Collect() {
	c.pile.GenericCollect()
}

func (c *Cell) Conformant() bool {
	return true
}

func (c *Cell) Complete() bool {
	return c.pile.Empty()
}

func (c *Cell) UnsortedPairs() int {
	return 0
}
