package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Easy struct {
	foundations, tableaux []*Pile
	waste                 *Pile
}

func (ez *Easy) BuildPiles() {
	NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	ez.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	for x := 9; x < 13; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		ez.foundations = append(ez.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 13; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		ez.tableaux = append(ez.tableaux, t)
	}
}

func (ez *Easy) StartGame() {
	MoveNamedCard(CLUB, 1, ez.foundations[0])
	MoveNamedCard(DIAMOND, 1, ez.foundations[1])
	MoveNamedCard(HEART, 1, ez.foundations[2])
	MoveNamedCard(SPADE, 1, ez.foundations[3])
	for _, pile := range ez.tableaux {
		for i := 0; i < 2; i++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
		MoveCard(TheBaize.stock, pile)
	}
	if s, ok := (TheBaize.stock.subtype).(*Stock); ok {
		s.recycles = 32767
	}
	TheBaize.stock.SetRune(RECYCLE_RUNE)
	MoveCard(TheBaize.stock, ez.waste)
}

func (ez *Easy) AfterMove() {
	if ez.waste.Len() == 0 && TheBaize.stock.Len() != 0 {
		MoveCard(TheBaize.stock, ez.waste)
	}
}

func (*Easy) TailMoveError(tail []*Card) (bool, error) {
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

func (*Easy) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Waste:
		return false, errors.New("Waste can only accept cards from the Stock")
	case *Foundation:
		if dst.Empty() {
			if tail[0].Ordinal() != 1 {
				return false, errors.New("Empty Foundations can only accept an Ace")
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Easy) UnsortedPairs(pile *Pile) int {
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

func (ez *Easy) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		c := pile.Pop()
		ez.waste.Push(c)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (ez *Easy) PileTapped(pile *Pile) {
	if s, ok := (pile.subtype).(*Stock); ok {
		if s.recycles > 0 {
			for ez.waste.Len() > 0 {
				MoveCard(ez.waste, TheBaize.stock)
			}
			s.recycles--
			if s.recycles == 0 {
				TheBaize.stock.SetRune(NORECYCLE_RUNE)
			}
		} else {
			TheUI.Toast("No more recycles")
		}
	}
}

func (*Easy) PercentComplete() int {
	return TheBaize.PercentComplete()
}

func (*Easy) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Solitaire"
}

func (*Easy) Discards() []*Pile {
	return nil
}

func (ez *Easy) Foundations() []*Pile {
	return ez.foundations
}

func (ez *Easy) Waste() *Pile {
	return ez.waste
}
