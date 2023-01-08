package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
	"log"
)

type EightOff struct {
	ScriptBase
}

func (self *EightOff) BuildPiles() {

	self.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.cells = nil
	for x := 0; x < 8; x++ {
		self.cells = append(self.cells, NewCell(image.Point{x, 0}))
	}

	self.foundations = nil
	for y := 0; y < 4; y++ {
		pile := NewFoundation(image.Point{9, y})
		self.foundations = append(self.foundations, pile)
		pile.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		pile := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, pile)
		pile.SetLabel("K")
	}
}

func (self *EightOff) StartGame() {
	for i := 0; i < 4; i++ {
		MoveCard(self.stock, self.cells[i])
	}
	for _, pile := range self.tableaux {
		for i := 0; i < 6; i++ {
			MoveCard(self.stock, pile)
		}
	}
	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (*EightOff) AfterMove() {}

func (*EightOff) TailMoveError(tail []*Card) (bool, error) {
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

func (*EightOff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	}
	return true, nil
}

func (*EightOff) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*EightOff) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (*EightOff) PileTapped(*Pile) {}
