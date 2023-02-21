package dark

import "errors"

// Baize holds the state of the baize, piles and cards therein.
// Baize is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Baize object as immutable, hence the unexported fields and getters.
type Baize struct {
	variant    string
	script     scripter
	pack       []Card // needed for undo/Pile.UpdateFromSavable
	powerMoves bool
	// undoStack
	// statistics (for all variants)
	piles    []*Pile // needed by LIGHT to display piles and cards
	recycles int     // needed by LIGHT to determine Stock rune
	bookmark int     // needed by LIGHT to grey out goto bookmark menu item
}

func (d *dark) NewBaize(variant string) (*Baize, error) {
	return nil, errors.New("not implemented")
}

func (d *dark) LoadBaize(variant string) (*Baize, error) {
	return nil, errors.New("not implemented")
}

// Baize public interface ////////////////////////////////////////////////////////////

func (b *Baize) Bookmark() int {
	return b.bookmark
}

func (b *Baize) Complete() bool {
	return false
}

func (b *Baize) Conformant() bool {
	return false
}

func (b *Baize) GotoBookmark() (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) PercentComplete() int {
	return 0
}

func (b *Baize) Piles() []*Pile {
	return b.piles
}

func (b *Baize) PileTapped(pile *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) Recycles() int {
	return b.recycles
}

func (b *Baize) RestartGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) SaveGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) SetPowerMoves(value bool) {
	b.powerMoves = value
}

func (b *Baize) Statistics() []string {
	return []string{}
}

func (b *Baize) TailDragged(src *Pile, card *Card, dst *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) TailTapped(pile *Pile, card *Card) (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) Undo() (bool, error) {
	return false, errors.New("not implemented")
}

func (b *Baize) UndoStackSize() int {
	return 0
}
