package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
)

type Freecell struct{}

func (*Freecell) BuildPiles() {
	TheBaize.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil)
	TheBaize.piles = append(TheBaize.piles, TheBaize.stock)

	for x := 0; x < 4; x++ {
		c := NewCell(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, c)

	}
	for x := 4; x < 8; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		TheBaize.piles = append(TheBaize.piles, f)
		TheBaize.foundations = append(TheBaize.foundations, f)
		f.SetLabel("A")
	}

	for x := 0; x < 8; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS)
		TheBaize.piles = append(TheBaize.piles, t)
		TheBaize.tableaux = append(TheBaize.tableaux, t)
	}
}

func (*Freecell) StartGame() {
	for i := 0; i < 4; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < 7; j++ {
			MoveCard(TheBaize.stock, pile)
		}
	}
	for i := 4; i < 8; i++ {
		pile := TheBaize.tableaux[i]
		for j := 0; j < 6; j++ {
			MoveCard(TheBaize.stock, pile)
		}
	}

	if TheBaize.stock.Len() > 0 {
		println("*** still", TheBaize.stock.Len(), "cards in Stock")
	}
}

func (*Freecell) AfterMove() {
}

func (*Freecell) TailMoveError(tail []*Card) (bool, error) {
	var c1 *Card = tail[0]
	var pile Pile = c1.Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch pile.(type) {
	case *Tableau:
		var c2 *Card
		for i := 1; i < len(tail); i++ {
			c2 = tail[i]
			ok, err := CardCompare_DownAltColor(c1, c2)
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

func (*Freecell) TailAppendError(dst Pile, tail []*Card) (bool, error) {
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
		} else {
			c1 := dst.Peek()
			c2 := tail[0]
			return CardCompare_DownAltColor(c1, c2)
		}
	default:
		println("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Freecell) UnsortedPairs(pile Pile) int {
	var unsorted int
	var c1 *Card = pile.Get(0) // may be nil
	var c2 *Card
	for i := 1; i < pile.Len(); i++ {
		c2 = pile.Get(i)
		if c1.Prone() || c2.Prone() {
			unsorted++
		} else {
			ok, _ := CardCompare_DownAltColor(c1, c2)
			if !ok {
				unsorted++
			}
		}
		c1 = c2
	}
	return unsorted
}

func (*Freecell) TailTapped(tail []*Card) {
	tail[0].Owner().TailTapped(tail)
}

func (*Freecell) PileTapped(pile Pile) {
}

func (*Freecell) PercentComplete() int {
	return Script_PercentComplete()
}

func (*Freecell) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/FreeCell"
}
