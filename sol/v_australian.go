package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Australian struct {
	ScriptBase
}

func (*Australian) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Australian_Patience",
		relaxable:   false,
	}
}

func (aus *Australian) BuildPiles() {
	aus.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	aus.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	aus.foundations = nil
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0})
		aus.foundations = append(aus.foundations, f)
		f.SetLabel("A")
	}

	aus.tableaux = nil
	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		aus.tableaux = append(aus.tableaux, t)
		t.SetLabel("K")
	}
}

func (aus *Australian) StartGame() {
	TheBaize.SetRecycles(0)
	for _, pile := range aus.tableaux {
		for i := 0; i < 4; i++ {
			MoveCard(aus.stock, pile)
		}
	}
	MoveCard(aus.stock, aus.waste)
}

func (aus *Australian) AfterMove() {
	if aus.waste.Len() == 0 && aus.stock.Len() != 0 {
		MoveCard(aus.stock, aus.waste)
	}
}

func (*Australian) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Australian) TailAppendError(dst Pile, tail []*Card) (bool, error) {
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

func (*Australian) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (aus *Australian) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	if pile == aus.stock && len(tail) == 1 {
		c := pile.Pop()
		aus.waste.Push(c)
	} else {
		pile.TailTapped(tail)
	}
}

func (*Australian) PileTapped(Pile) {}
