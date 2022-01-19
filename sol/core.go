package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/sound"
)

const (
	coremagic uint32 = 0xdeadbeef
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
	MOVE_NONE MoveType = iota
	MOVE_ANY
	MOVE_ONE
	MOVE_ONE_PLUS
	MOVE_ONE_OR_ALL
)

const (
	CARD_FACE_FAN_FACTOR_V = 3.7
	CARD_FACE_FAN_FACTOR_H = 4
	CARD_BACK_FAN_FACTOR   = 8
)

var DefaultFanFactor [7]float64 = [7]float64{
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
type Core struct {
	magic     uint32
	category  string
	slot      image.Point
	fanType   FanType
	moveType  MoveType
	cards     []*Card
	pos       image.Point
	pos1      image.Point // waste pos #1
	pos2      image.Point // waste pos #1
	fanFactor float64
	// buddyPos    image.Point
	label  string
	symbol rune
	img    *ebiten.Image
	target bool // experimental, might delete later, IDK
}

func NewCore(category string, slot image.Point, fanType FanType, moveType MoveType) Core {
	self := Core{
		// static
		magic:    coremagic,
		category: category,
		slot:     slot,
		fanType:  fanType,
		moveType: moveType,
		// dynamic
		fanFactor: DefaultFanFactor[fanType],
	}
	return self
}

func (self *Core) Valid() bool {
	return self != nil && self.magic == coremagic
}

func (self *Core) Reset() {
	self.cards = self.cards[:0]
	self.fanFactor = DefaultFanFactor[self.fanType]
}

// Hidden returns true if this is off screen
func (self *Core) Hidden() bool {
	return self.slot.X < 0 || self.slot.Y < 0
}

func (self *Core) IsStock() bool {
	return self.category == "Stock"
}

func (self *Core) IsTableau() bool {
	return self.category == "Tableau"
}

func (self *Core) Cards() []*Card {
	return self.cards
}

func (self *Core) FanType() FanType {
	return self.fanType
}

func (self *Core) SetFanType(fanType FanType) {
	self.fanType = fanType
}

func (self *Core) MoveType() MoveType {
	return self.moveType
}

func (self *Core) Label() string {
	return self.label
}

func (self *Core) SetLabel(label string) {
	if self.label != label {
		if self.IsTableau() && ThePreferences.Relaxed && TheBaize.script.Info().relaxable {
		} else {
			self.label = label
			TheBaize.setFlag(dirtyPileBackgrounds)
		}
	}
}

func (self *Core) Rune() rune {
	return self.symbol
}

func (self *Core) SetRune(symbol rune) {
	if self.symbol != symbol {
		self.symbol = symbol
		TheBaize.setFlag(dirtyPileBackgrounds)
	}
}

func (self *Core) Target() bool {
	return self.target
}

func (self *Core) SetTarget(target bool) {
	self.target = target
}

// Empty returns true if this core is empty.
// for use outside this chunk
func (self *Core) Empty() bool {
	return len(self.cards) == 0
}

// Len returns the number of cards in this core.
// Len satisfies the sort.Interface interface.
// for use outside this chunk
func (self *Core) Len() int {
	return len(self.cards)
}

// Less satisfies the sort.Interface interface
func (self *Core) Less(i, j int) bool {
	c1 := self.cards[i]
	c2 := self.cards[j]
	return c1.Suit() < c2.Suit() && c1.Ordinal() < c2.Ordinal()
}

// Swap satisfies the sort.Interface interface
func (self *Core) Swap(i, j int) {
	self.cards[i], self.cards[j] = self.cards[j], self.cards[i]
}

// Get a *Card from this collection
func (self *Core) Get(i int) *Card {
	return self.cards[i]
}

// Append a *Card to this collection
func (self *Core) Append(c *Card) {
	self.cards = append(self.cards, c)
}

// Delete a *Card from this collection
func (self *Core) Delete(index int) {
	self.cards = append(self.cards[:index], self.cards[index+1:]...)
}

// Peek topmost Card of this Pile (a stack)
func (self *Core) Peek() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	return self.cards[len(self.cards)-1]
}

