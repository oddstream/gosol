package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"
)

type Klondike struct {
	foundations, tableaux []*Pile
	waste                 *Pile
	draw, recycles        int
}

func (kl *Klondike) BuildPiles() {
	if kl.draw == 0 {
		kl.draw = 1
	}
	NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil)
	NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		kl.foundations = append(kl.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		t.SetLabel("K")
		kl.tableaux = append(kl.tableaux, t)
	}
}

func (kl *Klondike) StartGame() {
	var dealDown = 0
	for _, pile := range kl.tableaux {
		for i := 0; i < dealDown; i++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
		dealDown++
		MoveCard(TheBaize.stock, pile)
	}
	// MoveCard(TheBaize.stock, kl.waste)
	if s, ok := (TheBaize.stock.subtype).(*Stock); ok {
		s.recycles = kl.recycles
	}
	TheBaize.stock.SetRune(RECYCLE_RUNE)
	MoveCard(TheBaize.stock, kl.waste)
}

func (kl *Klondike) AfterMove() {
	if kl.waste.Len() == 0 && TheBaize.stock.Len() != 0 {
		MoveCard(TheBaize.stock, kl.waste)
	}
}

func (*Klondike) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile.subtype).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownAltColor(); !ok {
				return false, err
			}
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*Klondike) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Waste:
		return false, errors.New("Waste can only accept cards from the Stock")
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
			c1 := tail[0]
			if c1.Ordinal() != 13 {
				return false, errors.New("Empty Tableaux can only accept a King")
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Klondike) UnsortedPairs(pile *Pile) int {
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

func (kl *Klondike) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		for i := 0; i < kl.draw; i++ {
			MoveCard(TheBaize.stock, kl.waste)
		}
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (kl *Klondike) PileTapped(pile *Pile) {
	if s, ok := (pile.subtype).(*Stock); ok {
		if s.recycles > 0 {
			for kl.waste.Len() > 0 {
				MoveCard(kl.waste, TheBaize.stock)
			}
			s.recycles--
			switch {
			case s.recycles == 0:
				TheBaize.stock.SetRune(NORECYCLE_RUNE)
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
	return TheBaize.PercentComplete()
}

func (*Klondike) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Solitaire"
}

func (*Klondike) Discards() []*Pile {
	return nil
}

func (kl *Klondike) Foundations() []*Pile {
	return kl.foundations
}

func (kl *Klondike) Waste() *Pile {
	return kl.waste
}
