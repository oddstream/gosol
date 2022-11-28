package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Scorpion struct {
	ScriptBase
	wikipedia string
}

func (sp *Scorpion) BuildPiles() {

	sp.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)

	sp.discards = nil
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		sp.discards = append(sp.discards, d)
	}

	sp.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.SetLabel("K")
		sp.tableaux = append(sp.tableaux, t)
	}
}

func (sp *Scorpion) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for _, tab := range sp.tableaux {
		for i := 0; i < 7; i++ {
			MoveCard(sp.stock, tab)
		}
	}

	for i := 0; i < 4; i++ {
		tab := sp.tableaux[i]
		for j := 0; j < 3; j++ {
			tab.cards[j].FlipDown()
		}
	}
	TheBaize.SetRecycles(0)
	if DebugMode {
		println(sp.stock.Len(), "cards in stock")
	}
}

func (*Scorpion) AfterMove() {
}

func (*Scorpion) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Scorpion) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).category {
	case "Discard":
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
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

func (*Scorpion) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (sp *Scorpion) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	switch (pile).category {
	case "Stock":
		if !sp.stock.Empty() {
			for _, tab := range sp.tableaux {
				MoveCard(sp.stock, tab)
			}
		}
	default:
		tail[0].Owner().vtable.TailTapped(tail)
	}
}

func (*Scorpion) PileTapped(*Pile) {}

func (sp *Scorpion) Wikipedia() string {
	return sp.wikipedia
}

func (*Scorpion) CardColors() int {
	return 4
}
