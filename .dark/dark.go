package dark

import (
	"errors"

	"oddstream.games/gosol/sol"
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
	GetBaize() Baize
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

	PileTapped(int) (bool, error)            // pileIndex
	TailTapped(int, int) (bool, error)       // pileIndex, cardIndex
	TailDragged(int, int, int) (bool, error) // srcPileIndex, srcCardIndex, dstPileIndex
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
type Baize struct {
	Piles    []Pile // needed by LIGHT to display piles and cards
	Recycles int    // needed by LIGHT to determine Stock rune
	Bookmark int    // needed by LIGHT to grey out goto bookmark menu item
}

// Pile holds the state of the piles and cards therein.
// Pile is exported from this package because it's used to pass between light and dark.
type Pile struct {
	Category string // needed by LIGHT when creating Pile Placeholder (switch)
	Label    string // needed by LIGHT when creating Pile Placeholder
	moveType MoveType
	Cards    []Card
	vtable   PileVtabler
}

// Card holds the state of the cards.
// Card is exported from this package because it's used to pass between light and dark.
type Card struct {
	Card   sol.CardID
	Weight int
}

func NewDark() Darker {
	return &dark{}
}

func (d *dark) SetPowerMoves(value bool) {
	d.powerMoves = value
}

func (d *dark) Bookmark() {
	d.baize.Bookmark = 32
}

func (d *dark) Complete() bool {
	return false
}

func (d *dark) Conformant() bool {
	return false
}

func (d *dark) GetBaize() Baize {
	return d.baize // here: have a copy of the dark baize object
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

func (d *dark) PileTapped(pileIndex int) (bool, error) {
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

func (d *dark) TailDragged(srcPileIndex, cardIndex, dstPileIndex int) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) TailTapped(pileIndex, cardIndex int) (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) Undo() (bool, error) {
	return false, errors.New("not implemented")
}

func (d *dark) UndoStackSize() int {
	return 0
}
