package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"

	"oddstream.games/gosol/util"
)

type Toad struct {
	ScriptBase
}

func (self *Toad) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.reserves = nil
	self.reserves = append(self.reserves, NewReserve(image.Point{3, 0}, FAN_RIGHT))

	self.foundations = nil
	for x := 0; x < 8; x++ {
		self.foundations = append(self.foundations, NewFoundation(image.Point{x, 1}))
	}

	self.tableaux = nil
	for x := 0; x < 8; x++ {
		// When moving tableau piles, you must either move the whole pile or only the top card.
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ONE_OR_ALL))
	}
}

func (self *Toad) StartGame() {

	TheGame.Baize.SetRecycles(1)

	for n := 0; n < 20; n++ {
		MoveCard(self.stock, self.reserves[0])
		self.reserves[0].Peek().FlipDown()
	}
	self.reserves[0].Peek().FlipUp()

	for _, pile := range self.tableaux {
		MoveCard(self.stock, pile)
	}
	// One card is dealt onto the first foundation. This rank will be used as a base for the other foundations.
	c := MoveCard(self.stock, self.foundations[0])
	for _, pile := range self.foundations {
		pile.SetLabel(util.OrdinalToShortString(c.Ordinal()))
	}
	MoveCard(self.stock, self.waste)
}

func (self *Toad) AfterMove() {
	// Empty spaces are filled automatically from the reserve.
	for _, p := range self.tableaux {
		if p.Empty() {
			MoveCard(self.reserves[0], p)
		}
	}
	if self.waste.Len() == 0 && self.stock.Len() != 0 {
		MoveCard(self.stock, self.waste)
	}

}

func (*Toad) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (self *Toad) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, card)
		} else {
			return CardPair{dst.Peek(), card}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			// Once the reserve is empty, spaces in the tableau can be filled with a card from the Deck [Stock/Waste], but NOT from another tableau pile.
			// pointless rule, since tableuax move rule is MOVE_ONE_OR_ALL
			if card.Owner() != self.waste {
				return false, errors.New("Empty tableaux must be filled with cards from the waste")
			}
		} else {
			return CardPair{dst.Peek(), card}.Compare_DownSuitWrap()
		}
	}
	return true, nil
}

func (*Toad) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuitWrap)
}

func (self *Toad) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		c := pile.Pop()
		self.waste.Push(c)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Toad) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.waste, self.stock)
	}
}
