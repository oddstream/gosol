package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gomps5/util"
)

type Duchess struct {
	ScriptBase
}

func (*Duchess) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Duchess_(solitaire)",
		relaxable:   false,
	}
}

func (du *Duchess) BuildPiles() {

	du.stock = NewStock(image.Point{1, 1}, FAN_NONE, 1, 4, nil)

	du.reserves = nil
	for i := 0; i < 4; i++ {
		du.reserves = append(du.reserves, NewReserve(image.Point{i * 2, 0}, FAN_RIGHT))
	}

	du.waste = NewWaste(image.Point{1, 2}, FAN_DOWN3)

	du.foundations = nil
	for x := 3; x < 7; x++ {
		du.foundations = append(du.foundations, NewFoundation(image.Point{x, 1}, FAN_NONE))
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
}

func (*Duchess) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile Pile = tail[0].Owner()
	switch (pile).(type) {
	case *Tableau:
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

func (du *Duchess) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).(type) {
	case *Foundation:
		if dst.Empty() {
			c := tail[0]
			ord := util.OrdinalToShortString(c.Ordinal())
			if dst.Label() == "" {
				if _, ok := (c.owner).(*Reserve); !ok {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
				for _, pile := range du.foundations {
					pile.SetLabel(ord)
				}
			}
			if ord != dst.Label() {
				return false, fmt.Errorf("Foundations can only accept an %s, not a %s", dst.Label(), ord)
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			var rescards int = 0
			for _, p := range du.reserves {
				rescards += p.Len()
			}
			if rescards > 0 {
				// Spaces that occur on the tableau are filled with any top card in the reserve
				c := tail[0]
				if _, ok := (c.owner).(*Reserve); !ok {
					return false, errors.New("An empty Tableau must be filled from a Reserve")
				}
			}
			return true, nil
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (*Duchess) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColorWrap)
}

func (du *Duchess) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	if pile == du.stock && len(tail) == 1 {
		MoveCard(du.stock, du.waste)
	} else {
		pile.TailTapped(tail)
	}
}

func (du *Duchess) PileTapped(pile Pile) {
	if pile == du.stock {
		RecycleWasteToStock(du.waste, du.stock)
	}
}
