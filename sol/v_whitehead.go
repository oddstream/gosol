package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Whitehead struct {
	ScriptBase
}

func (*Whitehead) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		relaxable:   false,
	}
}

func (wh *Whitehead) BuildPiles() {

	wh.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	wh.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	wh.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		wh.foundations = append(wh.foundations, f)
		f.SetLabel("A")
	}

	wh.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		wh.tableaux = append(wh.tableaux, t)
	}
}

func (wh *Whitehead) StartGame() {
	var deal = 1
	for _, pile := range wh.tableaux {
		for i := 0; i < deal; i++ {
			MoveCard(wh.stock, pile)
		}
		deal++
	}
	if s, ok := (wh.stock.subtype).(*Stock); ok {
		s.recycles = 0
	}
	wh.stock.SetRune(NORECYCLE_RUNE)
	MoveCard(wh.stock, wh.waste)
}

func (wh *Whitehead) AfterMove() {
	if wh.waste.Len() == 0 && wh.stock.Len() != 0 {
		MoveCard(wh.stock, wh.waste)
	}
}

func (*Whitehead) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile.subtype).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*Whitehead) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			return CardPair{dst.Peek(), tail[0]}.Compare_DownColor()
		}
	}
	return true, nil
}

func (*Whitehead) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownColor(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (wh *Whitehead) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile.IsStock() && len(tail) == 1 {
		MoveCard(wh.stock, wh.waste)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (wh *Whitehead) PileTapped(pile *Pile) {
	// https://politaire.com/help/whitehead
	// Only one pass through the Stock is permitted
}

func (*Whitehead) Discards() []*Pile {
	return nil
}

func (wh *Whitehead) Foundations() []*Pile {
	return wh.foundations
}

func (wh *Whitehead) Stock() *Pile {
	return wh.stock
}

func (wh *Whitehead) Waste() *Pile {
	return wh.waste
}
