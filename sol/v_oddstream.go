package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
)

type Oddstream struct {
	ScriptBase
}

func (self *Oddstream) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil, 0)

	self.cells = []*Pile{}
	for x := 1; x < 4; x++ {
		c := NewCell(image.Point{x, 0})
		self.cells = append(self.cells, c)
	}
	self.foundations = []*Pile{}
	for x := 4; x < 12; x++ {
		f := NewFoundation(image.Point{x, 0})
		f.SetLabel("A")
		self.foundations = append(self.foundations, f)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 12; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		if x%2 == 0 {
			t.SetLabel("K")
		} else {
			t.SetLabel("X")
		}
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Oddstream) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for _, tab := range self.tableaux {
		for i := 0; i < 3; i++ {
			MoveCard(self.stock, tab)
		}
	}

	TheGame.Baize.SetRecycles(0)
}

func (*Oddstream) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := TailConformant(tail, CardPair.Compare_DownSuit)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (*Oddstream) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	}
	return true, nil
}

func (*Oddstream) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Oddstream) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Stock:
		if !self.stock.Empty() {
			for _, tab := range self.tableaux {
				MoveCard(self.stock, tab)
			}
		}
	default:
		tail[0].Owner().vtable.TailTapped(tail)
	}
}

func (*Oddstream) PileTapped(*Pile) {}
