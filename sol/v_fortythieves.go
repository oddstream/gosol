package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
	"log"
)

type FortyThieves struct {
	tabs, cardsPerTab int
}

func (ft *FortyThieves) BuildPiles() {
	TheBaize.stock = NewStock(image.Point{0, 0}, FAN_NONE, 2, 4, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	TheBaize.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)
	TheBaize.piles = append(TheBaize.piles, TheBaize.waste)

	for x := ft.tabs - 8; x < ft.tabs; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, f)
		TheBaize.foundations = append(TheBaize.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < ft.tabs; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		TheBaize.piles = append(TheBaize.piles, t)
		TheBaize.tableaux = append(TheBaize.tableaux, t)
	}
}

func (ft *FortyThieves) StartGame() {
	for _, pile := range TheBaize.tableaux {
		for i := 0; i < ft.cardsPerTab; i++ {
			MoveCard(TheBaize.stock, pile)
		}
	}
	s, ok := (TheBaize.stock).(*Stock)
	if !ok {
		log.Fatal("cannot get Stock from it's interface")
	}
	s.recycles = 0
}

func (*FortyThieves) AfterMove() {
}

func (*FortyThieves) TailMoveError(tail []*Card) (bool, error) {
	var pile Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch pile.(type) {
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

func (*FortyThieves) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch v := dst.(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Waste:
		return false, errors.New("You cannot move cards to the Waste")
	case *Foundation:
		if v.Empty() {
			if tail[0].Ordinal() != 1 {
				return false, errors.New("Empty Foundations can only accept an Ace")
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if v.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*FortyThieves) UnsortedPairs(pile Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.Cards()) {
		if ok, _ := pair.Compare_DownSuit(); !ok {
			unsorted++
		}
	}
	return unsorted
}

func (*FortyThieves) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	if _, ok := pile.(*Stock); ok && len(tail) == 1 {
		MoveCard(TheBaize.stock, TheBaize.waste)
	} else {
		pile.TailTapped(tail)
	}
}

func (*FortyThieves) PileTapped(pile Pile) {
}

func (*FortyThieves) PercentComplete() int {
	return Script_PercentComplete()
}

func (*FortyThieves) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)"
}
