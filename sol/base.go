package sol

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/schriftbank"
)

const (
	BASE_MAGIC uint32 = 0xdeadbeef
)

type FanType int

const (
	FAN_NONE FanType = iota
	FAN_DOWN
	FAN_LEFT
	FAN_RIGHT
	FAN_DOWN3
	FAN_LEFT3
	FAN_RIGHT3
)

type MoveType int

const (
	MOVE_ANY MoveType = iota
	MOVE_ONE
	MOVE_ONE_PLUS
	MOVE_ONE_OR_ALL
)

const (
	CARD_FACE_FAN_FACTOR_V = 4
	CARD_FACE_FAN_FACTOR_H = 4
	CARD_BACK_FAN_FACTOR   = 8
)

var FanFactors [7]float64 = [7]float64{
	1.0,                    // FAN_NONE
	CARD_FACE_FAN_FACTOR_V, // FAN_DOWN
	CARD_FACE_FAN_FACTOR_H, // FAN_LEFT,
	CARD_FACE_FAN_FACTOR_H, // FAN_RIGHT,
	CARD_FACE_FAN_FACTOR_V, // FAN_DOWN3,
	CARD_FACE_FAN_FACTOR_H, // FAN_LEFT3,
	CARD_FACE_FAN_FACTOR_H, // FAN_RIGHT3,
}

const (
	RECYCLE_RUNE   = rune(0x2672)
	NORECYCLE_RUNE = rune(0x2613)
)

// Base is a generic container for cards
type Base struct {
	magic            uint32
	category         string
	slot             image.Point
	pos              image.Point
	pos1             image.Point // waste pos #1
	pos2             image.Point // waste pos #1
	fanType          FanType
	fanFactor        float64
	defaultFanFactor float64
	label            string
	symbol           rune
	cards            []*Card
	iface            Pile // TODO not sure about doing this
	img              *ebiten.Image
	scrunchDims      image.Point
	buddyPos         image.Point
}

func (base *Base) Ctor(iface Pile, category string, slot image.Point, fanType FanType) {
	base.magic = BASE_MAGIC
	base.category = category
	base.slot = slot
	base.fanType = fanType
	base.defaultFanFactor = FanFactors[fanType]
	base.fanFactor = base.defaultFanFactor
	base.cards = nil
	base.iface = iface
}

func (base *Base) Valid() bool {
	return base.magic == BASE_MAGIC
}

// Hidden returns true if this is off screen
func (base *Base) Hidden() bool {
	return base.slot.X < 0 || base.slot.Y < 0
}

// Empty returns true if this pile is empty.
// for use outside this chunk
func (base *Base) Empty() bool {
	return len(base.cards) == 0
}

// Len returns the number of cards in this pile.
// Len satisfies the sort.Interface interface.
// for use outside this chunk
func (base *Base) Len() int {
	return len(base.cards)
}

// Less satisfies the sort.Interface interface
func (base *Base) Less(i, j int) bool {
	c1 := base.cards[i]
	c2 := base.cards[j]
	return c1.Suit() < c2.Suit() && c1.Ordinal() < c2.Ordinal()
}

// Swap satisfies the sort.Interface interface
func (base *Base) Swap(i, j int) {
	base.cards[i], base.cards[j] = base.cards[j], base.cards[i]
}

func (base *Base) Label() string {
	return base.label
}

func (base *Base) SetLabel(label string) {
	base.label = label
	TheBaize.setFlag(dirtyPileBackgrounds)
}

func (base *Base) Rune() rune {
	return base.symbol
}

func (base *Base) SetRune(symbol rune) {
	base.symbol = symbol
	TheBaize.setFlag(dirtyPileBackgrounds)
}

// Get a *Card from this collection
func (base *Base) Get(i int) *Card {
	return base.cards[i]
}

// Append a *Card to this collection
func (base *Base) Append(c *Card) {
	base.cards = append(base.cards, c)
}

// Peek topmost Card of this Pile (a stack)
func (base *Base) Peek() *Card {
	if len(base.cards) == 0 {
		return nil
	}
	return base.cards[len(base.cards)-1]
}

