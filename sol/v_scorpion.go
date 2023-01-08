package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
	"log"
)

type Scorpion struct {
	ScriptBase
}

func (self *Scorpion) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)

	self.discards = []*Pile{}
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.SetLabel("K")
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Scorpion) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for _, tab := range self.tableaux {
		for i := 0; i < 7; i++ {
			MoveCard(self.stock, tab)
		}
	}

	for i := 0; i < 4; i++ {
		tab := self.tableaux[i]
		for j := 0; j < 3; j++ {
			tab.cards[j].FlipDown()
		}
	}
	TheBaize.SetRecycles(0)
	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (*Scorpion) AfterMove() {}

func (*Scorpion) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Scorpion) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch dst.vtable.(type) {
	case *Discard:
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		ok, err := TailConformant(tail, CardPair.Compare_DownSuit)
		if !ok {
			return ok, err
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

func (*Scorpion) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Scorpion) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Stock:
		if !self.stock.Empty() {
			for _, tab := range self.tableaux {
				MoveCard(self.stock, tab)
			}
		}
	default:
		tail[0].Owner().vtable.TailTapped(tail)
	}
}

func (*Scorpion) PileTapped(*Pile) {}

func (self *Scorpion) Complete() bool {
	return self.SpiderComplete()
}
