package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
	"log"
)

type Yukon struct {
	ScriptBase
	extraCells int
}

func (self *Yukon) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.foundations = nil
	for y := 0; y < 4; y++ {
		f := NewFoundation(image.Point{8, y})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.cells = nil
	y := 4
	for i := 0; i < self.extraCells; i++ {
		c := NewCell(image.Point{8, y})
		self.cells = append(self.cells, c)
		y += 1
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
		t.SetLabel("K")
	}
}

func (self *Yukon) StartGame() {

	MoveCard(self.stock, self.tableaux[0])
	var dealDown int = 1
	for x := 1; x < 7; x++ {
		for i := 0; i < dealDown; i++ {
			MoveCard(self.stock, self.tableaux[x])
			if c := self.tableaux[x].Peek(); c == nil {
				break
			} else {
				c.FlipDown()
			}
		}
		dealDown++
		for i := 0; i < 5; i++ {
			MoveCard(self.stock, self.tableaux[x])
		}
	}
	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (*Yukon) TailMoveError([]*Card) (bool, error) {
	return true, nil
}

func (*Yukon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Yukon) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
}

func (*Yukon) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (*Yukon) PileTapped(*Pile) {}
