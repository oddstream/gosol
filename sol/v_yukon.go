package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Yukon struct {
	ScriptBase
	extraCells int
}

func (*Yukon) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "portrait",
		wikipedia:   "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
	}
}

func (yuk *Yukon) BuildPiles() {

	yuk.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	yuk.foundations = nil
	for y := 0; y < 4; y++ {
		f := NewFoundation(image.Point{8, y})
		yuk.foundations = append(yuk.foundations, f)
		f.SetLabel("A")
	}

	yuk.cells = nil
	y := 4
	for i := 0; i < yuk.extraCells; i++ {
		c := NewCell(image.Point{8, y})
		yuk.cells = append(yuk.cells, c)
		y += 1
	}

	yuk.tableaux = nil
	for x := 0; x < 7; x++ {
		t := NewTableau(image.Point{x, 0}, FAN_DOWN, MOVE_ANY)
		yuk.tableaux = append(yuk.tableaux, t)
		t.SetLabel("K")
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

func (*Yukon) AfterMove() {}

func (*Yukon) TailMoveError([]*Card) (bool, error) {
	return true, nil
}

func (*Yukon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch (dst).category {
	case "Foundation":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuit()
		}
	case "Tableau":
		if dst.Empty() {
			return Compare_Empty(dst, tail[0])
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColor()
		}
	}
	return true, nil
}

func (*Yukon) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
}

func (*Yukon) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (*Yukon) PileTapped(*Pile) {}
