package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type SimpleSimon struct {
	ScriptBase
}

func (*SimpleSimon) PreferredWindow() string {
	return "landscape"
}

func (ss *SimpleSimon) BuildPiles() {
	ss.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil)

	ss.discards = nil
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		ss.discards = append(ss.discards, d)
	}

	ss.tableaux = nil
	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		ss.tableaux = append(ss.tableaux, t)
	}
}

func (ss *SimpleSimon) StartGame() {
	// 3 piles of 8 cards each
	for i := 0; i < 3; i++ {
		pile := ss.tableaux[i]
		for j := 0; j < 8; j++ {
			MoveCard(ss.stock, pile)
		}
	}
	var deal int = 7
	for i := 3; i < 10; i++ {
		pile := ss.tableaux[i]
		for j := 0; j < deal; j++ {
			MoveCard(ss.stock, pile)
		}
		deal--
	}

	if ss.stock.Len() > 0 {
		println("*** still", ss.stock.Len(), "cards in Stock")
	}
}

func (*SimpleSimon) AfterMove() {
}

func (*SimpleSimon) TailMoveError(tail []*Card) (bool, error) {
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

func (*SimpleSimon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
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
	}
	return true, nil
}

func (*SimpleSimon) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if ok, _ := pair.Compare_DownSuit(); !ok {
			unsorted++
		}
	}
	return unsorted
}

func (*SimpleSimon) TailTapped(tail []*Card) {
}

func (*SimpleSimon) PileTapped(*Pile) {
}

func (*SimpleSimon) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)"
}

func (ss *SimpleSimon) Discards() []*Pile {
	return ss.discards
}

func (*SimpleSimon) Foundations() []*Pile {
	return nil
}

func (ss *SimpleSimon) Stock() *Pile {
	return ss.stock
}

func (*SimpleSimon) Waste() *Pile {
	return nil
}
