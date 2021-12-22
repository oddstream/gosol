package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"

	"oddstream.games/gomps5/util"
)

type Agnes struct {
	ScriptBase
	sorel bool
}

func (ag *Agnes) BuildPiles() VariantInfo {

	ag.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	ag.waste = nil

	ag.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		ag.foundations = append(ag.foundations, f)
	}

	ag.reserves = nil
	if !ag.sorel {
		for x := 0; x < 7; x++ {
			r := NewReserve(image.Point{x, 1}, FAN_NONE)
			ag.reserves = append(ag.reserves, r)
		}
	}

	var taby int = 2
	if ag.sorel {
		taby = 1
	}
	ag.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, taby}, FAN_DOWN, MOVE_ANY)
		ag.tableaux = append(ag.tableaux, t)
	}

	return VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Agnes_(solitaire)",
	}
}

func (ag *Agnes) StartGame() {

	for _, pile := range ag.reserves {
		MoveCard(ag.stock, pile)
	}

	var deal = 1
	for _, pile := range ag.tableaux {
		for i := 0; i < deal; i++ {
			MoveCard(ag.stock, pile)
		}
		deal++
	}

	c := MoveCard(ag.stock, ag.foundations[0])
	ord := c.Ordinal()
	for _, pile := range ag.foundations {
		pile.SetLabel(util.OrdinalToShortString(ord))
	}
	ord -= 1
	if ord == 0 {
		ord = 13
	}
	for _, pile := range ag.tableaux {
		pile.SetLabel(util.OrdinalToShortString(ord))
	}
}

func (ag *Agnes) AfterMove() {
}

func (*Agnes) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch (pile.subtype).(type) {
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

func (*Agnes) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (*Agnes) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownAltColorWrap(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (ag *Agnes) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		for _, pile := range ag.reserves {
			MoveCard(ag.stock, pile)
		}
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (ag *Agnes) PileTapped(pile *Pile) {
}

func (*Agnes) Discards() []*Pile {
	return nil
}

func (ag *Agnes) Foundations() []*Pile {
	return ag.foundations
}

func (ag *Agnes) Stock() *Pile {
	return ag.stock
}

func (ag *Agnes) Waste() *Pile {
	return ag.waste
}
