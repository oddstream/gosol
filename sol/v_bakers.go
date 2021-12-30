package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type BakersDozen struct {
	ScriptBase
}

func (*BakersDozen) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
		relaxable:   true,
	}
}

func (bd *BakersDozen) BuildPiles() {

	bd.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil)

	bd.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ONE)
		bd.tableaux = append(bd.tableaux, t)
		t.SetLabel("X")
	}
	for x := 0; x < 6; x++ {
		t := NewTableau(image.Point{x, 3}, FAN_DOWN, MOVE_ONE)
		bd.tableaux = append(bd.tableaux, t)
		t.SetLabel("X")
	}

	bd.foundations = nil
	for y := 0; y < 4; y++ {
		f := NewFoundation(image.Point{9, y}, FAN_NONE)
		bd.foundations = append(bd.foundations, f)
		f.SetLabel("A")
	}
}

func (bd *BakersDozen) StartGame() {

	for _, tab := range bd.tableaux {
		for x := 0; x < 4; x++ {
			MoveCard(bd.stock, tab)
		}
		// demote kings
		tab.BuryCards(13)
	}

	if bd.stock.Len() > 0 {
		println("*** still", bd.stock.Len(), "cards in Stock")
	}

}

func (*BakersDozen) AfterMove() {
}

func (*BakersDozen) TailMoveError(tail []*Card) (bool, error) {
	// attempt to move more than one card will be caught before this
	return true, nil
}

func (*BakersDozen) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	switch (dst).(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
			return false, errors.New("Cannot move a card to an empty Tableau")
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
		}
	}
	return true, nil
}

func (*BakersDozen) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*BakersDozen) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	pile.TailTapped(tail)
}

func (*BakersDozen) PileTapped(Pile) {}
