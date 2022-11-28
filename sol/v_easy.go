package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Easy struct {
	ScriptBase
}

func (ez *Easy) BuildPiles() {

	ez.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	ez.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	ez.foundations = nil
	for x := 9; x < 13; x++ {
		f := NewFoundation(image.Point{x, 0})
		ez.foundations = append(ez.foundations, f)
		f.SetLabel("A")
	}

	ez.tableaux = nil
	for x := 0; x < 13; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		ez.tableaux = append(ez.tableaux, t)
		t.SetLabel("K")
	}
}

func (ez *Easy) StartGame() {

	MoveNamedCard(ez.stock, CLUB, 1, ez.foundations[0])
	MoveNamedCard(ez.stock, DIAMOND, 1, ez.foundations[1])
	MoveNamedCard(ez.stock, HEART, 1, ez.foundations[2])
	MoveNamedCard(ez.stock, SPADE, 1, ez.foundations[3])

	MoveNamedCard(ez.stock, CLUB, 2, ez.foundations[0])
	MoveNamedCard(ez.stock, DIAMOND, 2, ez.foundations[1])
	MoveNamedCard(ez.stock, HEART, 2, ez.foundations[2])
	MoveNamedCard(ez.stock, SPADE, 2, ez.foundations[3])

	MoveNamedCard(ez.stock, CLUB, 3, ez.foundations[0])
	MoveNamedCard(ez.stock, DIAMOND, 3, ez.foundations[1])
	MoveNamedCard(ez.stock, HEART, 3, ez.foundations[2])
	MoveNamedCard(ez.stock, SPADE, 3, ez.foundations[3])

	for _, pile := range ez.tableaux {
		for i := 0; i < 1; i++ {
			MoveCard(ez.stock, pile).FlipDown()
		}
		MoveCard(ez.stock, pile)
	}
	TheBaize.SetRecycles(32767)
	MoveCard(ez.stock, ez.waste)
}

func (ez *Easy) AfterMove() {
	if ez.waste.Len() == 0 && ez.stock.Len() != 0 {
		MoveCard(ez.stock, ez.waste)
	}
}

func (*Easy) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch pile.category {
	case "Tableau":
		var cpairs CardPairs = NewCardPairs(tail)
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*Easy) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch dst.category {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	}
	return true, nil
}

func (*Easy) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (ez *Easy) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == ez.stock && len(tail) == 1 {
		c := pile.Pop()
		ez.waste.Push(c)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (ez *Easy) PileTapped(pile *Pile) {
	if pile == ez.stock {
		RecycleWasteToStock(ez.waste, ez.stock)
	}
}
