package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
	"log"

	"oddstream.games/gosol/util"
)

type CanThieves struct {
	ScriptBase
}

func (self *CanThieves) BuildPiles() {

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	if self.reserves != nil {
		log.Println("*** reserves is not nil ***")
	}
	self.reserves = nil
	self.reserves = append(self.reserves, NewReserve(image.Point{0, 1}, FAN_DOWN))

	self.foundations = nil
	for x := 3; x < 11; x++ {
		self.foundations = append(self.foundations, NewFoundation(image.Point{x, 0}))
	}

	self.tableaux = nil
	for x := 2; x < 6; x++ {
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY))
	}
	for x := 7; x < 12; x++ {
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY))
	}
}

func (self *CanThieves) StartGame() {
	for _, pile := range self.foundations {
		pile.SetLabel("")
	}

	// "At the start of the game 13 cards are dealt here."
	for i := 0; i < 12; i++ {
		MoveCard(self.stock, self.reserves[0]).FlipDown()
	}
	MoveCard(self.stock, self.reserves[0])

	// "At the start of the game 1 card is dealt each to the left 4 piles,
	// and 8 cards are dealt to each of the remaining 5 piles."
	for i := 0; i < 4; i++ {
		MoveCard(self.stock, self.tableaux[i])
	}
	for i := 4; i < 9; i++ {
		for j := 0; j < 8; j++ {
			MoveCard(self.stock, self.tableaux[i])
		}
	}

	TheBaize.SetRecycles(2)
}

func (self *CanThieves) AfterMove() {
	if self.foundations[0].label == "" {
		// The first card played to a foundation will determine the starting ordinal for all the foundations
		var ord int = 0
		for _, f := range self.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.Peek().ID.Ordinal()
				break
			}
		}
		if ord != 0 {
			for _, f := range self.foundations {
				f.SetLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

func (self *CanThieves) inFirstFour(tab *Pile) bool {
	for i := 0; i < 4; i++ {
		if tab == self.tableaux[i] {
			return true
		}
	}
	return false
}

func (self *CanThieves) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			ok, err := TailConformant(tail, CardPair.Compare_DownAltColorWrap)
			if !ok {
				return ok, err
			}
		} else {
			ok, err := TailConformant(tail, CardPair.Compare_DownSuitWrap)
			if !ok {
				return ok, err
			}
		}
	}
	return true, nil
}

func (self *CanThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			if self.inFirstFour(dst) {
				return TailConformant(tail, CardPair.Compare_DownAltColorWrap)
			} else {
				return TailConformant(tail, CardPair.Compare_DownSuitWrap)
			}
		} else {
			if self.inFirstFour(dst) {
				ok, err := TailConformant(tail, CardPair.Compare_DownAltColorWrap)
				if !ok {
					return ok, err
				}
				return CardPair{dst.Peek(), card}.Compare_DownAltColorWrap()
			} else {
				ok, err := TailConformant(tail, CardPair.Compare_DownSuitWrap)
				if !ok {
					return ok, err
				}
				return CardPair{dst.Peek(), card}.Compare_DownSuitWrap()
			}
		}
	}
	return true, nil
}

func (self *CanThieves) UnsortedPairs(pile *Pile) int {
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			return UnsortedPairs(pile, CardPair.Compare_DownAltColorWrap)
		} else {
			return UnsortedPairs(pile, CardPair.Compare_DownSuitWrap)
		}
	default:
		log.Println("*** eh?", pile.category)
	}
	return 0
}

func (self *CanThieves) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *CanThieves) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.waste, self.stock)
	}
}

func (self *CanThieves) Wikipedia() string {
	return "https://www.goodsol.com/pgshelp/index.html?demons_and_thieves.htm"
}

func (self *CanThieves) CardColors() int {
	return 2
}
