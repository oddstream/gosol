package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Yukon struct {
	ScriptBase
	twocells bool
}

func (*Yukon) PreferredWindow() string {
	return "portrait"
}

func (yuk *Yukon) BuildPiles() {
	yuk.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil)

	yuk.foundations = nil
	for y := 0; y < 4; y++ {
		f := NewFoundation(image.Point{8, y}, FAN_NONE)
		yuk.foundations = append(yuk.foundations, f)
		f.SetLabel("A")
	}

	yuk.cells = nil
	if yuk.twocells {
		for y := 4; y < 6; y++ {
			c := NewCell(image.Point{8, y})
			yuk.cells = append(yuk.cells, c)
		}
	}

	yuk.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ANY)
		yuk.tableaux = append(yuk.tableaux, t)
		if !yuk.relaxed {
			t.SetLabel("K")
		}
	}
}

func (yuk *Yukon) StartGame() {

	MoveCard(yuk.stock, yuk.tableaux[0])
	var dealDown int = 1
	for x := 1; x < 7; x++ {
		for i := 0; i < dealDown; i++ {
			MoveCard(yuk.stock, yuk.tableaux[x])
			if c := yuk.tableaux[x].Peek(); c == nil {
				break
			} else {
				c.FlipDown()
			}
		}
		dealDown++
		for i := 0; i < 5; i++ {
			MoveCard(yuk.stock, yuk.tableaux[x])
		}
	}

	if yuk.stock.Len() > 0 {
		println("*** still", yuk.stock.Len(), "cards in Stock")
	}

}

func (*Yukon) AfterMove() {
}

func (*Yukon) TailMoveError(tail []*Card) (bool, error) {
	return true, nil
}

func (*Yukon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch (dst.subtype).(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case *Tableau:
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Yukon) UnsortedPairs(pile *Pile) int {
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

func (*Yukon) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	pile.subtype.TailTapped(tail)
}

func (*Yukon) PileTapped(pile *Pile) {
}

func (*Yukon) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Yukon_(solitaire)"
}

func (*Yukon) Discards() []*Pile {
	return nil
}

func (yuk *Yukon) Foundations() []*Pile {
	return yuk.foundations
}

func (yuk *Yukon) Stock() *Pile {
	return yuk.stock
}

func (*Yukon) Waste() *Pile {
	return nil
}
