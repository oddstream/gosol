package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Australian struct {
	ScriptBase
}

func (ez *Australian) BuildPiles() {
	ez.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	ez.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	ez.foundations = nil
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		ez.foundations = append(ez.foundations, f)
		f.SetLabel("A")
	}

	ez.tableaux = nil
	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		ez.tableaux = append(ez.tableaux, t)
		t.SetLabel("K")
	}
}

func (ez *Australian) StartGame() {
	if s, ok := (ez.stock.subtype).(*Stock); ok {
		s.recycles = 0
	}
	for _, pile := range ez.tableaux {
		for i := 0; i < 4; i++ {
			MoveCard(ez.stock, pile)
		}
	}
}

func (ez *Australian) AfterMove() {
}

func (*Australian) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Australian) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	}
	return true, nil
}

func (*Australian) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownSuit(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (ez *Australian) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		c := pile.Pop()
		ez.waste.Push(c)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (ez *Australian) PileTapped(pile *Pile) {
}

func (*Australian) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Australian_Patience"
}

func (*Australian) Discards() []*Pile {
	return nil
}

func (ez *Australian) Foundations() []*Pile {
	return ez.foundations
}

func (ez *Australian) Stock() *Pile {
	return ez.stock
}

func (ez *Australian) Waste() *Pile {
	return ez.waste
}
