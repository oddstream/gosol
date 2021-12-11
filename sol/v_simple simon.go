package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type SimpleSimon struct{}

func (*SimpleSimon) BuildPiles() {
	TheBaize.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	for x := 3; x < 7; x++ {
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

func (*SimpleSimon) StartGame() {
	// 3 piles of 8 cards each
	for i := 0; i < 3; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < 8; j++ {
			MoveCard(TheBaize.stock, pile)
		}
	}
	var deal int = 7
	for i := 3; i < 10; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < deal; j++ {
			MoveCard(TheBaize.stock, pile)
		}
		deal--
	}

	if TheBaize.stock.Len() > 0 {
		println("*** still", TheBaize.stock.Len(), "cards in Stock")
	}
}

func (*SimpleSimon) AfterMove() {
}

func (*SimpleSimon) TailMoveError(tail []*Card) (bool, error) {
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

func (*SimpleSimon) TailAppendError(dst Pile, tail []*Card) (bool, error) {
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

func (*SimpleSimon) UnsortedPairs(pile Pile) int {
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

func (*SimpleSimon) TailTapped(tail []*Card) {
}

func (*SimpleSimon) PileTapped(pile Pile) {
}

func (*SimpleSimon) PercentComplete() int {
	return Script_PercentComplete()
}

func (*SimpleSimon) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)"
}
