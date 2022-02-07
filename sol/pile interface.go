package sol

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// "the bigger the interface, the weaker the abstraction"

/*
	3. Open, extensible object models
	An object typically describes both state and behaviour that
	acts on (or is influenced by) that state. We might account
	for both state and behaviour in our object model, but it
	would be simpler to model just one of them and then use
	it to provide the other indirectly. We choose to model (and
	expose) behaviour as a set of methods that are invoked in
	an object by name; access to state, if appropriate, is then
	provided through ‘accessor’ methods.
*/

type Pile interface {
	// implemented by Core
	Valid() bool
	Reset()
	Hidden() bool
	IsStock() bool
	IsTableau() bool
	Cards() []*Card
	MoveType() MoveType
	FanType() FanType
	SetFanType(FanType)
	Label() string
	SetLabel(string)
	Rune() rune
	SetRune(rune)
	Target() bool
	SetTarget(bool)
	Empty() bool
	Len() int
	Less(int, int) bool
	Swap(int, int)
	Get(int) *Card
	Append(*Card)
	Delete(int)
	Peek() *Card
	Pop() *Card
	Push(*Card)
	Slot() image.Point
	SetSlot(image.Point)
	SetBaizePos(image.Point)
	BaizePos() image.Point
	BaizeRect() image.Rectangle
	ScreenRect() image.Rectangle
	FannedBaizeRect() image.Rectangle
	FannedScreenRect() image.Rectangle
	PosAfter(*Card) image.Point
	Scrunch()
	Refan()
	IndexOf(*Card) int
	CanMoveTail([]*Card) (bool, error)
	MakeTail(*Card) []*Card
	ApplyToCards(func(*Card))
	BuryCards(int)

	Savable() *SavablePile
	UpdateFromSavable(*SavablePile)

	DrawStaticCards(*ebiten.Image)
	DrawTransitioningCards(*ebiten.Image)
	DrawFlippingCards(*ebiten.Image)
	DrawDraggingCards(*ebiten.Image)

	Update()
	CreateBackgroundImage()
	Draw(*ebiten.Image)

	// implemented by Cell, Discard, Foundation, Reserve, Stock, Tableau, Waste
	CanAcceptCard(*Card) (bool, error)
	CanAcceptTail([]*Card) (bool, error)
	TailTapped([]*Card)
	Collect()
	Conformant() bool
	Complete() bool
	UnsortedPairs() int
}
