package sol

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"image"
	"log"

	"oddstream.games/gosol/cardid"
	"oddstream.games/gosol/sound"
)

// The CardID contains everything we need to serialize the card: pack, ordinal, suit and prone flag

type SavablePile struct {
	Category string          // for readability and sanity checks
	Label    string          `json:",omitempty"`
	Cards    []cardid.CardID `json:",omitempty"`
}

type SavableBaize struct {
	Piles    []*SavablePile `json:",omitempty"`
	Bookmark int            `json:",omitempty"`
	Recycles int            `json:",omitempty"`
}

func (self *Pile) Savable() *SavablePile {
	sp := &SavablePile{Category: self.category, Label: self.label}
	for _, c := range self.cards {
		sp.Cards = append(sp.Cards, c.id)
	}
	return sp
}

func (self *Pile) UpdateFromSavable(sp *SavablePile, cardPositionMap map[cardid.CardID]image.Point) {
	if self.category != sp.Category {
		log.Panicf("Baize pile (%s) and SavablePile (%s) are different", self.category, sp.Category)
	}
	self.Reset()
	// undo
	// a card that was face up (eg because it was top of pile)
	// now needs to be face down (because it is no longer top of pile)
	// card is created same face up/down as it was saved
	// when it should be created face up, then flipped.

	for _, cid := range sp.Cards {
		// pos will be set by Pile.Push(), and then lerped to
		// but we need to set a pos to lerp from
		// the default of {0,0} looks wrong
		// so we get the last known position to lerp from
		// nb a lot of the time, the pos won't change
		// ie we lerp from pos to pos, which is handled efficiently
		pos, ok := cardPositionMap[cid]
		if !ok {
			pos = cardPositionMap[cid.PackSuitOrdinal()]
			// card was put in the map face up, but is now face down
			// if still not ok, pos will be image.Point{0, 0}
		}
		var c Card = Card{id: cid, pos: pos}
		self.Push(&c) // will always flip down if pile is Stock
	}
	if len(self.cards) != len(sp.Cards) {
		log.Panicf("%s cards rebuilt incorrectly", self.category)
	}
	self.SetLabel(sp.Label)
}

func (b *Baize) NewSavableBaize() *SavableBaize {
	sb := &SavableBaize{Bookmark: b.bookmark, Recycles: b.recycles}
	for _, p := range b.piles {
		sb.Piles = append(sb.Piles, p.Savable())
	}
	return sb
}

func (b *Baize) UndoPush() {
	sb := b.NewSavableBaize()
	b.undoStack = append(b.undoStack, sb)
}

func (b *Baize) UndoPeek() *SavableBaize {
	if len(b.undoStack) == 0 {
		return nil
	}
	return b.undoStack[len(b.undoStack)-1]
}

func (b *Baize) UndoPop() (*SavableBaize, bool) {
	if len(b.undoStack) == 0 {
		return &SavableBaize{}, false
	}
	sav := b.undoStack[len(b.undoStack)-1]
	b.undoStack = b.undoStack[:len(b.undoStack)-1]
	return sav, true
}

func (b *Baize) IsSavableOk(sb *SavableBaize) bool {
	if len(b.piles) != len(sb.Piles) {
		log.Printf("Baize piles (%d) and SavableBaize piles (%d) are different", len(b.piles), len(sb.Piles))
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

func (b *Baize) IsSavableStackOk(stack []*SavableBaize) bool {
	if stack == nil {
		log.Print("No savable stack")
		return false
	}
	for i := 0; i < len(stack); i++ {
		if !b.IsSavableOk(stack[i]) {
			return false
		}
	}
	return true
}

func (b *Baize) UpdateFromSavable(sb *SavableBaize) {
	if len(b.piles) != len(sb.Piles) {
		log.Panicf("Baize piles (%d) and SavableBaize piles (%d) are different", len(b.piles), len(sb.Piles))
	}
	var cardPositionMap map[cardid.CardID]image.Point = make(map[cardid.CardID]image.Point)
	b.ForeachCard(func(c *Card) { cardPositionMap[c.id] = c.pos })

	for i := 0; i < len(sb.Piles); i++ {
		b.piles[i].UpdateFromSavable(sb.Piles[i], cardPositionMap)
	}
	sound.Play("TakeOutPackage")
	b.bookmark = sb.Bookmark
	b.recycles = sb.Recycles
	b.setFlag(dirtyCardPositions)
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.undoStack) < 2 {
		TheGame.UI.ToastError("Nothing to undo")
		return
	}
	if b.Complete() {
		TheGame.UI.ToastError("Cannot undo a completed game") // otherwise the stats can be cooked
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
	b.FindDestinations()
}

func (b *Baize) RestartDeal() {
	if b.Complete() {
		TheGame.UI.ToastError("Cannot restart a completed game") // otherwise the stats can be cooked
		return
	}
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
	b.FindDestinations()
}

// SavePosition saves the current Baize state
func (b *Baize) SavePosition() {
	if b.Complete() {
		TheGame.UI.ToastError("Cannot bookmark a completed game") // otherwise the stats can be cooked
		return
	}
	b.bookmark = len(b.undoStack)
	sb := b.UndoPeek()
	sb.Bookmark = b.bookmark
	sb.Recycles = b.recycles
	TheGame.UI.ToastInfo("Position bookmarked")
}

// LoadPosition loads a previously saved Baize state
func (b *Baize) LoadPosition() {
	if b.bookmark == 0 || b.bookmark > len(b.undoStack) {
		// println("bookmark", b.bookmark, "undostack", len(b.undoStack))
		TheGame.UI.ToastError("No bookmark")
		return
	}
	if b.Complete() {
		TheGame.UI.ToastError("Cannot undo a completed game") // otherwise the stats can be cooked
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
	b.FindDestinations()
}
