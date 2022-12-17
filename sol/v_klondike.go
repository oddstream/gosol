package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
)

type Klondike struct {
	ScriptBase
	wikipedia      string
	draw, recycles int
	thoughtful     bool
}

func (self *Klondike) BuildPiles() {

	if self.draw == 0 {
		self.draw = 1
	}
	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.SetLabel("K")
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *Klondike) StartGame() {
	var dealDown = 0
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			card := MoveCard(self.stock, pile)
			if !self.thoughtful {
				card.FlipDown()
			}
		}
		dealDown++
		MoveCard(self.stock, pile)
	}
	TheBaize.SetRecycles(self.recycles)
	for i := 0; i < self.draw; i++ {
		MoveCard(self.stock, self.waste)
	}
}

func (self *Klondike) AfterMove() {
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		for i := 0; i < self.draw; i++ {
			MoveCard(self.stock, self.waste)
		}
	}
}

func (*Klondike) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile).category {
	case "Tableau":
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownAltColor(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*Klondike) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).category {
	case "Foundation":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case "Tableau":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Klondike) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
}

func (self *Klondike) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		for i := 0; i < self.draw; i++ {
			MoveCard(self.stock, self.waste)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Klondike) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.waste, self.stock)
	}
}

func (self *Klondike) Wikipedia() string {
	return self.wikipedia
}

func (self *Klondike) CardColors() int {
	return 2
}
