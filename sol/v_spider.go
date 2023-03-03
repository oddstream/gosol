package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type Spider struct {
	ScriptBase
}

func (self *Spider) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, self.packs, self.suits, nil, 0)

	self.discards = nil
	for x := 2; x < 10; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Spider) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for i := 0; i < 4; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 6; j++ {
			MoveCard(self.stock, pile).FlipDown()
		}
	}
	for i := 4; i < 10; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 5; j++ {
			MoveCard(self.stock, pile).FlipDown()
		}
	}
	for _, pile := range self.tableaux {
		if c := pile.Peek(); c != nil {
			c.FlipUp()
		}
	}
	TheGame.Baize.SetRecycles(0)
}

func (*Spider) TailMoveError(tail []*Card) (bool, error) {
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

func (*Spider) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
		}
	}
	return true, nil
}

func (*Spider) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Spider) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Stock:
		var tabCards, emptyTabs int
		for _, tab := range self.tableaux {
			if tab.Len() == 0 {
				emptyTabs++
			} else {
				tabCards += tab.Len()
			}
		}
		if emptyTabs > 0 && tabCards >= len(self.tableaux) {
			TheGame.UI.ToastError("All empty tableaux must be filled before dealing a new row")
		} else {
			for _, tab := range self.tableaux {
				MoveCard(self.stock, tab)
			}
		}
	default:
		tail[0].Owner().vtable.TailTapped(tail)

	}
}

// func (*Spider) PileTapped(*Pile) {}

func (self *Spider) Complete() bool {
	return self.SpiderComplete()
}
