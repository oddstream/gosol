package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"image"
)

type FortyThieves struct {
	ScriptBase
	wikipedia      string
	cardColors     int
	packs          int
	founds         []int
	tabs           []int
	proneRows      []int
	cardsPerTab    int
	recycles       int
	dealAces       bool
	moveType       MoveType
	tabCompareFunc func(CardPair) (bool, error)
}

func (ft *FortyThieves) BuildPiles() {

	if ft.packs == 0 {
		ft.packs = 2
	}
	if ft.moveType == MOVE_NONE /* 0 */ {
		ft.moveType = MOVE_ONE_PLUS
	}
	if ft.cardColors == 0 {
		ft.cardColors = 2
	}
	if ft.tabCompareFunc == nil {
		ft.tabCompareFunc = CardPair.Compare_DownSuit
	}

	ft.stock = NewStock(image.Point{0, 0}, FAN_NONE, ft.packs, 4, nil, 0)
	ft.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	ft.foundations = nil
	for _, x := range ft.founds {
		f := NewFoundation(image.Point{x, 0})
		ft.foundations = append(ft.foundations, f)
		f.SetLabel("A")
	}

	ft.tableaux = nil
	for _, x := range ft.tabs {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, ft.moveType)
		ft.tableaux = append(ft.tableaux, t)
	}
}

func (ft *FortyThieves) StartGame() {
	if ft.dealAces {
		if c := ft.stock.Extract(1, CLUB); c != nil {
			ft.foundations[0].Push(c)
		}
		if c := ft.stock.Extract(1, DIAMOND); c != nil {
			ft.foundations[1].Push(c)
		}
		if c := ft.stock.Extract(1, HEART); c != nil {
			ft.foundations[2].Push(c)
		}
		if c := ft.stock.Extract(1, SPADE); c != nil {
			ft.foundations[3].Push(c)
		}
		if c := ft.stock.Extract(1, CLUB); c != nil {
			ft.foundations[4].Push(c)
		}
		if c := ft.stock.Extract(1, DIAMOND); c != nil {
			ft.foundations[5].Push(c)
		}
		if c := ft.stock.Extract(1, HEART); c != nil {
			ft.foundations[6].Push(c)
		}
		if c := ft.stock.Extract(1, SPADE); c != nil {
			ft.foundations[7].Push(c)
		}
	}
	for _, pile := range ft.tableaux {
		for i := 0; i < ft.cardsPerTab; i++ {
			MoveCard(ft.stock, pile)
		}
	}
	for _, row := range ft.proneRows {
		for _, pile := range ft.tableaux {
			pile.Get(row).FlipDown()
		}
	}
	TheBaize.SetRecycles(ft.recycles)
	MoveCard(ft.stock, ft.waste)
}

func (ft *FortyThieves) AfterMove() {
	if ft.waste.Empty() && !ft.stock.Empty() {
		MoveCard(ft.stock, ft.waste)
	}
}

func (ft *FortyThieves) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch (pile).category {
	case "Tableau":
		var cpairs CardPairs = NewCardPairs(tail)
		for _, pair := range cpairs {
			// if ok, err := pair.Compare_DownSuit(); !ok {
			if ok, err := ft.tabCompareFunc(pair); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (ft *FortyThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
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
			// return CardPair{dst.Peek(), tail[0]}.Compare_DownSuit()
			return ft.tabCompareFunc(CardPair{dst.Peek(), tail[0]})
		}
	}
	return true, nil
}

func (ft *FortyThieves) UnsortedPairs(pile *Pile) int {
	// return UnsortedPairs(pile, CardPair.Compare_DownSuit)
	return UnsortedPairs(pile, ft.tabCompareFunc)
}

func (ft *FortyThieves) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == ft.stock && len(tail) == 1 {
		MoveCard(ft.stock, ft.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (ft *FortyThieves) PileTapped(pile *Pile) {
	if pile == ft.stock {
		RecycleWasteToStock(ft.waste, ft.stock)
	}
}

func (ft *FortyThieves) Wikipedia() string {
	return ft.wikipedia
}

func (ft *FortyThieves) CardColors() int {
	return ft.cardColors
}
