package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type Spiderette struct {
	ScriptBase
}

func (self *Spiderette) BuildPiles() {

	if self.cardColors == 0 {
		self.cardColors = 4
	}

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, self.packs, self.suits, nil, 0)

	self.discards = []*Pile{}
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Spiderette) StartGame() {
	var deal int = 1
	for _, pile := range self.tableaux {
		for i := 0; i < deal; i++ {
			if c := MoveCard(self.stock, pile); c != nil {
				c.FlipDown()
			}
		}
		deal++
		MoveCard(self.stock, pile)
	}
	for _, pile := range self.tableaux {
		if c := pile.Peek(); c != nil {
			c.FlipUp()
		}
	}
	TheGame.Baize.SetRecycles(0)
}

func (*Spiderette) TailMoveError(tail []*Card) (bool, error) {
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

func (*Spiderette) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Spiderette) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Spiderette) TailTapped(tail []*Card) {
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

// func (*Spiderette) PileTapped(*Pile) {}

func (self *Spiderette) Complete() bool {
	return self.SpiderComplete()
}
