package sol

import (
	"fmt"
	"log"

	"oddstream.games/gomps5/sound"
)

// The CardID contains everything we need to serialize the card: pack, ordinal, suit and prone flag

type SavablePile struct {
	Category string   // for readability and sanity checks
	Label    string   `json:",omitempty"`
	Symbol   rune     `json:",omitempty"`
	Cards    []CardID `json:",omitempty"`
}

type SavableBaize struct {
	Piles    []*SavablePile
	Bookmark int
	Recycles int
}

func (base *Base) Savable() *SavablePile {
	sp := &SavablePile{Category: base.category, Label: base.label, Symbol: base.symbol}
	for _, c := range base.cards {
		sp.Cards = append(sp.Cards, c.ID)
	}
	return sp
}

func (base *Base) UpdateFromSavable(sp *SavablePile) {
	if base.category != sp.Category {
		log.Panic("Baize pile and SavablePile are different")
	}
	base.Reset()
	for _, cid := range sp.Cards {
		for i := 0; i < len(TheBaize.cardLibrary); i++ {
			if SameCardAndPack(cid, TheBaize.cardLibrary[i].ID) {
				c := &TheBaize.cardLibrary[i]
				// c.SetProne(cid.Prone())
				base.Push(c)
				if cid.Prone() {
					c.FlipDown()
				} else {
					c.FlipUp()
				}
				break
			}
		}
	}
	if len(base.cards) != len(sp.Cards) {
		log.Panic("cards rebuilt incorrectly")
	}
	base.label = sp.Label
	base.symbol = sp.Symbol
}

func (b *Baize) NewSavableBaize() *SavableBaize {
	ss := &SavableBaize{}
	for _, p := range b.piles {
		ss.Piles = append(ss.Piles, p.Savable())
	}
	if stockobject, ok := b.stock.(*Stock); ok {
		ss.Recycles = stockobject.recycles
	}
	ss.Bookmark = b.bookmark
	return ss
}

func (b *Baize) UndoPush() {
	ss := b.NewSavableBaize()
	b.undoStack = append(b.undoStack, ss)
	b.UpdateStatusbar()
}

func (b *Baize) UndoPeek() *SavableBaize {
	if len(b.undoStack) == 0 {
		return nil
	}
	return b.undoStack[len(b.undoStack)-1]
}

func (b *Baize) UndoPop() (*SavableBaize, bool) {
	if len(b.undoStack) > 0 {
		sav := b.undoStack[len(b.undoStack)-1]
		b.undoStack = b.undoStack[:len(b.undoStack)-1]
		return sav, true
	}
	return &SavableBaize{}, false
}

func (b *Baize) UpdateFromSavable(ss *SavableBaize) {
	if len(b.piles) != len(ss.Piles) {
		log.Panic("Baize piles and SavableBaize piles are different")
	}
	sound.Play("OpenPackage")
	for i := 0; i < len(ss.Piles); i++ {
		b.piles[i].UpdateFromSavable(ss.Piles[i])
		b.piles[i].Scrunch()
	}
	if stockobject, ok := b.stock.(*Stock); ok {
		stockobject.recycles = ss.Recycles
	}
	b.bookmark = ss.Bookmark
	b.setFlag(dirtyCardPositions)
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.undoStack) < 2 {
		TheUI.Toast("Nothing to undo")
		return
	}
	_, ok := b.UndoPop() // removes current state
	if !ok {
		log.Panic("error popping current state from undo stack")
	}

	sav, ok := b.UndoPop() // removes previous state for examination
	if !ok {
		log.Panic("error popping second state from undo stack")
	}
	b.UpdateFromSavable(sav)
	b.UndoPush() // replace current state
}

func (b *Baize) RestartDeal() {
	var sav *SavableBaize
	var ok bool
	for len(b.undoStack) > 0 {
		sav, ok = b.UndoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.UpdateFromSavable(sav)
	b.bookmark = 0 // do this AFTER UpdateFromSavable
	b.UndoPush()   // replace current state
}

// SavePosition saves the current Baize state
func (b *Baize) SavePosition() {
	b.bookmark = len(b.undoStack)
	sb := b.UndoPeek()
	sb.Bookmark = b.bookmark
	stockobject, ok := (b.stock).(*Stock)
	if ok {
		sb.Recycles = stockobject.recycles
	}
	TheUI.Toast(fmt.Sprintf("Position %d bookmarked", b.bookmark))
}

// LoadPosition loads a previously saved Baize state
func (b *Baize) LoadPosition() {
	if b.bookmark == 0 || b.bookmark > len(b.undoStack) {
		println("bookmark", b.bookmark, "undostack", len(b.undoStack))
		TheUI.Toast("No bookmark")
		return
	}
	var sav *SavableBaize
	var ok bool
	for len(b.undoStack)+1 > b.bookmark {
		sav, ok = b.UndoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.UpdateFromSavable(sav)
	b.UndoPush() // replace current state
}
