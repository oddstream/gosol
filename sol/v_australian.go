package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Australian struct {
	ScriptBase
}

func (self *Australian) BuildPiles() {
	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.SetLabel("K")
	}
}

func (self *Australian) StartGame() {
	for _, pile := range self.tableaux {
		for i := 0; i < 4; i++ {
			MoveCard(self.stock, pile)
		}
	}
	MoveCard(self.stock, self.waste)
	TheGame.Baize.SetRecycles(0)
}

func (self *Australian) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		MoveCard(self.stock, self.waste)
	}
}

func (*Australian) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Australian) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Australian) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Australian) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		c := pile.Pop()
		self.waste.Push(c)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

// func (*Australian) PileTapped(*Pile) {}
