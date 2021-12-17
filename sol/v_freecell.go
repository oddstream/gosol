package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Freecell struct {
	ScriptPiles
}

func (fc *Freecell) BuildPiles() {
	fc.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil)

	for x := 0; x < 4; x++ {
		NewCell(image.Point{x, 0})
	}
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		fc.foundations = append(fc.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		fc.tableaux = append(fc.tableaux, t)
	}
}

func (fc *Freecell) StartGame() {
	for i := 0; i < 4; i++ {
		pile := fc.tableaux[i]
		for j := 0; j < 7; j++ {
			MoveCard(fc.stock, pile)
		}
	}
	for i := 4; i < 8; i++ {
		pile := fc.tableaux[i]
		for j := 0; j < 6; j++ {
			MoveCard(fc.stock, pile)
		}
	}

	if fc.stock.Len() > 0 {
		println("*** still", fc.stock.Len(), "cards in Stock")
	}
}

func (*Freecell) AfterMove() {
}

func (*Freecell) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch (pile.subtype).(type) {
	case *Tableau:
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownAltColor(); !ok {
				return false, err
			}
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*Freecell) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Foundation:
		if dst.Empty() {
			c1 := tail[0]
			if c1.Ordinal() != 1 {
				return false, errors.New("Empty Foundations can only accept an Ace")
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Freecell) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownAltColor(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (*Freecell) TailTapped(tail []*Card) {
	tail[0].Owner().subtype.TailTapped(tail)
}

func (*Freecell) PileTapped(*Pile) {
}

func (*Freecell) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/FreeCell"
}

func (*Freecell) Discards() []*Pile {
	return nil
}

func (fc *Freecell) Foundations() []*Pile {
	return fc.foundations
}

func (*Freecell) Waste() *Pile {
	return nil
}