// Pop a Card off the end of this Pile (a stack)
func (self *Core) Pop() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	c := self.cards[len(self.cards)-1]
	self.cards = self.cards[:len(self.cards)-1]
	c.SetOwner(nil)
	c.FlipUp()
	TheBaize.setFlag(dirtyCardPositions)
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (self *Core) Push(c *Card) {
	var pos image.Point
	if len(self.cards) == 0 {
		pos = self.pos
	} else {
		pos = self.PosAfter(self.Peek()) // get this BEFORE appending card
	}

	self.cards = append(self.cards, c)
	c.SetOwner(FindCardOwner(c))
	// c.SetOwner(self)
	c.TransitionTo(pos)

	if self.IsStock() {
		c.FlipDown()
	}
	TheBaize.setFlag(dirtyCardPositions)
}

// Slot returns the virtual slot this core is positioned at
// TODO to use fractional slots, scale the slot values up by, say, 10
func (self *Core) Slot() image.Point {
	return self.slot
}

func (self *Core) SetSlot(slot image.Point) {
	self.slot = slot
}

// SetBaizePos sets the position of this Pile in Baize coords,
// and also sets the auxillary waste core fanned positions
func (self *Core) SetBaizePos(pos image.Point) {
	self.pos = pos
	switch self.fanType {
	case FAN_DOWN3:
		self.pos1.X = self.pos.X
		self.pos1.Y = self.pos.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
		self.pos2.X = self.pos.X
		self.pos2.Y = self.pos1.Y + int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V)
	case FAN_LEFT3:
		self.pos1.X = self.pos.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		self.pos1.Y = self.pos.Y
		self.pos2.X = self.pos1.X - int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		self.pos2.Y = self.pos.Y
	case FAN_RIGHT3:
		self.pos1.X = self.pos.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		self.pos1.Y = self.pos.Y
		self.pos2.X = self.pos1.X + int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H)
		self.pos2.Y = self.pos.Y
	}
	// println(base.category, base.pos.X, base.pos.Y)
}

func (self *Core) BaizePos() image.Point {
	return self.pos
}

func (self *Core) ScreenPos() image.Point {
	return self.pos.Add(TheBaize.dragOffset)
}

func (self *Core) BaizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = self.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

