package dark

import (
	"errors"
)

// "private struct that implements a public interface"
// Conceptually, a value of an interface type, or interface value, has 2 components:
// a concrete type (type descriptor) and a value of that type.
// The descriptor is a pointer to virtual table and the interface value is the pointer
// to the instance of the concrete type that implements the interface.
//
// interface struct {
//		*vtable of functions, listed in the Darker interface declaration
//		*dark struct, as created by NewDark()
//	}

type Darker interface {
	ListVariantGroups() []string
	ListVariants(string) []string
	NewGame(string) (bool, error)
	LoadGame() (bool, error)
	SaveGame() (bool, error)
	Baize() *Baize
	RestartGame() (bool, error)
	Undo() (bool, error)
	UndoStackSize() int
	Bookmark()
	GotoBookmark() (bool, error)
	PercentComplete() int
	Complete() bool
	Conformant() bool
	Statistics() []string
	SetPowerMoves(bool)

	PileTapped(*Pile) (bool, error)
	TailTapped(*Pile, *Card) (bool, error)
	TailDragged(*Pile, *Card, *Pile) (bool, error)
}

// dark holds the state for the current game/baize in play. It is NOT exported
// from this package, making it opaque to the client.
// All access to this struct is through the Darker interface
type dark struct {
	baize                Baize
	variant              string
	script               scripter
	numberOfCardsInStock int
	powerMoves           bool
	// undoStack
	// statistics (for all variants)
}

// Baize holds the state of the baize, piles and cards therein.
// Baize is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Baize object as immutable, hence the unexported fields and getters.
type Baize struct {
	piles    []*Pile // needed by LIGHT to display piles and cards
	recycles int     // needed by LIGHT to determine Stock rune
	bookmark int     // needed by LIGHT to grey out goto bookmark menu item
}

func (b *Baize) Piles() []*Pile {
	return b.piles
}

func (b *Baize) Recycles() int {
	return b.recycles
}

func (b *Baize) Bookmark() int {
	return b.bookmark
}

func NewDark() Darker {
	return &dark{}
}

func (d *dark) SetPowerMoves(value bool) {
	d.powerMoves = value
}

func (d *dark) Bookmark() {
	d.baize.bookmark = 32
}

func (d *dark) Complete() bool {
	return false
}

func (d *dark) Conformant() bool {
	return false
}

func (d *dark) Baize() *Baize {
	return &(d.baize)
}

func (d *dark) GotoBookmark() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) LoadGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) NewGame(variant string) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) PercentComplete() int {
	return 0
}

func (d *dark) PileTapped(pile *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) RestartGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) SaveGame() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) Statistics() []string {
	return []string{}
}

func (d *dark) TailDragged(src *Pile, card *Card, dst *Pile) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) TailTapped(pile *Pile, card *Card) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) Undo() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) UndoStackSize() int {
	return 0
}
