package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"

	"oddstream.games/gosol/util"
)

type Agnes struct {
	ScriptBase
	wikipedia string
}

func (ag *Agnes) BuildPiles() {

	ag.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	ag.waste = nil

	ag.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0})
		ag.foundations = append(ag.foundations, f)
	}

	ag.reserves = nil
	for x := 0; x < 7; x++ {
		r := NewReserve(image.Point{x, 1}, FAN_NONE)
		ag.reserves = append(ag.reserves, r)
	}

	ag.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY)
		ag.tableaux = append(ag.tableaux, t)
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

func (*Agnes) AfterMove() {}

func (ag *Agnes) TailMoveError(tail []*Card) (bool, error) {
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

func (ag *Agnes) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch dst.category {
	case "Foundation":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuitWrap()
		}
	case "Tableau":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColorWrap()
		}
	}
	return true, nil
}

func (ag *Agnes) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColorWrap)
}

func (ag *Agnes) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == ag.stock && len(tail) == 1 {
		for _, pile := range ag.reserves {
			MoveCard(ag.stock, pile)
		}
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (ag *Agnes) PileTapped(*Pile) {}

func (ag *Agnes) Wikipedia() string {
	return ag.wikipedia
}

func (ag *Agnes) CardColors() int {
	return 2
}
