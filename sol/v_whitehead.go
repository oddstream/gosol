package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
)

type Whitehead struct {
	ScriptBase
}

func (self *Whitehead) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Whitehead) StartGame() {
	var deal = 1
	for _, pile := range self.tableaux {
		for i := 0; i < deal; i++ {
			MoveCard(self.stock, pile)
		}
		deal++
	}
	TheGame.Baize.SetRecycles(0)
	MoveCard(self.stock, self.waste)
}

func (self *Whitehead) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		MoveCard(self.stock, self.waste)
	}
}

func (*Whitehead) TailMoveError(tail []*Card) (bool, error) {
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

func (*Whitehead) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownColor()
		}
	}
	return true, nil
}

func (*Whitehead) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownColor)
}

func (self *Whitehead) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Whitehead) PileTapped(*Pile) {
	// https://politaire.com/help/whitehead
	// Only one pass through the Stock is permitted
}
