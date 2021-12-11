package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
	"log"
)

type Spider struct {
	packs, suits int
}

func (s *Spider) BuildPiles() {
	TheBaize.stock = NewStock(image.Point{0, 0}, FAN_NONE, s.packs, s.suits, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	for x := 2; x < 10; x++ {
		c := NewDiscard(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, c)
		TheBaize.discards = append(TheBaize.discards, c)

	}

	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		TheBaize.piles = append(TheBaize.piles, t)
		TheBaize.tableaux = append(TheBaize.tableaux, t)
	}
}

func (*Spider) StartGame() {
	// The Tableau consists of 10 stacks with 6 cards in the first 4 stacks, with the 6th card face up,
	// and 5 cards in the remaining 6 stacks, with the 5th card face up.

	for i := 0; i < 4; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < 6; j++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
	}
	for i := 4; i < 10; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < 5; j++ {
			MoveCard(TheBaize.stock, pile)
			pile.Peek().FlipDown()
		}
	}
	for _, pile := range TheBaize.tableaux {
		pile.Peek().FlipUp()
	}
	s, ok := (TheBaize.stock).(*Stock)
	if !ok {
		log.Fatal("cannot get Stock from interface")
	}
	s.recycles = 0

	println(TheBaize.stock.Len(), "cards in stock")
}

func (*Spider) AfterMove() {
}

func (*Spider) TailMoveError(tail []*Card) (bool, error) {
	var c1 *Card = tail[0]
	var pile Pile = c1.Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch pile.(type) {
	case *Tableau:
		var c2 *Card
		for i := 1; i < len(tail); i++ {
			c2 = tail[i]
			ok, err := CardCompare_DownSuit(c1, c2)
			if !ok {
				return false, err
			}
			c1 = c2
		}
	default:
		println("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (*Spider) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch v := dst.(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Discard:
		c1 := tail[0]
		if c1.Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		for i := 1; i < len(tail); i++ {
			c2 := tail[i]
			ok, err := CardCompare_DownSuit(c1, c2)
			if !ok {
				return false, err
			}
			c1 = c2
		}
	case *Tableau:
		if v.Empty() {
		} else {
			c1 := dst.Peek()
			c2 := tail[0]
			return CardCompare_Down(c1, c2)
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Spider) UnsortedPairs(pile Pile) int {
	var unsorted int
	var c1 *Card = pile.Get(0) // may be nil
	var c2 *Card
	for i := 1; i < pile.Len(); i++ {
		c2 = pile.Get(i)
		if c1.Prone() || c2.Prone() {
			unsorted++
		} else {
			ok, _ := CardCompare_DownSuit(c1, c2)
			if !ok {
				unsorted++
			}
		}
		c1 = c2
	}
	return unsorted
}

func (*Spider) TailTapped(tail []*Card) {
	c := tail[0]
	p := c.Owner()
	switch p.(type) {
	case *Stock:
		var tabCards, emptyTabs int
		for _, tab := range TheBaize.tableaux {
			if tab.Len() == 0 {
				emptyTabs++
			} else {
				tabCards += tab.Len()
			}
		}
		if emptyTabs > 0 && tabCards >= len(TheBaize.tableaux) {
			TheUI.Toast("All empty tableaux must be filled before dealing a new row")
		} else {
			for _, tab := range TheBaize.tableaux {
				MoveCard(TheBaize.stock, tab)
			}
		}
	}
}

func (*Spider) PileTapped(pile Pile) {
}

func (*Spider) PercentComplete() int {
	return Script_PercentComplete()
}

func (*Spider) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Spider_(solitaire)"
}
