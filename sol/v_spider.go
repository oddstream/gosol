package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Spider struct {
	discards, tableaux []*Pile
	packs, suits       int
}

func (sp *Spider) BuildPiles() {
	NewStock(image.Point{0, 0}, FAN_NONE, sp.packs, sp.suits, nil)

	for x := 2; x < 10; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		sp.discards = append(sp.discards, d)
	}

	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		sp.tableaux = append(sp.tableaux, t)
	}
}

func (ss *Spider) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for i := 0; i < 4; i++ {
		pile := ss.tableaux[i]
		for j := 0; j < 6; j++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
	}
	for i := 4; i < 10; i++ {
		pile := ss.tableaux[i]
		for j := 0; j < 5; j++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
	}
	for _, pile := range ss.tableaux {
		pile.Peek().FlipUp()
	}
	if s, ok := (TheBaize.stock.subtype).(*Stock); ok {
		s.recycles = 0
	}
	if DebugMode {
		println(TheBaize.stock.Len(), "cards in stock")
	}
}

func (*Spider) AfterMove() {
}

func (*Spider) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile.subtype).(type) {
	case *Tableau:
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*Spider) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Discard:
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Spider) UnsortedPairs(pile *Pile) int {
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

func (ss *Spider) TailTapped(tail []*Card) {
	pile := tail[0].Owner()
	switch (pile.subtype).(type) {
	case *Stock:
		var tabCards, emptyTabs int
		for _, tab := range ss.tableaux {
			if tab.Len() == 0 {
				emptyTabs++
			} else {
				tabCards += tab.Len()
			}
		}
		if emptyTabs > 0 && tabCards >= len(ss.tableaux) {
			TheUI.Toast("All empty tableaux must be filled before dealing a new row")
		} else {
			for _, tab := range ss.tableaux {
				MoveCard(TheBaize.stock, tab)
			}
		}
	}
}

func (*Spider) PileTapped(*Pile) {
}

func (*Spider) PercentComplete() int {
	return TheBaize.PercentComplete()
}

func (*Spider) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Spider_(solitaire)"
}

func (sp *Spider) Discards() []*Pile {
	return sp.discards
}

func (*Spider) Foundations() []*Pile {
	return nil
}

func (*Spider) Waste() *Pile {
	return nil
}

func (sp *Spider) AuxPile(subtype SubtypeAPI) []*Pile {
	switch subtype.(type) {
	case *Discard:
		return sp.discards
	case *Tableau:
		return sp.tableaux
	}
	return nil
}
