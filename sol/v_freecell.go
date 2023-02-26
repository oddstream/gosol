package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
	"log"

	"oddstream.games/gosol/cardid"
)

type Freecell struct {
	ScriptBase
	tabCompareFunc CardPairCompareFunc
	blind, easy    bool
}

func (self *Freecell) BuildPiles() {

	if self.cardColors == 0 {
		self.cardColors = 2
	}
	if self.tabCompareFunc == nil {
		self.tabCompareFunc = CardPair.Compare_DownAltColor
	}

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.cells = []*Pile{}
	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, NewCell(image.Point{x, 0}))
	}

	self.foundations = []*Pile{}
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = []*Pile{}
	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Freecell) StartGame() {
	if self.easy {
		self.foundations[0].Push(self.stock.Extract(0, 1, cardid.CLUB))
		self.foundations[1].Push(self.stock.Extract(0, 1, cardid.DIAMOND))
		self.foundations[2].Push(self.stock.Extract(0, 1, cardid.HEART))
		self.foundations[3].Push(self.stock.Extract(0, 1, cardid.SPADE))
		for _, t := range self.tableaux {
			for i := 0; i < 6; i++ {
				MoveCard(self.stock, t)
			}
		}
	} else {
		// 4 piles of 7 cards
		// 4 piles of 6 cards
		for i := 0; i < 4; i++ {
			t := self.tableaux[i]
			for j := 0; j < 7; j++ {
				MoveCard(self.stock, t)
			}
		}
		for i := 4; i < 8; i++ {
			t := self.tableaux[i]
			for j := 0; j < 6; j++ {
				MoveCard(self.stock, t)
			}
		}
	}
	if self.blind {
		for _, t := range self.tableaux {
			topCard := t.Peek()
			for _, card := range t.cards {
				if card != topCard {
					card.FlipDown()
				}
			}
		}
	}
	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (self *Freecell) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := TailConformant(tail, self.tabCompareFunc)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *Freecell) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return self.tabCompareFunc(CardPair{dst.Peek(), tail[0]})
		}
	}
	return true, nil
}

func (self *Freecell) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, self.tabCompareFunc)
}

func (*Freecell) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

// func (*Freecell) PileTapped(*Pile) {}
