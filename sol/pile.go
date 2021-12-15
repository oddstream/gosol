package sol

import (
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
type Pile struct {
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
	subtype          SubtypeAPI
	img              *ebiten.Image
	scrunchDims      image.Point
	buddyPos         image.Point
	target           bool // experimental, might delete later, IDK
}

func (p *Pile) Ctor(subtype SubtypeAPI, category string, slot image.Point, fanType FanType) {
	p.magic = BASE_MAGIC
	p.category = category
	p.slot = slot
	p.fanType = fanType
	p.defaultFanFactor = FanFactors[fanType]
	p.fanFactor = p.defaultFanFactor
	p.cards = nil
	p.subtype = subtype
	TheBaize.piles = append(TheBaize.piles, p) // TODO nasty
}

func (p *Pile) Valid() bool {
	return p.magic == BASE_MAGIC
}

// Hidden returns true if this is off screen
func (p *Pile) Hidden() bool {
	return p.slot.X < 0 || p.slot.Y < 0
}

// Empty returns true if this pile is empty.
// for use outside this chunk
func (p *Pile) Empty() bool {
	return len(p.cards) == 0
}

// Len returns the number of cards in this pile.
// Len satisfies the sort.Interface interface.
// for use outside this chunk
func (p *Pile) Len() int {
	return len(p.cards)
}

// Less satisfies the sort.Interface interface
func (p *Pile) Less(i, j int) bool {
	c1 := p.cards[i]
	c2 := p.cards[j]
	return c1.Suit() < c2.Suit() && c1.Ordinal() < c2.Ordinal()
}

// Swap satisfies the sort.Interface interface
func (p *Pile) Swap(i, j int) {
	p.cards[i], p.cards[j] = p.cards[j], p.cards[i]
}

func (p *Pile) Label() string {
	return p.label
}

func (p *Pile) SetLabel(label string) {
	p.label = label
	TheBaize.setFlag(dirtyPileBackgrounds)
}

func (p *Pile) Rune() rune {
	return p.symbol
}

func (p *Pile) SetRune(symbol rune) {
	p.symbol = symbol
	TheBaize.setFlag(dirtyPileBackgrounds)
}

// Get a *Card from this collection
func (p *Pile) Get(i int) *Card {
	return p.cards[i]
}

// Append a *Card to this collection
func (p *Pile) Append(c *Card) {
	p.cards = append(p.cards, c)
}

// Peek topmost Card of this Pile (a stack)
func (p *Pile) Peek() *Card {
	if len(p.cards) == 0 {
		return nil
	}
	return p.cards[len(p.cards)-1]
}

// Pop a Card off the end of this Pile (a stack)
func (p *Pile) Pop() *Card {
	if len(p.cards) == 0 {
		return nil
	}
	c := p.cards[len(p.cards)-1]
	p.cards = p.cards[:len(p.cards)-1]
	c.SetOwner(nil)
	c.FlipUp()
	p.Scrunch()
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (p *Pile) Push(c *Card) {
	// c.StopSpinning()

	var pos image.Point
	if len(p.cards) == 0 {
		pos = p.pos
	} else {
		pos = p.PosAfter(p.Peek()) // get this BEFORE appending card
	}

	p.cards = append(p.cards, c)
	c.SetOwner(p)
	c.TransitionTo(pos)

	if _, ok := (p.subtype).(*Stock); ok {
		c.FlipDown()
	}
	p.Scrunch()
}

// Slot returns the virtual slot this pile is positioned at
// TODO to use fractional slots, scale the slot values up by, say, 10
func (p *Pile) Slot() image.Point {
	return p.slot
}

// SetBaizePos sets the position of this Pile in Baize coords,
// and also sets the auxillary waste pile fanned positions
func (p *Pile) SetBaizePos(pos image.Point) {
	p.pos = pos
	switch p.fanType {
	case FAN_DOWN3:
		p.pos1.X = p.pos.X
		p.pos1.Y = p.pos.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
		p.pos2.X = p.pos.X
		p.pos2.Y = p.pos1.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
	case FAN_LEFT3:
		p.pos1.X = p.pos.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	case FAN_RIGHT3:
		p.pos1.X = p.pos.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos1.Y = p.pos.Y
		p.pos2.X = p.pos1.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		p.pos2.Y = p.pos.Y
	}
	// println(base.category, base.pos.X, base.pos.Y)
}

func (p *Pile) BaizePos() image.Point {
	return p.pos
}

func (p *Pile) ScreenPos() image.Point {
	return p.pos.Add(TheBaize.dragOffset)
}

func (p *Pile) BaizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = p.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

func (p *Pile) ScreenRect() image.Rectangle {
	var r image.Rectangle = p.BaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

func (p *Pile) FannedBaizeRect() image.Rectangle {
	var r image.Rectangle = p.BaizeRect()
	if len(p.cards) > 1 {
		var c *Card = p.Peek()
		if c.Dragging() {
			return r
		}
		var cPos = c.BaizePos()
		switch p.fanType {
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

func (p *Pile) FannedScreenRect() image.Rectangle {
	var r image.Rectangle = p.FannedBaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

// PosAfter returns the position of the next card
func (p *Pile) PosAfter(c *Card) image.Point {
	if len(p.cards) == 0 {
		return p.pos
	}
	var pos image.Point
	if c.Transitioning() {
		pos = c.dst
	} else {
		pos = c.pos
	}
	if pos.X == 0 && pos.Y == 0 {
		println("zero pos in PosAfter", p.category)
	}
	switch p.fanType {
	case FAN_NONE:
		// nothing to do
	case FAN_DOWN:
		if c.Prone() {
			pos.Y += int(float64(CardHeight) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.Y += int(float64(CardHeight) / p.fanFactor)
		}
	case FAN_LEFT:
		if c.Prone() {
			pos.X -= int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X -= int(float64(CardWidth) / p.fanFactor)
		}
	case FAN_RIGHT:
		if c.Prone() {
			pos.X += int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X += int(float64(CardWidth) / p.fanFactor)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		switch len(p.cards) {
		case 0:
			// nothing to do
		case 1:
			pos = p.pos1 // incoming card at slot 1
		case 2:
			pos = p.pos2 // incoming card at slot 2
		default:
			pos = p.pos2 // incoming card at slot 2
			// top card needs to transition from slot[2] to slot[1]
			i := len(p.cards) - 1
			p.cards[i].TransitionTo(p.pos1)
			// mid card needs to transition from slot[1] to slot[0]
			// all other cards to slot[0]
			for i > 0 {
				i--
				p.cards[i].TransitionTo(p.pos)
			}
		}
	}
	return pos
}

func (p *Pile) Refan() {
	var doFan3 bool = false
	switch p.fanType {
	case FAN_NONE:
		for _, c := range p.cards {
			c.TransitionTo(p.pos)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		for _, c := range p.cards {
			c.TransitionTo(p.pos)
		}
		doFan3 = true
	case FAN_DOWN, FAN_LEFT, FAN_RIGHT:
		var pos = p.pos
		var i = 0
		for _, c := range p.cards {
			c.TransitionTo(pos)
			pos = p.PosAfter(p.cards[i])
			i++
		}
	}

	if doFan3 {
		switch len(p.cards) {
		case 0:
		case 1:
			// nothing to do
		case 2:
			c := p.cards[1]
			c.TransitionTo(p.pos1)
		default:
			i := len(p.cards)
			i--
			c := p.cards[i]
			c.TransitionTo(p.pos2)
			i--
			c = p.cards[i]
			c.TransitionTo(p.pos1)
		}
	}
}

func (p *Pile) IndexOf(card *Card) int {
	for i, c := range p.cards {
		if c == card {
			return i
		}
	}
	return -1
}

func (p *Pile) MakeTail(c *Card) []*Card {
	var tail []*Card
	if len(p.cards) > 0 {
		for i, pc := range p.cards {
			if pc == c {
				tail = p.cards[i:]
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
func (p *Pile) ApplyToCards(fn func(*Card)) {
	for _, c := range p.cards {
		fn(c)
	}
}

func (p *Pile) GenericTailTapped(tail []*Card) {
	if len(tail) != 1 {
		return
	}
	c := tail[0]
	for _, fp := range TheBaize.foundations {
		if ok, _ := fp.subtype.CanAcceptCard(c); ok {
			MoveCard(p, fp)
			break
		}
	}
}

func (p *Pile) GenericCollect() {
	for _, fp := range TheBaize.foundations {
		for {
			// loop to get as many cards as possible from this pile
			if p.Empty() {
				return
			}
			if ok, _ := fp.subtype.CanAcceptCard(p.Peek()); !ok {
				// this foundation doesn't want this card; onto the next one
				break
			}
			MoveCard(p, fp)
		}
	}
}

func (p *Pile) GenericReset() {
	p.cards = p.cards[:0]
	p.label = ""
	p.symbol = 0
}

func (p *Pile) DrawStaticCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if !(c.Transitioning() || c.Flipping() || c.Dragging()) {
			c.Draw(screen)
		}
	}
}

func (p *Pile) DrawTransitioningCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.Transitioning() {
			c.Draw(screen)
		}
	}
}

func (p *Pile) DrawFlippingCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.Flipping() {
			c.Draw(screen)
		}
	}
}

func (p *Pile) DrawDraggingCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.Dragging() {
			c.Draw(screen)
		}
	}
}

func (p *Pile) Update() {
	for _, card := range p.cards {
		card.Update()
	}
}

func (p *Pile) CreateBackgroundImage() *ebiten.Image {
	if CardWidth == 0 || CardHeight == 0 {
		println("zero dimension in CreateCardShadowImage, unliked in wasm")
		return nil
		// log.Panic("zero dimension in CreateCardShadowImage, unliked in wasm")
	}
	if p.Hidden() {
		return nil
	}
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	switch (p.subtype).(type) {
	case *Discard:
		dc.Fill()
	default:
		if p.symbol != 0 {
			dc.SetFontFace(schriftbank.CardSymbolLarge)
			dc.DrawStringAnchored(string(p.symbol), float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
		} else if p.label != "" {
			dc.SetFontFace(schriftbank.CardOrdinalLarge)
			dc.DrawStringAnchored(p.label, float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
		}
	}
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func (p *Pile) Draw(screen *ebiten.Image) {
	if p.img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.pos.X+TheBaize.dragOffset.X), float64(p.pos.Y+TheBaize.dragOffset.Y))
	if p.target && len(p.cards) == 0 {
		// op.GeoM.Translate(-4, -4)
		// screen.DrawImage(CardHighlightImage, op)
		// op.GeoM.Translate(4, 4)
		op.ColorM.Scale(0.75, 0.75, 0.75, 1)
	}
	screen.DrawImage(p.img, op)
}
