package sol

import "log"

// UndoPush pushes the current state onto the undo stack
func (b *Baize) UndoPush() {
	b.UndoStack = append(b.UndoStack, b.Saveable())
}

// UndoPop pops a state off the undo stack
func (b *Baize) UndoPop() (SaveableBaize, bool) {
	if len(b.UndoStack) > 0 {
		sav := b.UndoStack[len(b.UndoStack)-1]
		b.UndoStack = b.UndoStack[:len(b.UndoStack)-1]
		return sav, true
	}
	return SaveableBaize{}, false
}

// UndoPeekChecksum peeks the state at the top of the undo stack
func (b *Baize) UndoPeekChecksum() (uint32, bool) {
	if len(b.UndoStack) > 0 {
		sav := b.UndoStack[len(b.UndoStack)-1]
		return sav.Checksum, true
	}
	return 0, false
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.UndoStack) < 2 {
		b.ui.Toast("Nothing to undo")
		return
	}
	sav, ok := b.UndoPop() // removes current state
	if !ok {
		log.Fatal("error popping current state from undo stack")
	}

	sav, ok = b.UndoPop() // removes previous state for examination
	if !ok {
		log.Fatal("error popping second from undo stack")
	}
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}

// SavePosition saves the current Baize state
func (b *Baize) SavePosition() {
	b.SavedPosition = len(b.UndoStack)
	b.ui.Toast("Position bookmarked")
}

// LoadPosition loads a previously saved Baize state
func (b *Baize) LoadPosition() {
	if b.SavedPosition == 0 {
		b.ui.Toast("No bookmark")
		return
	}
	if b.SavedPosition > len(b.UndoStack) {
		println("error with saved position")
		return
	}
	var sav SaveableBaize
	var ok bool
	for len(b.UndoStack)+1 > b.SavedPosition {
		sav, ok = b.UndoPop()
		if !ok {
			log.Fatal("error popping from undo stack")
		}
	}
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}
