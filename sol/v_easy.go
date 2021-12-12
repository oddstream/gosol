package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
	"log"
)

type Easy struct{}

func (*Easy) BuildPiles() {
	TheBaize.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	TheBaize.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)
	TheBaize.piles = append(TheBaize.piles, TheBaize.waste)

	for x := 9; x < 13; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, f)
		TheBaize.foundations = append(TheBaize.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 13; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		TheBaize.piles = append(TheBaize.piles, t)
		TheBaize.tableaux = append(TheBaize.tableaux, t)
	}
}

func (*Easy) StartGame() {
	for _, pile := range TheBaize.tableaux {
		for i := 0; i < 2; i++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
		MoveCard(TheBaize.stock, pile)
	}
	s, ok := (TheBaize.stock.subtype).(*Stock)
	if !ok {
		log.Fatal("cannot get Stock from it's interface")
	}
	s.recycles = 32767
	TheBaize.stock.SetRune(RECYCLE_RUNE)
}

func (*Easy) AfterMove() {
}

func (*Easy) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile.subtype).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		cpairs.Print()
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
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
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

func (*Easy) TailTapped(tail []*Card) {
	var c1 *Card = tail[0]
	var pile *Pile = c1.Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		c2 := pile.Pop()
		if c1 != c2 {
			println("Ooops")
		}
		TheBaize.waste.Push(c2)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (*Easy) PileTapped(pile *Pile) {
	if s, ok := (pile.subtype).(*Stock); ok {
		if s.recycles > 0 {
			for TheBaize.waste.Len() > 0 {
				MoveCard(TheBaize.waste, TheBaize.stock)
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
	return Script_PercentComplete()
}

func (*Easy) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Solitaire"
}
