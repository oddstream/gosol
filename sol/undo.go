package sol

import "log"

// UndoPush pushes the current state onto the undo stack
func (b *Baize) UndoPush() {
	b.UndoStack = append(b.UndoStack, b.Saveable())
	b.MarkMovable()
	b.percentComplete = b.calcPercentComplete()
	b.ui.SetMoves(len(b.UndoStack) - 1)
	b.ui.SetPercent(b.percentComplete)
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
		// sav := b.UndoStack[len(b.UndoStack)-1]
		// return sav.Checksum, true
		return b.UndoStack[len(b.UndoStack)-1].Checksum, true

	}
	return 0, false
}

// Undo reverts the Baize state to it's previous state
func (b *Baize) Undo() {
	if len(b.UndoStack) < 2 {
		b.ui.Toast("Nothing to undo")
		return
	}
	if b.State == Complete {
		b.ui.Toast("Cannot undo a completed game")
		return
	}
	sav, ok := b.UndoPop() // removes current state
	if !ok {
		log.Panic("error popping current state from undo stack")
	}

	sav, ok = b.UndoPop() // removes previous state for examination
	if !ok {
		log.Panic("error popping second state from undo stack")
	}
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}

// SavePosition saves the current Baize state
func (b *Baize) SavePosition() {
	if b.State == Complete {
		b.ui.Toast("Cannot bookmark a completed game")
		return
	}
	b.SavedPosition = len(b.UndoStack)
	b.ui.Toast("Position bookmarked")
}

// LoadPosition loads a previously saved Baize state
func (b *Baize) LoadPosition() {
	if b.SavedPosition == 0 || b.SavedPosition > len(b.UndoStack) {
		b.ui.Toast("No bookmark")
		return
	}
	if b.State == Complete {
		b.ui.Toast("Cannot undo a completed game")
		return
	}
	var sav SaveableBaize
	var ok bool
	for len(b.UndoStack)+1 > b.SavedPosition {
		sav, ok = b.UndoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}

// RestartGame loads a previously saved Baize state
func (b *Baize) RestartGame() {
	var sav SaveableBaize
	var ok bool
	for len(b.UndoStack) > 0 {
		sav, ok = b.UndoPop()
		if !ok {
			log.Panic("error popping from undo stack")
		}
	}
	b.SavedPosition = 0
	b.UpdateFromSaveable(sav)
	b.UndoPush() // replace current state
}
