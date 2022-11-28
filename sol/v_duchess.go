package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"

	"oddstream.games/gosol/util"
)

type Duchess struct {
	ScriptBase
	wikipedia string
}

func (du *Duchess) BuildPiles() {

	du.stock = NewStock(image.Point{1, 1}, FAN_NONE, 1, 4, nil, 0)

	du.reserves = nil
	for i := 0; i < 4; i++ {
		du.reserves = append(du.reserves, NewReserve(image.Point{i * 2, 0}, FAN_RIGHT))
	}

	du.waste = NewWaste(image.Point{1, 2}, FAN_DOWN3)

	du.foundations = nil
	for x := 3; x < 7; x++ {
		du.foundations = append(du.foundations, NewFoundation(image.Point{x, 1}))
	}

	du.tableaux = nil
	for x := 3; x < 7; x++ {
		du.tableaux = append(du.tableaux, NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY))
	}
}

func (du *Duchess) StartGame() {
	TheBaize.SetRecycles(1)
	for _, pile := range du.foundations {
		pile.SetLabel("")
	}
	for _, pile := range du.reserves {
		MoveCard(du.stock, pile)
		MoveCard(du.stock, pile)
		MoveCard(du.stock, pile)
	}

	for _, pile := range du.tableaux {
		MoveCard(du.stock, pile)
	}
	TheUI.Toast("Move a Reserve card to a Foundation")
}

func (du *Duchess) AfterMove() {
	if du.foundations[0].label == "" {
		// To start the game, the player will choose among the top cards of the reserve fans which will start the first foundation pile.
		// Once he/she makes that decision and picks a card, the three other cards with the same rank,
		// whenever they become available, will start the other three foundations.
		var ord int = 0
		for _, f := range du.foundations {
			// find where the first card landed
			if len(f.cards) > 0 {
				ord = f.Peek().ID.Ordinal()
				break
			}
		}
		if ord == 0 {
			TheUI.Toast("Move a Reserve card to a Foundation")
		} else {
			for _, f := range du.foundations {
				f.SetLabel(util.OrdinalToShortString(ord))
			}
		}
	}
}

func (*Duchess) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].Owner()
	switch pile.category {
	case "Tableau":
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownAltColorWrap(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (du *Duchess) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	card := tail[0]
	switch (dst).category {
	case "Foundation":
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
	case "Tableau":
		if dst.Empty() {
			var rescards int = 0
			for _, p := range du.reserves {
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

func (du *Duchess) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == du.stock && len(tail) == 1 {
		MoveCard(du.stock, du.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (du *Duchess) PileTapped(pile *Pile) {
	if pile == du.stock {
		RecycleWasteToStock(du.waste, du.stock)
	}
}

func (du *Duchess) Wikipedia() string {
	return du.wikipedia
}

func (du *Duchess) CardColors() int {
	return 2
}
