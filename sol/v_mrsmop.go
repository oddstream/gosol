package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
	"log"
)

type MrsMop struct {
	ScriptBase
	easy bool
}

func (self *MrsMop) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 2, 4, nil, 0)

	self.discards = []*Pile{}
	for x := 0; x < 4; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
		d = NewDiscard(image.Point{x + 9, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 13; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}

	self.cells = []*Pile{}
	if self.easy {
		for x := 5; x < 8; x++ {
			t := NewCell(image.Point{x, 0})
			self.cells = append(self.cells, t)
		}
	}
}

func (self *MrsMop) StartGame() {
	// 13 piles of 8 cards each
	for _, pile := range self.tableaux {
		for i := 0; i < 8; i++ {
			MoveCard(self.stock, pile)
		}
	}
	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (*MrsMop) TailMoveError(tail []*Card) (bool, error) {
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

func (*MrsMop) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
	case *Discard:
		// Discard.CanAcceptTail() has already checked
		// (1) pile is empty
		// (2) no prone cards in tail
		// (3) tail is the length of a complete set (eg 13)
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		ok, err := TailConformant(tail, CardPair.Compare_DownSuit)
		if !ok {
			return ok, err
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
		}
	}
	return true, nil
}

func (*MrsMop) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*MrsMop) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

// func (*MrsMop) PileTapped(*Pile) {}

func (self *MrsMop) Complete() bool {
	return self.SpiderComplete()
}