// Pop a Card off the end of this Pile (a stack)
func (base *Base) Pop() *Card {
	if len(base.cards) == 0 {
		return nil
	}
	c := base.cards[len(base.cards)-1]
	base.cards = base.cards[:len(base.cards)-1]
	c.SetOwner(nil)
	c.FlipUp()
	base.Scrunch()
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (base *Base) Push(c *Card) {
	c.StopSpinning()
	pos := base.PosAfter(base.Peek()) // get this BEFORE appending card
	base.cards = append(base.cards, c)
	c.TransitionTo(pos)
	c.SetOwner(base.iface)
	if _, ok := (base.iface).(*Stock); ok {
		c.FlipDown()
	}
	base.Scrunch()
}

// func (base *Base) FindCard(ordinal, suit int) (*Card, int) {
// for i, c := range base.cards {
// if c.Ordinal() == ordinal && c.Suit() == suit {
// return c, i
// }
// }
// return nil, 0
// }

// Slot returns the virtual slot this pile is positioned at
// TODO to use fractional slots, scale the slot values up by, say, 10
func (base *Base) Slot() image.Point {
	return base.slot
}

// SetBaizePos sets the position of this Pile in Baize coords,
// and also sets the auxillary waste pile fanned positions
func (base *Base) SetBaizePos(pos image.Point) {
	base.pos = pos
	switch base.fanType {
	case FAN_DOWN3:
		base.pos1.X = base.pos.X
		base.pos1.Y = base.pos.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
		base.pos2.X = base.pos.X
		base.pos2.Y = base.pos1.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
	case FAN_LEFT3:
		base.pos1.X = base.pos.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		base.pos1.Y = base.pos.Y
		base.pos2.X = base.pos1.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		base.pos2.Y = base.pos.Y
	case FAN_RIGHT3:
		base.pos1.X = base.pos.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		base.pos1.Y = base.pos.Y
		base.pos2.X = base.pos1.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		base.pos2.Y = base.pos.Y
	}
	// println(base.category, base.pos.X, base.pos.Y)
}

func (base *Base) BaizePos() image.Point {
	return base.pos
}

func (base *Base) ScreenPos() image.Point {
	return base.pos.Add(TheBaize.dragOffset)
}

func (base *Base) BaizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = base.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

func (base *Base) ScreenRect() image.Rectangle {
	var r image.Rectangle = base.BaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

func (base *Base) FannedBaizeRect() image.Rectangle {
	var r image.Rectangle = base.BaizeRect()
	if len(base.cards) > 1 {
		var c *Card = base.Peek()
		if c.Dragging() {
			return r
		}
		var cPos = c.BaizePos()
		switch base.fanType {
		case FAN_NONE:
			// do nothing
		case FAN_RIGHT, FAN_RIGHT3:
			r.Max.X = cPos.X + CardWidth
		case FAN_LEFT, FAN_LEFT3:
			r.Max.X = cPos.X - CardWidth
		case FAN_DOWN, FAN_DOWN3:
			r.Max.Y = cPos.Y + CardHeight
		}
	}
	return r
}

func (base *Base) FannedScreenRect() image.Rectangle {
	var r image.Rectangle = base.FannedBaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

// PosAfter returns the position of the next card
func (base *Base) PosAfter(c *Card) image.Point {
	if len(base.cards) == 0 {
		return base.pos
	}
	var pos image.Point
	if c.Transitioning() {
		pos = c.dst
	} else {
		pos = c.pos
	}
	if pos.X == 0 && pos.Y == 0 {
		println("zero pos in PosAfter", base.category)
	}
	switch base.fanType {
	case FAN_NONE:
		// nothing to do
	case FAN_DOWN:
		if c.Prone() {
			pos.Y += int(float64(CardHeight) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.Y += int(float64(CardHeight) / base.fanFactor)
		}
	case FAN_LEFT:
		if c.Prone() {
			pos.X -= int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X -= int(float64(CardWidth) / base.fanFactor)
		}
	case FAN_RIGHT:
		if c.Prone() {
			pos.X += int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X += int(float64(CardWidth) / base.fanFactor)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		switch len(base.cards) {
		case 0:
			// nothing to do
		case 1:
			pos = base.pos1 // incoming card at slot 1
		case 2:
			pos = base.pos2 // incoming card at slot 2
		default:
			pos = base.pos2 // incoming card at slot 2
			// top card needs to transition from slot[2] to slot[1]
			i := len(base.cards) - 1
			base.cards[i].TransitionTo(base.pos1)
			// mid card needs to transition from slot[1] to slot[0]
			// all other cards to slot[0]
			for i > 0 {
				i--
				base.cards[i].TransitionTo(base.pos)
			}
		}
	}
	return pos
}

func (base *Base) Refan() {
	var doFan3 bool = false
	switch base.fanType {
	case FAN_NONE:
		for _, c := range base.cards {
			c.TransitionTo(base.pos)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		for _, c := range base.cards {
			c.TransitionTo(base.pos)
		}
		doFan3 = true
	case FAN_DOWN, FAN_LEFT, FAN_RIGHT:
		var pos = base.pos
		var i = 0
		for _, c := range base.cards {
			c.TransitionTo(pos)
			pos = base.PosAfter(base.cards[i])
			i++
		}
	}

	if doFan3 {
		switch len(base.cards) {
		case 0:
		case 1:
			// nothing to do
		case 2:
			c := base.cards[1]
			c.TransitionTo(base.pos1)
		default:
			i := len(base.cards)
			i--
			c := base.cards[i]
			c.TransitionTo(base.pos2)
			i--
			c = base.cards[i]
			c.TransitionTo(base.pos1)
		}
	}
}

func (base *Base) IndexOf(card *Card) (int, error) {
	for i, c := range base.cards {
		if c == card {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%s not found in %s", card.String(), base.category)
}

func (base *Base) MakeTail(c *Card) []*Card {
	var tail []*Card
	if len(base.cards) > 0 {
		for i, pc := range base.cards {
			if pc == c {
				tail = base.cards[i:]
				break
			}
		}
	}
	if len(tail) == 0 {
		log.Panic("Pile.makeTail made an empty tail")
	}
	return tail
}

// ApplyToCards applies a function to each card in the pile
// caller must use a method expression, eg (*Card).StartSpinning, yielding a function value
// with a regular first parameter taking the place of the receiver
func (base *Base) ApplyToCards(fn func(*Card)) {
	for _, c := range base.cards {
		fn(c)
	}
}

func (base *Base) TailTapped(tail []*Card) {
	if len(tail) != 1 {
		return
	}
	c := tail[0]
	for _, fp := range TheBaize.foundations {
		if ok, _ := fp.CanAcceptCard(c); ok {
			MoveCard(base.iface, fp)
			break
		}
	}
}

func (base *Base) Collect() {
	for _, fp := range TheBaize.foundations {
		for {
			// loop to get as many cards as possible from this pile
			if base.Empty() {
				return
			}
			if ok, _ := fp.CanAcceptCard(base.Peek()); !ok {
				// this foundation doesn't want this card; onto the next one
				break
			}
			MoveCard(base.iface, fp)
		}
	}
}

func (base *Base) Reset() {
	base.cards = base.cards[:0]
}

func (base *Base) DrawStaticCards(screen *ebiten.Image) {
	for _, c := range base.cards {
		if !(c.Transitioning() || c.Flipping() || c.Dragging()) {
			c.Draw(screen)
		}
	}
}

func (base *Base) DrawTransitioningCards(screen *ebiten.Image) {
	for _, c := range base.cards {
		if c.Transitioning() {
			c.Draw(screen)
		}
	}
}

func (base *Base) DrawFlippingCards(screen *ebiten.Image) {
	for _, c := range base.cards {
		if c.Flipping() {
			c.Draw(screen)
		}
	}
}

func (base *Base) DrawDraggingCards(screen *ebiten.Image) {
	for _, c := range base.cards {
		if c.Dragging() {
			c.Draw(screen)
		}
	}
}

func (base *Base) Update() {
	for _, card := range base.cards {
		card.Update()
	}
}

func (b *Base) CreateBackgroundImage() {
	if b.Hidden() {
		return
	}
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	if b.symbol != 0 {
		dc.SetFontFace(schriftbank.CardSymbolLarge)
		dc.DrawStringAnchored(string(b.symbol), float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
	} else if b.label != "" {
		dc.SetFontFace(schriftbank.CardOrdinalLarge)
		dc.DrawStringAnchored(b.label, float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
	}
	dc.Stroke()
	b.img = ebiten.NewImageFromImage(dc.Image())
}

func (b *Base) Draw(screen *ebiten.Image) {
	if b.img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.pos.X+TheBaize.dragOffset.X), float64(b.pos.Y+TheBaize.dragOffset.Y))
	screen.DrawImage(b.img, op)
}
