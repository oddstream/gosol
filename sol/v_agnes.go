package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"

	"oddstream.games/gosol/util"
)

type Agnes struct {
	ScriptBase
}

func (self *Agnes) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	self.waste = nil

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
	}

	self.reserves = nil
	for x := 0; x < 7; x++ {
		r := NewReserve(image.Point{x, 1}, FAN_NONE)
		self.reserves = append(self.reserves, r)
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Agnes) StartGame() {

	for _, pile := range self.reserves {
		MoveCard(self.stock, pile)
	}

	var dealDown int = 0
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			card := MoveCard(self.stock, pile)
			card.FlipDown()
		}
		dealDown++
		MoveCard(self.stock, pile)
	}

	c := MoveCard(self.stock, self.foundations[0])
	ord := c.Ordinal()
	for _, pile := range self.foundations {
		pile.SetLabel(util.OrdinalToShortString(ord))
	}
	ord -= 1
	if ord == 0 {
		ord = 13
	}
	for _, pile := range self.tableaux {
		pile.SetLabel(util.OrdinalToShortString(ord))
	}
}

func (self *Agnes) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := TailConformant(tail, CardPair.Compare_DownAltColorWrap)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *Agnes) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (*Agnes) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColorWrap)
}

func (self *Agnes) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		for _, pile := range self.reserves {
			MoveCard(self.stock, pile)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (*Agnes) PileTapped(*Pile) {}
