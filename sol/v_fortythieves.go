package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"

	"oddstream.games/gosol/cardid"
)

type FortyThieves struct {
	ScriptBase
	founds         []int
	tabs           []int
	proneRows      []int
	cardsPerTab    int
	recycles       int
	dealAces       bool
	moveType       MoveType
	tabCompareFunc CardPairCompareFunc
}

func (self *FortyThieves) BuildPiles() {

	if self.moveType == MOVE_NONE /* 0 */ {
		self.moveType = MOVE_ONE_PLUS
	}
	if self.cardColors == 0 {
		self.cardColors = 2
	}
	if self.tabCompareFunc == nil {
		self.tabCompareFunc = CardPair.Compare_DownSuit
	}

	self.stock = NewStock(image.Point{0, 0}, FAN_NONE, self.packs, 4, nil, 0)
	self.waste = NewWaste(image.Point{1, 0}, FAN_RIGHT3)

	self.foundations = nil
	for _, x := range self.founds {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for _, x := range self.tabs {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, self.moveType)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *FortyThieves) StartGame() {
	if self.dealAces {
		if c := self.stock.Extract(0, 1, cardid.CLUB); c != nil {
			self.foundations[0].Push(c)
		}
		if c := self.stock.Extract(0, 1, cardid.DIAMOND); c != nil {
			self.foundations[1].Push(c)
		}
		if c := self.stock.Extract(0, 1, cardid.HEART); c != nil {
			self.foundations[2].Push(c)
		}
		if c := self.stock.Extract(0, 1, cardid.SPADE); c != nil {
			self.foundations[3].Push(c)
		}
		if c := self.stock.Extract(1, 1, cardid.CLUB); c != nil {
			self.foundations[4].Push(c)
		}
		if c := self.stock.Extract(1, 1, cardid.DIAMOND); c != nil {
			self.foundations[5].Push(c)
		}
		if c := self.stock.Extract(1, 1, cardid.HEART); c != nil {
			self.foundations[6].Push(c)
		}
		if c := self.stock.Extract(1, 1, cardid.SPADE); c != nil {
			self.foundations[7].Push(c)
		}
	}
	for _, pile := range self.tableaux {
		for i := 0; i < self.cardsPerTab; i++ {
			MoveCard(self.stock, pile)
		}
	}
	for _, row := range self.proneRows {
		for _, pile := range self.tableaux {
			pile.Get(row).FlipDown()
		}
	}
	TheGame.Baize.SetRecycles(self.recycles)
	MoveCard(self.stock, self.waste)
}

func (self *FortyThieves) AfterMove() {
	if self.waste.Empty() && !self.stock.Empty() {
		MoveCard(self.stock, self.waste)
	}
}

func (self *FortyThieves) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		ok, err := TailConformant(tail, self.tabCompareFunc)
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

func (self *FortyThieves) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	switch dst.vtable.(type) {
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
			return self.tabCompareFunc(CardPair{dst.Peek(), tail[0]})
		}
	}
	return true, nil
}

func (self *FortyThieves) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, self.tabCompareFunc)
}

func (self *FortyThieves) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if pile == self.stock && len(tail) == 1 {
		MoveCard(self.stock, self.waste)
	} else {
		pile.vtable.TailTapped(tail)
	}
}

func (self *FortyThieves) PileTapped(pile *Pile) {
	if pile == self.stock {
		RecycleWasteToStock(self.waste, self.stock)
	}
}
