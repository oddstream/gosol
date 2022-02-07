package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
)

type Crimean struct {
	ScriptBase
	ukranian bool
}

func (*Crimean) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://old.reddit.com/r/solitaire/comments/s8ce8d/help_me_find_this_solitaire_online_please/",
		relaxable:   true,
	}
}

func (self *Crimean) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	if !self.ukranian {
		self.reserves = nil
		for x := 0; x < 3; x++ {
			self.reserves = append(self.reserves, NewReserve(image.Point{x, 0}, FAN_NONE))
		}
	}

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

func (self *Crimean) StartGame() {
	var dealDown = 0
	var dealUp = 7
	for _, pile := range self.tableaux {
		for i := 0; i < dealDown; i++ {
			MoveCard(self.stock, pile).FlipDown()
		}
		dealDown++
		for i := 0; i < dealUp; i++ {
			MoveCard(self.stock, pile)
		}
		dealUp--
	}
	if self.ukranian {
		MoveCard(self.stock, self.tableaux[4])
		MoveCard(self.stock, self.tableaux[5])
		MoveCard(self.stock, self.tableaux[6])
	} else {
		for _, pile := range self.reserves {
			MoveCard(self.stock, pile)
		}
	}
	if !self.stock.Empty() {
		println("ERROR: Stock length", self.stock.Len())
	}
}

func (self *Crimean) AfterMove() {}

func (*Crimean) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Crimean) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).(type) {
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

func (*Crimean) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (self *Crimean) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	pile.TailTapped(tail)
}

func (self *Crimean) PileTapped(pile Pile) {}
