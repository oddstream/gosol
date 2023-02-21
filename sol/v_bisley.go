package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"

	"oddstream.games/gosol/cardid"
)

type Bisley struct {
	ScriptBase
}

func (self *Bisley) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.foundations = nil

	for x := 0; x < 4; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("K")
	}

	for x := 0; x < 4; x++ {
		f := NewFoundation(image.Point{x, 1})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 13; x++ {
		t := NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ONE)
		self.tableaux = append(self.tableaux, t)
		t.SetLabel("X")
	}
}

func (self *Bisley) StartGame() {

	self.foundations[4].Push(self.stock.Extract(0, 1, cardid.CLUB))
	self.foundations[5].Push(self.stock.Extract(0, 1, cardid.DIAMOND))
	self.foundations[6].Push(self.stock.Extract(0, 1, cardid.HEART))
	self.foundations[7].Push(self.stock.Extract(0, 1, cardid.SPADE))

	// the first 4 tableaux have 3 cards
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			MoveCard(self.stock, self.tableaux[i])
		}
	}
	// the next 9 tableaux have 4 cards
	for i := 4; i < 13; i++ {
		for j := 0; j < 4; j++ {
			MoveCard(self.stock, self.tableaux[i])
		}
	}

	TheBaize.SetRecycles(0)
}

func (*Bisley) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Bisley) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			if dst.Label() == "A" {
				return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
			} else {
				return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
			}
		}
	case *Tableau:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			// return CardPair{dst.Peek(), tail[0]}.Compare_UpOrDownSuit()
			return CardPair{dst.Peek(), tail[0]}.ChainCall(CardPair.Compare_UpOrDown, CardPair.Compare_Suit)
		}
	}
	return true, nil
}

func (*Bisley) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownColor)
}

func (self *Bisley) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (*Bisley) PileTapped(*Pile) {}
