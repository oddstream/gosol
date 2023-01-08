package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I damn well like, thank you

import (
	"image"
)

type Blockade struct {
	ScriptBase
}

func (self *Blockade) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil, 0)

	self.foundations = nil
	for x := 4; x < 12; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 12; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Blockade) StartGame() {
	for _, pile := range self.tableaux {
		MoveCard(self.stock, pile)
	}
	TheBaize.SetRecycles(0)
}

func (self *Blockade) AfterMove() {
	// An empty pile will be filled up immediately by a card from the stock.
	for _, pile := range self.tableaux {
		if pile.Empty() {
			MoveCard(self.stock, pile)
		}
	}
}

func (*Blockade) TailMoveError(tail []*Card) (bool, error) {
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

func (*Blockade) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Blockade) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Blockade) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock {
		for _, tab := range self.tableaux {
			MoveCard(self.stock, tab)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (*Blockade) PileTapped(*Pile) {}
