package sol

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PileInterface interface { // TODO rename to "Pile"
	// implemented by Core
	Valid() bool
	Reset()
	Hidden() bool
	IsStock() bool
	IsTableau() bool
	Label() string
	SetLabel(string)
	Rune() rune
	SetRune(rune)
	Empty() bool
	Len() int
	Less(int, int) bool
	Swap(int, int)
	Get(int) *Card
	Append(*Card)
	Peek() *Card
	Pop() *Card
	Push(*Card)
	Slot() image.Point
	SetBaizePos(image.Point)
	BaizePos() image.Point
	BaizeRect() image.Rectangle
	ScreenRect() image.Rectangle
	FannedBaizeRect() image.Rectangle
	FannedScreenRect() image.Rectangle
	PosAfter(*Card) image.Point
	Refan()
	IndexOf(*Card) int
	CanMoveTail([]*Card) (bool, error)
	MakeTail(*Card) []*Card
	ApplyToCards(func(*Card))
	GenericTailTapped([]*Card)
	GenericCollect()
	BuryCards(int)

	DrawStaticCards(*ebiten.Image)
	DrawTransitioningCards(*ebiten.Image)
	DrawFlippingCards(*ebiten.Image)
	DrawDraggingCards(*ebiten.Image)

	Update()
	CreateBackgroundImage() *ebiten.Image
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
