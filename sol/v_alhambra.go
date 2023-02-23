package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"

	"oddstream.games/gosol/cardid"
)

type Alhambra struct {
	ScriptBase
}

func (self *Alhambra) BuildPiles() {

	self.stock = NewStock(image.Point{0, 3}, FAN_NONE, 2, 4, nil, 0)

	// waste pile implemented as a tableau because cards may be built on it
	self.tableaux = nil
	t := NewTableau(image.Point{1, 3}, FAN_RIGHT3, MOVE_ONE)
	self.tableaux = append(self.tableaux, t)

	self.foundations = nil
	for x := 0; x < 4; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("K")
	}

	self.reserves = nil
	for x := 0; x < 8; x++ {
		r := NewReserve(image.Point{x, 1}, FAN_DOWN)
		self.reserves = append(self.reserves, r)
	}
}

func (self *Alhambra) StartGame() {

	self.foundations[0].Push(self.stock.Extract(0, 1, cardid.CLUB))
	self.foundations[1].Push(self.stock.Extract(0, 1, cardid.DIAMOND))
	self.foundations[2].Push(self.stock.Extract(0, 1, cardid.HEART))
	self.foundations[3].Push(self.stock.Extract(0, 1, cardid.SPADE))
	self.foundations[4].Push(self.stock.Extract(0, 13, cardid.CLUB))
	self.foundations[5].Push(self.stock.Extract(0, 13, cardid.DIAMOND))
	self.foundations[6].Push(self.stock.Extract(0, 13, cardid.HEART))
	self.foundations[7].Push(self.stock.Extract(0, 13, cardid.SPADE))

	for _, r := range self.reserves {
		for i := 0; i < 4; i++ {
			MoveCard(self.stock, r)
		}
	}

	TheGame.Baize.SetRecycles(2)
}

func (*Alhambra) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Alhambra) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0]) // never happens
		} else {
			if dst.Label() == "A" {
				return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
			} else if dst.Label() == "K" {
				return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
			}
		}
	case *Tableau:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpOrDownSuitWrap()
		}
	}
	return true, nil
}

func (*Alhambra) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownColor)
}

func (self *Alhambra) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.tableaux[0])
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Alhambra) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.tableaux[0], self.stock)
	}
}
