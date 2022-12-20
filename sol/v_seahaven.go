package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Seahaven struct {
	ScriptBase
	wikipedia string
}

func (self *Seahaven) BuildPiles() {

	self.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.cells = nil
	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, NewCell(image.Point{x, 0}))
	}

	self.foundations = nil
	for x := 6; x < 10; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
		t.SetLabel("K")
	}
}

func (self *Seahaven) StartGame() {
	for _, t := range self.tableaux {
		for i := 0; i < 5; i++ {
			MoveCard(self.stock, t)
		}
	}
	MoveCard(self.stock, self.cells[1])
	MoveCard(self.stock, self.cells[2])
	if DebugMode && self.stock.Len() > 0 {
		println("*** still", self.stock.Len(), "cards in Stock ***")
	}
	TheBaize.SetRecycles(0)
}

func (*Seahaven) AfterMove() {}

func (self *Seahaven) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		for _, pair := range NewCardPairs(tail) {
			if ok, err := CardPair.Compare_DownSuit(pair); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (self *Seahaven) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Seahaven) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*Seahaven) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (*Seahaven) PileTapped(*Pile) {}

func (self *Seahaven) Wikipedia() string {
	return self.wikipedia
}

func (self *Seahaven) CardColors() int {
	return 4
}
