package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"

	"oddstream.games/gomps5/util"
)

type Toad struct {
	stock, waste, reserve *Pile
	foundations, tableaux []*Pile
}

func (t *Toad) BuildPiles() {
	t.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil)
	t.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)
	t.reserve = NewReserve(image.Point{3, 0}, FAN_RIGHT)

	for x := 0; x < 8; x++ {
		t.foundations = append(t.foundations, NewFoundation(image.Point{x, 1}, FAN_NONE))
	}

	for x := 0; x < 8; x++ {
		t.tableaux = append(t.tableaux, NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ONE_OR_ALL))
	}
}

func (t *Toad) StartGame() {

	if s, ok := (t.stock.subtype).(*Stock); ok {
		s.recycles = 1
	}
	t.stock.SetRune(RECYCLE_RUNE)

	for n := 0; n < 20; n++ {
		MoveCard(t.stock, t.reserve)
		t.reserve.Peek().FlipDown()
	}
	t.reserve.Peek().FlipUp()

	for _, pile := range t.tableaux {
		MoveCard(t.stock, pile)
	}
	MoveCard(t.stock, t.foundations[0])
	c := t.foundations[0].Peek()
	for _, pile := range t.foundations {
		pile.SetLabel(util.OrdinalToShortString(c.Ordinal()))
	}
}

func (*Toad) AfterMove() {
}

func (*Toad) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Toad) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Waste:
		return false, errors.New("Waste can only accept cards from the Stock")
	case *Foundation:
		if dst.Empty() {
			if util.OrdinalToShortString(tail[0].Ordinal()) != dst.label {
				return false, fmt.Errorf("Empty Foundations can only accept an %s", dst.label)
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuitWrap()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Toad) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownSuitWrap(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (t *Toad) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		c := pile.Pop()
		t.waste.Push(c)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (t *Toad) PileTapped(pile *Pile) {
	if s, ok := (pile.subtype).(*Stock); ok {
		if s.recycles > 0 {
			for t.waste.Len() > 0 {
				MoveCard(t.waste, t.stock)
			}
			s.recycles--
			if s.recycles == 0 {
				t.stock.SetRune(NORECYCLE_RUNE)
			}
		} else {
			TheUI.Toast("No more recycles")
		}
	}
}

func (*Toad) PercentComplete() int {
	return TheBaize.PercentComplete()
}

func (*Toad) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/American_Toad_(solitaire)"
}
