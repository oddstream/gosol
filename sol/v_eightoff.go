package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type EightOff struct {
	ScriptBase
}

func (*EightOff) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Eight_Off",
		relaxable:   true,
	}
}

func (eo *EightOff) BuildPiles() {

	eo.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil, 0)

	eo.cells = nil
	for x := 0; x < 8; x++ {
		eo.cells = append(eo.cells, NewCell(image.Point{x, 0}))
	}

	eo.foundations = nil
	for y := 0; y < 4; y++ {
		pile := NewFoundation(image.Point{9, y})
		eo.foundations = append(eo.foundations, pile)
		pile.SetLabel("A")
	}

	eo.tableaux = nil
	for x := 0; x < 8; x++ {
		pile := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		eo.tableaux = append(eo.tableaux, pile)
		pile.SetLabel("K")
	}
}

func (eo *EightOff) StartGame() {
	for i := 0; i < 4; i++ {
		MoveCard(eo.stock, eo.cells[i])
	}
	for _, pile := range eo.tableaux {
		for i := 0; i < 6; i++ {
			MoveCard(eo.stock, pile)
		}
	}
	if eo.stock.Len() > 0 {
		println("*** still", eo.stock.Len(), "cards in Stock")
	}
}

func (*EightOff) AfterMove() {}

func (*EightOff) TailMoveError(tail []*Card) (bool, error) {
	var pile Pile = tail[0].Owner()
	switch (pile).(type) {
	case *Tableau:
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*EightOff) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
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

func (*EightOff) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*EightOff) TailTapped(tail []*Card) {
	tail[0].Owner().TailTapped(tail)
}

func (*EightOff) PileTapped(Pile) {}
