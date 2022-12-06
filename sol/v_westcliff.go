package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Westcliff struct {
	ScriptBase
	wikipedia string
	variant   string
}

func (self *Westcliff) BuildPiles() {
	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	switch self.variant {
	case "Classic":
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
			self.tableaux = append(self.tableaux, t)
		}
	case "American":
		self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)
		self.foundations = nil
		for x := 6; x < 10; x++ {
			f := NewFoundation(image.Point{x, 0})
			self.foundations = append(self.foundations, f)
			f.SetLabel("A")
		}
		self.tableaux = nil
		for x := 0; x < 10; x++ {
			t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
		}
	case "Easthaven":
		self.waste = nil
		self.foundations = nil
		for x := 3; x < 7; x++ {
			f := NewFoundation(image.Point{x, 0})
			self.foundations = append(self.foundations, f)
			f.SetLabel("A")
		}
		self.tableaux = nil
		for x := 0; x < 7; x++ {
			t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
			self.tableaux = append(self.tableaux, t)
			t.SetLabel("K")
		}
	}
}

func (self *Westcliff) StartGame() {
	switch self.variant {
	case "Classic":
		MoveNamedCard(self.stock, CLUB, 1, self.foundations[0])
		MoveNamedCard(self.stock, DIAMOND, 1, self.foundations[1])
		MoveNamedCard(self.stock, HEART, 1, self.foundations[2])
		MoveNamedCard(self.stock, SPADE, 1, self.foundations[3])
		fallthrough
	case "American", "Easthaven":
		for _, pile := range self.tableaux {
			for i := 0; i < 2; i++ {
				card := MoveCard(self.stock, pile)
				card.FlipDown()
			}
		}
		for _, pile := range self.tableaux {
			MoveCard(self.stock, pile)
		}
		if self.waste != nil {
			MoveCard(self.stock, self.waste)
		}
	}
	TheBaize.SetRecycles(0)
}

func (self *Westcliff) AfterMove() {
	if self.waste != nil {
		if self.waste.Len() == 0 && self.stock.Len() != 0 {
			MoveCard(self.stock, self.waste)
		}
	}
}

func (*Westcliff) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
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

func (*Westcliff) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*Westcliff) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
}

func (self *Westcliff) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		switch self.variant {
		case "Classic", "American":
			MoveCard(self.stock, self.waste)
		case "Easthaven":
			for _, pile := range self.tableaux {
				MoveCard(self.stock, pile)
			}
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Westcliff) PileTapped(pile *Pile) {}

func (self *Westcliff) Wikipedia() string {
	return self.wikipedia
}

func (self *Westcliff) CardColors() int {
	return 2
}
