package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type FortyThieves struct {
	ScriptBase
	founds      []int
	tabs        []int
	cardsPerTab int
}

func (*FortyThieves) PreferredWindow() string {
	return "landscape"
}

func (ft *FortyThieves) BuildPiles() {
	ft.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil)
	ft.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	ft.foundations = nil
	for _, x := range ft.founds {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		ft.foundations = append(ft.foundations, f)
		f.SetLabel("A")
	}

	ft.tableaux = nil
	for _, x := range ft.tabs {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		ft.tableaux = append(ft.tableaux, t)
	}
}

func (ft *FortyThieves) StartGame() {
	for _, pile := range ft.tableaux {
		for i := 0; i < ft.cardsPerTab; i++ {
			MoveCard(ft.stock, pile)
		}
	}
	if s, ok := (ft.stock.subtype).(*Stock); ok {
		s.recycles = 0
	}
	MoveCard(ft.stock, ft.waste)
}

func (ft *FortyThieves) AfterMove() {
	if ft.waste.Len() == 0 && ft.stock.Len() != 0 {
		MoveCard(ft.stock, ft.waste)
	}
}

func (*FortyThieves) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile.subtype).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*FortyThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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

func (*FortyThieves) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if ok, _ := pair.Compare_DownSuit(); !ok {
			unsorted++
		}
	}
	return unsorted
}

func (ft *FortyThieves) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		MoveCard(ft.stock, ft.waste)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (*FortyThieves) PileTapped(*Pile) {
}

func (*FortyThieves) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)"
}

func (ft *FortyThieves) Discards() []*Pile {
	return nil
}

func (ft *FortyThieves) Foundations() []*Pile {
	return ft.foundations
}

func (ft *FortyThieves) Stock() *Pile {
	return ft.stock
}

func (ft *FortyThieves) Waste() *Pile {
	return ft.waste
}
