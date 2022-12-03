package sol

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"log"

	"oddstream.games/gosol/sound"
)

// The CardID contains everything we need to serialize the card: pack, ordinal, suit and prone flag

type SavablePile struct {
	Category string   // for readability and sanity checks
	Label    string   `json:",omitempty"`
	Symbol   rune     `json:",omitempty"`
	Cards    []CardID `json:",omitempty"`
}

type SavableBaize struct {
	Piles    []*SavablePile `json:",omitempty"`
	Bookmark int            `json:",omitempty"`
	Recycles int            `json:",omitempty"`
}

func (self *Pile) Savable() *SavablePile {
	sp := &SavablePile{Category: self.category, Label: self.label}
	for _, c := range self.cards {
		sp.Cards = append(sp.Cards, c.ID)
	}
	return sp
}

func (self *Pile) UpdateFromSavable(sp *SavablePile) {
	if self.category != sp.Category {
		log.Panic("Baize pile and SavablePile are different")
	}
	self.Reset()
	for _, cid := range sp.Cards {
		for i := 0; i < len(CardLibrary); i++ {
			if SameCardAndPack(cid, CardLibrary[i].ID) {
				c := &CardLibrary[i]
				self.Push(c)
				// Push() may have flipped the card, so do this afterwards ...
				if cid.Prone() {
					c.FlipDown()
				} else {
					c.FlipUp()
				}
				break
			}
		}
	}
	if len(self.cards) != len(sp.Cards) {
		log.Panic("cards rebuilt incorrectly")
	}
	self.SetLabel(sp.Label)
}

func (b *Baize) NewSavableBaize() *SavableBaize {
	ss := &SavableBaize{Bookmark: b.bookmark, Recycles: b.recycles}
	for _, p := range b.piles {
		ss.Piles = append(ss.Piles, p.Savable())
	}
	return ss
}

func (b *Baize) UndoPush() {
	ss := b.NewSavableBaize()
	b.undoStack = append(b.undoStack, ss)
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

func (b *Baize) UpdateFromSavable(sb *SavableBaize) {
	if len(b.piles) != len(sb.Piles) {
		log.Panic("Baize piles and SavableBaize piles are different")
	}
	sound.Play("OpenPackage")
	for i := 0; i < len(sb.Piles); i++ {
		b.piles[i].UpdateFromSavable(sb.Piles[i])
	}
	b.bookmark = sb.Bookmark
	b.recycles = sb.Recycles
	b.setFlag(dirtyCardPositions)
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.undoStack) < 2 {
		sound.Play("Blip")
		TheUI.Toast("Nothing to undo")
		return
	}
	if b.Complete() {
		TheUI.Toast("Cannot undo a completed game") // otherwise the stats can be cooked
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
	b.showMovableCards = false
	b.UpdateFromSavable(sav)
	b.UndoPush() // replace current state
	b.FindDestinations()
	b.UpdateToolbar()
	b.UpdateStatusbar()
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
	b.showMovableCards = false
	b.UpdateFromSavable(sav)
	b.bookmark = 0 // do this AFTER UpdateFromSavable
	b.UndoPush()   // replace current state
	b.FindDestinations()
	b.UpdateToolbar()
	b.UpdateStatusbar()
}

// SavePosition saves the current Baize state
func (b *Baize) SavePosition() {
	if b.Complete() {
		TheUI.Toast("Cannot bookmark a completed game") // otherwise the stats can be cooked
		sound.Play("Blip")
		return
	}
	b.bookmark = len(b.undoStack)
	sb := b.UndoPeek()
	sb.Bookmark = b.bookmark
	sb.Recycles = b.recycles
	TheUI.Toast("Position bookmarked")
}

// LoadPosition loads a previously saved Baize state
func (b *Baize) LoadPosition() {
	if b.bookmark == 0 || b.bookmark > len(b.undoStack) || b.Complete() {
		// println("bookmark", b.bookmark, "undostack", len(b.undoStack))
		TheUI.Toast("No bookmark")
		sound.Play("Blip")
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
	b.showMovableCards = false
	b.UpdateFromSavable(sav)
	b.UndoPush() // replace current state
	b.FindDestinations()
	b.UpdateToolbar()
	b.UpdateStatusbar()
}
