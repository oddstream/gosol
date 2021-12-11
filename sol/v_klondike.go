package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"
	"log"
)

type Klondike struct {
	draw, recycles int
}

func (k *Klondike) BuildPiles() {
	if k.draw == 0 {
		k.draw = 1
	}
	TheBaize.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	TheBaize.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)
	TheBaize.piles = append(TheBaize.piles, TheBaize.waste)

	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, f)
		TheBaize.foundations = append(TheBaize.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		TheBaize.piles = append(TheBaize.piles, t)
		TheBaize.tableaux = append(TheBaize.tableaux, t)
		t.SetLabel("K")
	}
}

func (k *Klondike) StartGame() {
	var dealDown = 0
	for _, pile := range TheBaize.tableaux {
		for i := 0; i < dealDown; i++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
		dealDown++
		MoveCard(TheBaize.stock, pile)
	}
	// MoveCard(TheBaize.stock, TheBaize.waste)
	s, ok := (TheBaize.stock).(*Stock)
	if !ok {
		log.Fatal("cannot get Stock from interface")
	}
	s.recycles = k.recycles
	s.SetRune(RECYCLE_RUNE)
}

func (*Klondike) AfterMove() {
}

func (*Klondike) TailMoveError(tail []*Card) (bool, error) {
	var pile Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch pile.(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		cpairs.Print()
		for _, pair := range cpairs {
			ok, err := CardCompare_DownAltColor(pair.c1, pair.c2)
			if !ok {
				return false, err
			}
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*Klondike) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch v := dst.(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Foundation:
		if v.Empty() {
			c1 := tail[0]
			if c1.Ordinal() != 1 {
				return false, errors.New("Empty Foundations can only accept an Ace")
			}
		} else {
			c1 := dst.Peek()
			c2 := tail[0]
			return CardCompare_UpSuit(c1, c2)
		}
	case *Tableau:
		if v.Empty() {
			c1 := tail[0]
			if c1.Ordinal() != 13 {
				return false, errors.New("Empty Tableaux can only accept a King")
			}
		} else {
			c1 := dst.Peek()
			c2 := tail[0]
			return CardCompare_DownAltColor(c1, c2)
		}
	case *Waste:
		return false, errors.New("Waste can only accept cards from the Stock")
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Klondike) UnsortedPairs(pile Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.Cards()) {
		if pair.EitherProne() {
			unsorted++
		} else {
			ok, _ := CardCompare_DownAltColor(pair.c1, pair.c2)
			if ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (k *Klondike) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	if _, ok := pile.(*Stock); ok && len(tail) == 1 {
		for i := 0; i < k.draw; i++ {
			MoveCard(TheBaize.stock, TheBaize.waste)
		}
	} else {
		pile.TailTapped(tail)
	}
}

func (*Klondike) PileTapped(pile Pile) {
	if s, ok := pile.(*Stock); ok {
		if s.recycles > 0 {
			for TheBaize.waste.Len() > 0 {
				MoveCard(TheBaize.waste, s)
			}
			s.recycles--
			switch {
			case s.recycles == 0:
				s.SetRune(NORECYCLE_RUNE)
				TheUI.Toast("No more recycles")
			case s.recycles == 1:
				TheUI.Toast(fmt.Sprintf("%d recycle remaining", s.recycles))
			case s.recycles < 10:
				TheUI.Toast(fmt.Sprintf("%d recycles remaining", s.recycles))
			}
		} else {
			TheUI.Toast("No more recycles")
		}
	}
}

func (*Klondike) PercentComplete() int {
	return Script_PercentComplete()
}

func (*Klondike) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Solitaire"
}
