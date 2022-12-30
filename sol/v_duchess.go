package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"

	"oddstream.games/gosol/util"
)

type Duchess struct {
	ScriptBase
	wikipedia string
}

func (self *Duchess) BuildPiles() {

	self.stock = NewStock(image.Point{1, 1}, FAN_NONE, 1, 4, nil, 0)

	self.reserves = []*Pile{}
	for i := 0; i < 4; i++ {
		self.reserves = append(self.reserves, NewReserve(image.Point{i * 2, 0}, FAN_RIGHT))
	}

	self.waste = NewWaste(image.Point{1, 2}, FAN_DOWN3)

	self.foundations = []*Pile{}
	for x := 3; x < 7; x++ {
		self.foundations = append(self.foundations, NewFoundation(image.Point{x, 1}))
	}

	self.tableaux = []*Pile{}
	for x := 3; x < 7; x++ {
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY))
	}
}

func (self *Duchess) StartGame() {
	TheBaize.SetRecycles(1)
	for _, pile := range self.foundations {
		pile.SetLabel("")
	}
	for _, pile := range self.reserves {
		MoveCard(self.stock, pile)
		MoveCard(self.stock, pile)
		MoveCard(self.stock, pile)
	}

	for _, pile := range self.tableaux {
		MoveCard(self.stock, pile)
	}
	TheUI.ToastInfo("Move a Reserve card to a Foundation")
}

func (self *Duchess) AfterMove() {
	if self.foundations[0].label == "" {
		// To start the game, the player will choose among the top cards of the reserve fans which will start the first foundation pile.
		// Once he/she makes that decision and picks a card, the three other cards with the same rank,
		// whenever they become available, will start the other three foundations.
		var ord int = 0
		for _, f := range self.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.Peek().ID.Ordinal()
				break
			}
		}
		if ord == 0 {
			TheUI.ToastInfo("Move a Reserve card to a Foundation")
		} else {
			for _, f := range self.foundations {
				f.SetLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

func (*Duchess) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := TailConformant(tail, CardPair.Compare_DownAltColorWrap)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *Duchess) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			if dst.Label() == "" {
				if card.owner.category != "Reserve" {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
			}
			return Compare_Empty(dst, card)
		} else {
			return CardPair{dst.Peek(), card}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			var rescards int = 0
			for _, p := range self.reserves {
				rescards += p.Len()
			}
			if rescards > 0 {
				// Spaces that occur on the tableau are filled with any top card in the reserve
				if card.owner.category != "Reserve" {
					return false, errors.New("An empty Tableau must be filled from a Reserve")
				}
			}
			return true, nil
		} else {
			return CardPair{dst.Peek(), card}.Compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (*Duchess) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColorWrap)
}

func (self *Duchess) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *Duchess) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.waste, self.stock)
	}
}

func (self *Duchess) Wikipedia() string {
	return self.wikipedia
}

func (self *Duchess) CardColors() int {
	return 2
}