func (self *Core) ScreenRect() image.Rectangle {
	var r image.Rectangle = self.BaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

func (self *Core) FannedBaizeRect() image.Rectangle {
	var r image.Rectangle = self.BaizeRect()
	if len(self.cards) > 1 {
		var c *Card = self.Peek()
		// if c.Dragging() {
		// 	return r
		// }
		var cPos = c.BaizePos()
		switch self.fanType {
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

func (self *Core) FannedScreenRect() image.Rectangle {
	var r image.Rectangle = self.FannedBaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

// PosAfter returns the position of the next card
func (self *Core) PosAfter(c *Card) image.Point {
	if len(self.cards) == 0 {
		println("Panic! PosAfter called in impossible way")
		return self.pos
	}
	var pos image.Point
	if c.Transitioning() {
		pos = c.dst
	} else {
		pos = c.pos
	}
	if pos.X <= 0 && pos.Y <= 0 {
		// the card is still at 0,0 where it started life
		// and is yet to have pos calculated from the pile slot
		// println("zero pos in PosAfter", self.category)
		return pos
	}
	switch self.fanType {
	case FAN_NONE:
		// nothing to do
	case FAN_DOWN:
		if c.Prone() {
			pos.Y += int(float64(CardHeight) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.Y += int(float64(CardHeight) / self.fanFactor)
		}
	case FAN_LEFT:
		if c.Prone() {
			pos.X -= int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X -= int(float64(CardWidth) / self.fanFactor)
		}
	case FAN_RIGHT:
		if c.Prone() {
			pos.X += int(float64(CardWidth) / float64(CARD_BACK_FAN_FACTOR))
		} else {
			pos.X += int(float64(CardWidth) / self.fanFactor)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		switch len(self.cards) {
		case 0:
			// nothing to do
		case 1:
			pos = self.pos1 // incoming card at slot 1
		case 2:
			pos = self.pos2 // incoming card at slot 2
		default:
			pos = self.pos2 // incoming card at slot 2
			// top card needs to transition from slot[2] to slot[1]
			i := len(self.cards) - 1
			self.cards[i].TransitionTo(self.pos1)
			// mid card needs to transition from slot[1] to slot[0]
			// all other cards to slot[0]
			for i > 0 {
				i--
				self.cards[i].TransitionTo(self.pos)
			}
		}
	}
	return pos
}

func (self *Core) Refan() {
	var doFan3 bool = false
	switch self.fanType {
	case FAN_NONE:
		for _, c := range self.cards {
			c.TransitionTo(self.pos)
		}
	case FAN_DOWN3, FAN_LEFT3, FAN_RIGHT3:
		for _, c := range self.cards {
			c.TransitionTo(self.pos)
		}
		doFan3 = true
	case FAN_DOWN, FAN_LEFT, FAN_RIGHT:
		var pos = self.pos
		var i = 0
		for _, c := range self.cards {
			c.TransitionTo(pos)
			pos = self.PosAfter(self.cards[i])
			i++
		}
	}

	if doFan3 {
		switch len(self.cards) {
		case 0:
		case 1:
			// nothing to do
		case 2:
			c := self.cards[1]
			c.TransitionTo(self.pos1)
		default:
			i := len(self.cards)
			i--
			c := self.cards[i]
			c.TransitionTo(self.pos2)
			i--
			c = self.cards[i]
			c.TransitionTo(self.pos1)
		}
	}
}

func (self *Core) IndexOf(card *Card) int {
	for i, c := range self.cards {
		if c == card {
			return i
		}
	}
	return -1
}

func (self *Core) CanMoveTail(tail []*Card) (bool, error) {
	if AnyCardsProne(tail) {
		if _, isStock := (tail[0].Owner()).(*Stock); !isStock {
			return false, errors.New("Cannot move a face down card")
		}
	}
	switch self.moveType {
	case MOVE_NONE:
		return false, fmt.Errorf("Cannot move a card from a %s", self.category)
	case MOVE_ANY:
		// well, that was easy
	case MOVE_ONE:
		if len(tail) > 1 {
			return false, fmt.Errorf("Can only move one card from a %s", self.category)
		}
	case MOVE_ONE_PLUS:
		// don't know destination, so we allow this as MOVE_ANY
	case MOVE_ONE_OR_ALL:
		if len(tail) == 1 {
			// that's okay
		} else if len(tail) == self.Len() {
			// that's okay too
		} else {
			return false, errors.New("Only move one card, or the whole pile")
		}
	}
	return true, nil
}

func (self *Core) MakeTail(c *Card) []*Card {
	var tail []*Card
	if len(self.cards) > 0 {
		for i, pc := range self.cards {
			if pc == c {
				tail = self.cards[i:]
				break
			}
		}
	}
	if len(tail) == 0 {
		log.Panic("Pile.makeTail made an empty tail")
	}
	return tail
}

// ApplyToCards applies a function to each card in the core
// caller must use a method expression, eg (*Card).StartSpinning, yielding a function value
// with a regular first parameter taking the place of the receiver
func (self *Core) ApplyToCards(fn func(*Card)) {
	for _, c := range self.cards {
		fn(c)
	}
}

// BuryCards moves cards with the specified ordinal to the beginning of the core
func (self *Core) BuryCards(ordinal int) {
	tmp := make([]*Card, 0, cap(self.cards))
	for _, c := range self.cards {
		if c.Ordinal() == ordinal {
			tmp = append(tmp, c)
		}
	}
	for _, c := range self.cards {
		if c.Ordinal() != ordinal {
			tmp = append(tmp, c)
		}
	}
	self.Reset()
	for i := 0; i < len(tmp); i++ {
		self.Push(tmp[i])
	}
	self.Refan()
	// nb the card owner does not change
}

// default behaviours for all pile types, that can be over-ridden by providing (eg) *Stock.Collect

func (self *Core) CanAcceptCard(*Card) (bool, error)   { return false, nil }
func (self *Core) CanAcceptTail([]*Card) (bool, error) { return false, nil }

func (self *Core) TailTapped(tail []*Card) {
	tappedCard := tail[0]
	src := tappedCard.Owner()
	if len(tail) == 1 {
		for _, fp := range TheBaize.script.Foundations() {
			if ok, _ := fp.CanAcceptCard(tappedCard); ok {
				MoveCard(self, fp)
				return
			}
		}
	}
	var longestPile *Tableau
	for _, tp := range TheBaize.script.Tableaux() {
		if tp == src {
			continue
		}
		// can the tail be moved in general (MoveType check)?
		if ok, _ := src.CanMoveTail(tail); ok {
			// can the dst accept the tail?
			if ok, _ := tp.CanAcceptTail(tail); ok {
				// is the tail conformant enough to move?
				if ok, _ := TheBaize.script.TailMoveError(tail); ok {
					// very annoying to move cards to an empty pile
					// in games where creating empty piles is useful
					if tp.Empty() && tp.Label() == "" {
						continue
					}
					if !tp.Empty() {
						if tail[0].Suit() == tp.Peek().Suit() {
							longestPile = tp
							break
						}
					}
					if longestPile == nil || tp.Len() > longestPile.Len() {
						longestPile = tp
					}
				}
			}
		}
	}
	if longestPile != nil {
		MoveCards(src, src.IndexOf(tappedCard), longestPile)
	} else {
		sound.Play("Blip")
	}
}

func (self *Core) Collect() {
	for _, fp := range TheBaize.script.Foundations() {
		for {
			// loop to get as many cards as possible from this pile
			if self.Empty() {
				return
			}
			if ok, _ := fp.CanAcceptCard(self.Peek()); !ok {
				// this foundation doesn't want this card; onto the next one
				break
			}
			MoveCard(self, fp)
		}
	}
}

func (self *Core) Conformant() bool   { return false }
func (self *Core) Complete() bool     { return false }
func (self *Core) UnsortedPairs() int { return 0 }

func (self *Core) DrawStaticCards(screen *ebiten.Image) {
	for _, c := range self.cards {
		if !(c.Transitioning() || c.Flipping() || c.Dragging()) {
			c.Draw(screen)
		}
	}
}

func (self *Core) DrawTransitioningCards(screen *ebiten.Image) {
	for _, c := range self.cards {
		if c.Transitioning() && !c.Flipping() {
			c.Draw(screen)
		}
	}
}

func (self *Core) DrawFlippingCards(screen *ebiten.Image) {
	for _, c := range self.cards {
		if c.Flipping() {
			c.Draw(screen)
		}
	}
}

func (self *Core) DrawDraggingCards(screen *ebiten.Image) {
	for _, c := range self.cards {
		if c.Dragging() {
			c.Draw(screen)
		}
	}
}

func (self *Core) Update() {
	for _, card := range self.cards {
		card.Update()
	}
}

func (self *Core) CreateBackgroundImage() {
	self.img = nil
	if CardWidth == 0 || CardHeight == 0 {
		println("zero dimension in CreateCardShadowImage, unliked in wasm")
		return
		// log.Panic("zero dimension in CreateCardShadowImage, unliked in wasm")
	}
	if self.Hidden() {
		// off-screen? don't bother
		return
	}
	if self.category == "Reserve" {
		// don't draw anything for reserve cores
		return
	}
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	switch self.category {
	case "Discard":
		dc.Fill()
	default:
		if self.symbol != 0 {
			// usually the recycle symbol
			dc.SetFontFace(schriftbank.CardSymbolLarge)
			dc.DrawStringAnchored(string(self.symbol), float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
		} else if self.label != "" {
			// usually the "place card in an empty core" constraint
			dc.SetFontFace(schriftbank.CardOrdinalLarge)
			dc.DrawStringAnchored(self.label, float64(CardWidth)*0.5, float64(CardHeight)*0.45, 0.5, 0.5)
		}
	}
	dc.Stroke()
	self.img = ebiten.NewImageFromImage(dc.Image())
}

func (self *Core) Draw(screen *ebiten.Image) {
	if self.img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(self.pos.X+TheBaize.dragOffset.X), float64(self.pos.Y+TheBaize.dragOffset.Y))
	if self.target && len(self.cards) == 0 {
		op.ColorM.Scale(0.75, 0.75, 0.75, 1)
	}

	if self.symbol != 0 {
		if pt := image.Pt(ebiten.CursorPosition()); pt.In(self.ScreenRect()) {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				op.GeoM.Translate(2, 2)
			}
		}
	}

	// if DebugMode {
	// 	if sz := self.SizeWithFanFactor(self.fanFactor); sz != 0 {
	// 		switch self.fanType {
	// 		case FAN_DOWN:
	// 			rect := self.FannedScreenRect()
	// 			ebitenutil.DrawRect(screen,
	// 				float64(rect.Min.X),
	// 				float64(rect.Min.Y),
	// 				float64(rect.Max.X-rect.Min.X),
	// 				float64(sz),
	// 				color.RGBA{0, 0, 0, 32})
	// 		}
	// 	}
	// }

	screen.DrawImage(self.img, op)
}
