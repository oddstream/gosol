package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type Klondike struct {
	ScriptBase
	draw, recycles int
	thoughtful     bool
}

func (*Klondike) Info() *VariantInfo {
	return &VariantInfo{
		windowShape: "square",
		wikipedia:   "https://en.wikipedia.org/wiki/Solitaire",
		relaxable:   true,
	}
}

func (kl *Klondike) BuildPiles() {

	if kl.draw == 0 {
		kl.draw = 1
	}
	kl.stock = NewStock(image.Point{0, 0}, FAN_NONE, 1, 4, nil, 0)
	kl.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	kl.foundations = nil
	for x := 3; x < 7; x++ {
		f := NewFoundation(image.Point{x, 0}, FAN_NONE)
		kl.foundations = append(kl.foundations, f)
		f.SetLabel("A")
	}

	kl.tableaux = nil
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
			card := MoveCard(kl.stock, pile)
			if !kl.thoughtful {
				card.FlipDown()
			}
		}
		dealDown++
		MoveCard(kl.stock, pile)
	}
	TheBaize.SetRecycles(kl.recycles)
	for i := 0; i < kl.draw; i++ {
		MoveCard(kl.stock, kl.waste)
	}
}

func (kl *Klondike) AfterMove() {
	if kl.waste.Len() == 0 && kl.stock.Len() != 0 {
		for i := 0; i < kl.draw; i++ {
			MoveCard(kl.stock, kl.waste)
		}
	}
}

func (*Klondike) TailMoveError(tail []*Card) (bool, error) {
	var pile Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownAltColor(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*Klondike) TailAppendError(dst Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).(type) {
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

func (*Klondike) UnsortedPairs(pile Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
}

func (kl *Klondike) TailTapped(tail []*Card) {
	var pile Pile = tail[0].Owner()
	if pile == kl.stock && len(tail) == 1 {
		for i := 0; i < kl.draw; i++ {
			MoveCard(kl.stock, kl.waste)
		}
	} else {
		pile.TailTapped(tail)
	}
}

func (kl *Klondike) PileTapped(pile Pile) {
	if pile == kl.stock {
		RecycleWasteToStock(kl.waste, kl.stock)
	}
}
