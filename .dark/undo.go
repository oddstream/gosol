package dark

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"log"

	"oddstream.games/gosol/cardid"
)

// The CardID contains everything we need to serialize the card: pack, ordinal, suit and prone flag

type savablePile struct {
	Category string          // for readability and sanity checks
	Label    string          `json:",omitempty"`
	Cards    []cardid.CardID `json:",omitempty"`
}

type savableBaize struct {
	Piles    []*savablePile `json:",omitempty"`
	Bookmark int            `json:",omitempty"`
	Recycles int            `json:",omitempty"`
}

func (self *Pile) savable() *savablePile {
	sp := &savablePile{Category: self.category, Label: self.label}
	for _, c := range self.cards {
		sp.Cards = append(sp.Cards, c.id)
	}
	return sp
}

func (self *Pile) updateFromSavable(sp *savablePile) {
	if self.category != sp.Category {
		log.Panicf("Baize pile (%s) and SavablePile (%s) are different", self.category, sp.Category)
	}
	self.reset()
	for _, cid := range sp.Cards {
		var c Card = Card{id: cid}
		self.push(&c) // will always flip down if pile is Stock
	}
	if len(self.cards) != len(sp.Cards) {
		log.Panicf("%s cards rebuilt incorrectly", self.category)
	}
	self.label = sp.Label
}

func (b *Baize) newSavableBaize() *savableBaize {
	sb := &savableBaize{Bookmark: b.bookmark, Recycles: b.recycles}
	for _, p := range b.piles {
		sb.Piles = append(sb.Piles, p.savable())
	}
	return sb
}

func (b *Baize) undoPush() {
	sb := b.newSavableBaize()
	b.undoStack = append(b.undoStack, sb)
}

func (b *Baize) undoPeek() *savableBaize {
	if len(b.undoStack) == 0 {
		return nil
	}
	return b.undoStack[len(b.undoStack)-1]
}

func (b *Baize) undoPop() (*savableBaize, bool) {
	if len(b.undoStack) == 0 {
		return &savableBaize{}, false
	}
	sav := b.undoStack[len(b.undoStack)-1]
	b.undoStack = b.undoStack[:len(b.undoStack)-1]
	return sav, true
}

func (b *Baize) isSavableOk(sb *savableBaize) bool {
	if len(b.piles) != len(sb.Piles) {
		log.Printf("Baize piles (%d) and savableBaize piles (%d) are different", len(b.piles), len(sb.Piles))
		return false
	}
	for i := 0; i < len(sb.Piles); i++ {
		if b.piles[i].category != sb.Piles[i].Category {
			log.Printf("Baize pile (%s) and SavablePile (%s) are different", b.piles[i].category, sb.Piles[i].Category)
			return false
		}
	}
	return true
}

func (b *Baize) isSavableStackOk(stack []*savableBaize) bool {
	if stack == nil {
		log.Print("No savable stack")
		return false
	}
	for i := 0; i < len(stack); i++ {
		if !b.isSavableOk(stack[i]) {
			return false
		}
	}
	return true
}

func (b *Baize) updateFromSavable(sb *savableBaize) {
	if len(b.piles) != len(sb.Piles) {
		log.Panicf("Baize piles (%d) and SavableBaize piles (%d) are different", len(b.piles), len(sb.Piles))
	}
	for i := 0; i < len(sb.Piles); i++ {
		b.piles[i].updateFromSavable(sb.Piles[i])
	}
	b.bookmark = sb.Bookmark
	b.recycles = sb.Recycles
}
