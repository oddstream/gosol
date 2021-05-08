package sol

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

const (
	backFanFactor int = 5
	faceFanFactor int = 3
)

// var DefaultBuild = map[string]int{
// 	"Cell":       0,
// 	"Foundation": 21,
// 	"Reserve":    0,
// 	"Stock":      0,
// 	"Tableau":    42,
// 	"Waste":      15,
// }

// var DefaultMove = map[string]int{
// 	"Cell":       0,
// 	"Foundation": 0,
// 	"Reserve":    0,
// 	"Stock":      0,
// 	"Tableau":    42,
// 	"Waste":      15,
// }

// Pile is a generic container for cards
type Pile struct {
	PileInfo
	localAccept       int           // ordinal this pile can accept when empty (0=accept anything, 99=won't accept anything)
	localRecycles     int           // number of recycles left (stock only, could be in Baize)
	Cards             []*Card       // array of cards, managed as a stack
	Tail              []*Card       // array of cards currently being dragged
	scrunchSize       int           // relative (in card positions) size of scrunch height/width
	scrunchPercentage int           // percentage of compression of fanned cards so they fit on screen (but are harder to read)
	backgroundImage   *ebiten.Image // rounded rect for this Pile, optionally contains Accept/Recycle symbol
}

// NewPile create and fills in a Pile object
func NewPile(pi PileInfo) *Pile {
	p := &Pile{PileInfo: pi}
	p.Reset()
	return p
}

// Reset the pile
func (p *Pile) Reset() {
	p.localAccept, _ = p.GetIntAttribute("Accept")
	p.localRecycles, _ = p.GetIntAttribute("Recycles")
	// p.Cards
	// p.Tail
	// leave scrunchSize, it's calculated
	p.scrunchPercentage = 100
	p.CreateBackgroundImage()
}

// GetIntAttribute gets an integer Pile attribute
func (p *Pile) GetIntAttribute(key string) (int, bool) {
	str, exists := p.Attributes[key]
	if !exists {
		return 0, false
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(str + " is not an int")
	}
	return i, true
}

// GetStringAttribute gets a string Pile attribute
func (p *Pile) GetStringAttribute(key string) string {
	str, exists := p.Attributes[key]
	if exists {
		return str
	}
	return ""
}

// GetBoolAttribute gets a boolean Pile attribute
func (p *Pile) GetBoolAttribute(key string) bool {
	str, exists := p.Attributes[key]
	if exists {
		value, err := strconv.ParseBool(str)
		if err != nil {
			log.Panic("expecting a bool, got ", str)
		}
		return value
	}
	return false
}

func (p *Pile) CreateBackgroundImage() {

	invisible := p.GetBoolAttribute("Invisible")
	if invisible {
		p.backgroundImage = nil
		return
	}

	dc := gg.NewContext(CardWidth, CardHeight)

	dc.SetColor(colorPile)
	dc.SetLineWidth(2)

	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), cardCornerRadius())
	dc.Stroke()

	switch {
	case p.localAccept > 0:
		var str string
		dc.SetFontFace(schriftbank.CardOrdinal)
		if p.localAccept <= 13 {
			str = util.OrdinalToShortString(p.localAccept)
		} else {
			str = "x"
		}
		dc.DrawStringAnchored(str, float64(CardWidth)/3.333, float64(CardHeight)/6.666, 0.5, 0.5)
		// dc.SetLineWidth(1)
		// dc.DrawLine(0, float64(CardHeight)/6.666, float64(CardWidth), float64(CardHeight)/6.666)
		// dc.DrawLine(float64(CardWidth)/3.333, 0, float64(CardWidth)/3.333, float64(CardHeight))
		// dc.Stroke()

	case strings.HasPrefix(p.Class, "Stock"): // never StockSpider?
		dc.SetFontFace(schriftbank.CardSymbolLarge)
		if p.localRecycles == 0 {
			// anything put here either doesn't render (0x1F6AB) or looks ugly
			dc.SetFontFace(schriftbank.CardSymbolLarge)
			dc.DrawStringAnchored(string(rune(0x2613)), float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.4)
			// dc.DrawStringAnchored("X", float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.5)
		} else {
			dc.DrawStringAnchored(string(rune(0x2672)), float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.4)
		}
		dc.Stroke()

	case p.Class == "Reserve":
		dc.SetFontFace(schriftbank.CardOrdinal)
		dc.DrawStringAnchored("x", float64(CardWidth)/3.333, float64(CardHeight)/6.666, 0.5, 0.5)
		dc.Stroke()
	}

	p.backgroundImage = ebiten.NewImageFromImage(dc.Image())
}

// BaizePosition returns the x,y baize coords of this pile
func (p *Pile) BaizePosition() (int, int) {
	return LeftMargin + int(p.X*PilePositionType(CardWidth+PilePaddingX)), TopMargin + int(p.Y*PilePositionType(CardHeight+PilePaddingY))
}

// ScreenPosition returns the x,y baize coords of this pile
func (p *Pile) ScreenPosition() (int, int) {
	x, y := p.BaizePosition()
	x += TheBaize.DragOffsetX
	y += TheBaize.DragOffsetY
	return x, y
}

// BaizeRect gives the x,y baize coords of the pile's top left and bottom right corners
func (p *Pile) BaizeRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0 = p.BaizePosition()
	x1 = x0 + CardWidth
	y1 = y0 + CardHeight
	return // using named return parameters
}

// ScreenRect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) ScreenRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0, x1, y1 = p.BaizeRect()
	x0 += TheBaize.DragOffsetX
	x1 += TheBaize.DragOffsetX
	y0 += TheBaize.DragOffsetY
	y1 += TheBaize.DragOffsetY
	return // using named return parameters
}

// FannedBaizeRect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) FannedBaizeRect() (x0 int, y0 int, x1 int, y1 int) {
	// cannot use position of top card, in case it's being dragged
	x0, y0, x1, y1 = p.BaizeRect()
	if p.Fan == "" || p.Fan == "None" {
		return
	}
	if p.CardCount() > 1 {
		var x, y int
		if p.Tail == nil {
			x, y = p.Peek().BaizePosition()
		} else {
			// do not include cards being dragged, stop before Tail[0]
			for i := 1; i < p.CardCount(); i++ {
				if p.Cards[i] == p.Tail[0] {
					if i > 0 {
						x, y = p.Cards[i-1].BaizePosition()
					}
					break
				}
			}
		}
		switch p.Fan {
		case "Right", "Waste":
			x1 = x + CardWidth
		case "Down", "WasteDown":
			y1 = y + CardHeight
		}
	}
	return // using named return parameters
}

// FannedScreenRect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) FannedScreenRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0, x1, y1 = p.FannedBaizeRect()
	x0 += TheBaize.DragOffsetX
	x1 += TheBaize.DragOffsetX
	y0 += TheBaize.DragOffsetY
	y1 += TheBaize.DragOffsetY
	return // using named return parameters
}

// Hidden returns true if this Pile is off screen
func (p *Pile) Hidden() bool {
	return p.X < 0 || p.Y < 0
}

// SetAccept updates the Accept for this pile and updates the background image
func (p *Pile) SetAccept(ord int) {
	if p.localAccept != ord {
		p.localAccept = ord
		p.CreateBackgroundImage()
	}
}

// SetRecycles updates the Recycles for this pile and updates the background image
func (p *Pile) SetRecycles(n int) {
	if p.localRecycles != n {
		p.localRecycles = n
		p.CreateBackgroundImage()
	}
}

// CardCount returns the number of cards in this Pile
func (p *Pile) CardCount() int {
	return len(p.Cards)
}

// Peek topmost Card of this Pile (a stack)
func (p *Pile) Peek() *Card {
	if p.CardCount() == 0 {
		return nil
	}
	return p.Cards[p.CardCount()-1]
}

// Pop a Card off the end of this Pile (a stack)
func (p *Pile) Pop() *Card {
	if p.CardCount() == 0 {
		return nil
	}
	c := p.Cards[p.CardCount()-1]
	p.Cards = p.Cards[:p.CardCount()-1]
	c.owner = nil
	c.FlipUp()

	// experimental turn over exposed card here
	// if len(p.cards) > 0 {
	// 	p.cards[len(p.cards)-1].FlipUp()
	// }

	return c
}

// Push a Card onto the end of this Pile (a stack)
func (p *Pile) Push(c *Card) {
	c.StopSpinning()
	if strings.HasPrefix(p.Class, "Stock") {
		c.FlipDown()
	}
	c.owner = p
	c.TransitionTo(p.PushedFannedPosition()) // do this BEFORE appending card to pile
	p.Cards = append(p.Cards, c)
}

// RepushAllCards takes all the cards off the pile and puts them back
func (p *Pile) RepushAllCards() {
	// because we're about to use copy(), tmp must have a length
	var tmp = make([]*Card, len(p.Cards), cap(p.Cards)) // https://github.com/golang/go/wiki/SliceTricks#copy
	// len(tmp) == len(p.Cards)
	copy(tmp, p.Cards)
	p.Cards = p.Cards[:0] // keep the underlying array, slice the slice to zero length
	for _, c := range tmp {
		p.Push(c)
	}
}

// Extract a Card from the Pile
func (p *Pile) Extract(idx int) *Card {
	c := p.Cards[idx]
	p.Cards = append(p.Cards[:idx], p.Cards[idx+1:]...)
	c.owner = nil
	c.FlipUp()
	return c
}

// CanAcceptCard returns true if this Pile can accept the Card
func (p *Pile) CanAcceptCard(c *Card) bool {
	return p.CanAcceptTail([]*Card{c}, true)
}

// PushedFannedPosition returns the x,y screen coords of a Card that will be pushed onto this Pile
func (p *Pile) PushedFannedPosition() (int, int) {
	x, y := p.BaizePosition()
	switch p.Fan {
	case "", "None":
		// do nothing
	case "Down":
		for _, c := range p.Cards {
			if c.Prone() {
				y = y + (CardHeight / backFanFactor * p.scrunchPercentage / 100)
			} else {
				y = y + (CardHeight / faceFanFactor * p.scrunchPercentage / 100)
			}
		}
	case "Right":
		for _, c := range p.Cards {
			if c.Prone() {
				x = x + (CardWidth / backFanFactor * p.scrunchPercentage / 100)
			} else {
				x = x + (CardHeight / faceFanFactor * p.scrunchPercentage / 100)
			}
		}
	case "Waste":
		x0, y0 := p.BaizePosition()
		x1 := x0 + CardWidth/faceFanFactor
		x2 := x1 + CardWidth/faceFanFactor
		switch p.CardCount() {
		case 0:
			// do nothing, incoming card will be at x,y
		case 1:
			// incoming card will be at slot [1]
			x = x1
		case 2:
			// incoming card will be at slot [2]
			x = x2
		default: // >=3 cards
			// incoming card will be at slot [2]
			x = x2
			// top card needs to transition from slot[2] to slot[1]
			i := p.CardCount() - 1
			p.Cards[i].TransitionTo(x1, y0)
			// mid card needs to transition from slot[1] to slot[0]
			i--
			// p.Cards[i].TransitionTo(x0, y0) not needed will be done by loop below
			// most cards will be at pile x0,y0
			for ; i >= 0; i-- {
				p.Cards[i].TransitionTo(x0, y0)
			}
		}
	case "WasteDown":
		x0, y0 := p.BaizePosition()
		y1 := y0 + CardWidth/faceFanFactor
		y2 := y1 + CardWidth/faceFanFactor
		switch p.CardCount() {
		case 0:
			// do nothing, incoming card will be at x,y
		case 1:
			// incoming card will be at slot [1]
			y = y1
		case 2:
			// incoming card will be at slot [2]
			y = y2
		default: // >=3 cards
			// incoming card will be at slot [2]
			y = y2
			// top card needs to transition from slot[2] to slot[1]
			i := p.CardCount() - 1
			p.Cards[i].TransitionTo(x0, y1)
			// mid card needs to transition from slot[1] to slot[0]
			i--
			// p.Cards[i].TransitionTo(x0, y0) not needed will be done by loop below
			// most cards will be at pile x0,y0
			for ; i >= 0; i-- {
				p.Cards[i].TransitionTo(x0, y0)
			}
		}
	}
	return x, y
}

func (p *Pile) makeTail(c *Card) []*Card {
	var tail []*Card
	for i, pc := range p.Cards {
		if pc == c {
			tail = p.Cards[i:]
			break
		}
	}
	if len(tail) == 0 {
		log.Panic("Pile.makeTail made an empty tail")
	}
	return tail
}

// CanAcceptTail returns true if the Pile can accept the tail of Cards from another Pile
func (p *Pile) CanAcceptTail(Tail []*Card, canToast bool) bool {

	if len(Tail) == 0 { // len() for nil slices is defined as zero
		log.Panic("empty tail passed to CanAcceptTail")
	}

	c0 := Tail[0]

	if c0.owner == p {
		return false // Cannot drag cards to yourself
	}

	targetClass := c0.owner.GetStringAttribute("Target")
	if targetClass != "" {
		if targetClass != p.Class {
			if canToast {
				TheBaize.ui.Toast("Cards from " + c0.owner.Class + " can only be dragged to " + targetClass + " not to " + p.Class)
			}
			return false
		}
	}

	switch p.Class {
	case "Waste":
		if len(Tail) == 1 && c0.owner.Class == "Stock" {
			// if ctm := c0.owner.GetStringAttribute("CardsToMove"); ctm == "" || ctm == "1" {
			return true
			// }
		}
		// if canToast {
		// 	TheBaize.ui.Toast("Can only drag one card from Stock to Waste")
		// }
		return false

	case "Foundation":
		if p.CardCount() == 13 {
			// with rank wrap, may get >13 cards in a foundation
			// if canToast {
			// 	TheBaize.ui.Toast("Can only have 13 cards in a Foundation")
			// }
			return false
		}
		// Duchess rule
		if afp := p.GetStringAttribute("AcceptFirstPush"); afp != "" {
			if p.localAccept == 0 {
				if c0.owner.Class != afp {
					if canToast {
						TheBaize.ui.Toast(fmt.Sprintf("%s can only accept first card from a %s", p.Class, afp))
					}
					return false
				}
			}
		}
		// Spider only accepts a tail of 13 cards, and onto an empty Foundation
		if p.Flags&BuildFlagSpider == BuildFlagSpider {
			if len(Tail) != 13 {
				if canToast {
					TheBaize.ui.Toast("You can only drag 13 cards to a Foundation")
				}
				return false
			}
			if p.CardCount() > 0 {
				if canToast {
					TheBaize.ui.Toast("The Foundation must be empty")
				}
				return false
			}
			return isTailConformant(p.Build, p.Flags, Tail)
		} else {
			if len(Tail) != 1 {
				if canToast {
					TheBaize.ui.Toast("CYou can only drag one card to a Foundation")
				}
				return false
			}
			if p.CardCount() == 0 {
				if p.localAccept > 0 {
					return c0.Ordinal() == p.localAccept
				}
				return true
			}
			return isCardPairConformant(p.Build, p.Flags, p.Peek(), c0)
		}

	case "Tableau":
		if p.Flags&DragFlagSingle == DragFlagSingle {
			if TheUserData.PowerMoves {
				pm := powerMoves(TheBaize.Piles, p)
				if len(Tail) > pm {
					if canToast {
						TheBaize.ui.Toast(fmt.Sprintf("Enough free space to move %s, not %d", util.Pluralize("card", pm), len(Tail)))
					}
					return false
				}
			} else {
				if len(Tail) > 1 {
					if canToast {
						TheBaize.ui.Toast("You can only drag a single card")
					}
					return false
				}
			}
		}
		if p.Flags&DragFlagSingleOrPile == DragFlagSingleOrPile {
			if !(len(Tail) == 1 || len(Tail) == c0.owner.CardCount()) {
				if canToast {
					TheBaize.ui.Toast("You can only drag a single card or the whole pile")
				}
				return false
			}
		}
		if p.CardCount() == 0 {
			if afAttrib := p.GetStringAttribute("AcceptFrom"); afAttrib != "" {
				afList := strings.Split(afAttrib, ",")
				for _, class := range afList {
					if c0.owner.Class == class {
						return true
					}
				}
				if canToast {
					TheBaize.ui.Toast(fmt.Sprintf("%s can only accept cards from %s", p.Class, afAttrib))
				}
				return false
			}
			if p.localAccept > 0 {
				return c0.Ordinal() == p.localAccept
			}
			return true
		}
		return isCardPairConformant(p.Build, p.Flags, p.Peek(), c0)

	case "Cell":
		ok := len(Tail) == 1 && p.CardCount() == 0
		if !ok && canToast {
			TheBaize.ui.Toast("You can only have one card in a Cell")
		}
		return ok
	case "Reserve":
		if canToast {
			TheBaize.ui.Toast("You cannot move a card to a Reserve")
		}
		return false
	}
	if canToast {
		TheBaize.ui.Toast("You cannot move a card there")
	}
	return false // Reserve, Stock, StockSpider, StockScorpion
}

// StartDrag this card and all the others after it in the stack
func (p *Pile) StartDrag(c *Card) bool {

	// no need for this with Foundation Drag=0
	// if strings.HasPrefix(c.owner.Class, "Foundation") {
	// 	return false // cannot take cards off foundation
	// }
	if c.Spinning() {
		c.Flip()
		return false
	}
	if c.Transitioning() || c.Flipping() {
		println("unwise to drag an animating or flipping card")
		return false
	}

	p.Tail = p.makeTail(c)

	if p.Flags&DragFlagSingle == DragFlagSingle {
		if TheUserData.PowerMoves && p.Class == "Tableau" {
			pm := powerMoves(TheBaize.Piles, p)
			if len(p.Tail) > pm {
				TheBaize.ui.Toast(fmt.Sprintf("Enough free space to move %s, not %d", util.Pluralize("card", pm), len(p.Tail)))
				p.ApplyToTail((*Card).Shake)
				p.Tail = nil
				return false
			}
		} else {
			if len(p.Tail) > 1 {
				TheBaize.ui.Toast(p.Class + " can only drag a single card")
				p.ApplyToTail((*Card).Shake)
				p.Tail = nil
				return false
			}
		}
	}
	if p.Flags&DragFlagSingleOrPile == DragFlagSingleOrPile {
		if !(len(p.Tail) == 1 || len(p.Tail) == len(p.Cards)) {
			TheBaize.ui.Toast("You can only drag a single card or the entire pile")
			p.ApplyToTail((*Card).Shake)
			p.Tail = nil
			return false
		}
	}
	if !isTailConformant(p.Drag, p.Flags, p.Tail) {
		p.ApplyToTail((*Card).Shake)
		p.Tail = nil
		return false
	}
	p.ApplyToTail((*Card).StartDrag)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	return true
}

// StopDrag this card and all the others after it in the stack
func (p *Pile) StopDrag(c *Card) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	p.ApplyToTail((*Card).StopDrag)
	p.Tail = nil
}

// CancelDrag this card and all the others after it in the stack
func (p *Pile) CancelDrag(c *Card) {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	p.ApplyToTail((*Card).CancelDrag)
	p.Tail = nil
}

// ApplyToTail applies a method func to this card and all the others after it in the tail
func (p *Pile) ApplyToTail(fn func(*Card)) {
	// https://golang.org/ref/spec#Method_expressions
	// (*Card).CancelDrag yields a function with the signature func(*Card)
	// fn passed as a method expression so add the receiver explicitly
	for _, tc := range p.Tail {
		fn(tc)
	}
}

// ApplyToCards applys a function to each card in the pile
// caller must use a method expression, eg (*Card).StartSpinning, yielding a function value
// with a regular first parameter taking the place of the receiver
func (p *Pile) ApplyToCards(fn func(*Card)) {
	for _, c := range p.Cards {
		fn(c)
	}
}

// ApplyToCards2 applies a method func to this card and all the others after it in the stack
// func (p *Pile) ApplyToCards2(fn func(*Card, int, int), dx, dy int) {
// 	for _, c := range p.Cards {
// 	  fn(c, dx, dy)
// 	}
// }

// DragTailBy repositions all the cards in the tail
func (p *Pile) DragTailBy(dx, dy int) {
	for _, tc := range p.Tail {
		tc.DragBy(dx, dy)
	}
}

// Complete returns true if this Pile is complete
func (p *Pile) Complete() bool {
	// a game is complete when all piles except foundations are empty

	// cw, ok := p.GetIntAttribute("CompleteWhen")
	// if ok {
	// 	return p.CardCount() == cw
	// }

	if !strings.HasPrefix(p.Class, "Foundation") {
		return p.CardCount() == 0
	}

	return true
}

// Conformant returns true if all cards in this pile are conformant
func (p *Pile) Conformant() bool {
	if len(p.Cards) == 0 {
		return true
	}
	if strings.HasPrefix(p.Class, "Stock") || p.Class == "Waste" {
		return false
	}
	// that leaves Cell, Reserve, Tableau
	return isTailConformant(p.Build, p.Flags, p.Cards)
}

// BuryCards moves cards with the specified ordinal to the bottom of the stack
func (p *Pile) BuryCards(ordinal int) {
	tmp := make([]*Card, 0, cap(p.Cards))
	for _, c := range p.Cards {
		if c.Ordinal() == ordinal {
			tmp = append(tmp, c)
		}
	}
	for _, c := range p.Cards {
		if c.Ordinal() != ordinal {
			tmp = append(tmp, c)
		}
	}
	p.Cards = p.Cards[:0] // keep the underlying array, slice the slice to zero length
	for i := 0; i < len(tmp); i++ {
		p.Push(tmp[i])
	}
}

// DisinterCards moves cards with the specified ordinal to the top of the stack
func (p *Pile) DisinterCards(ordinal int) {
	tmp := make([]*Card, 0, cap(p.Cards))
	for _, c := range p.Cards {
		if c.Ordinal() != ordinal {
			tmp = append(tmp, c)
		}
	}
	for _, c := range p.Cards {
		if c.Ordinal() == ordinal {
			tmp = append(tmp, c)
		}
	}
	p.Cards = p.Cards[:0] // keep the underlying array, slice the slice to zero length
	for i := 0; i < len(tmp); i++ {
		p.Push(tmp[i])
	}
}

// Layout the cards in this Pile
// func (p *Pile) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return outsideWidth, outsideHeight
// }

// Update the Pile state (transitions, user input)
func (p *Pile) Update() error {
	for _, c := range p.Cards {
		c.Update()
	}
	return nil
}

// Draw renders the Pile into the screen
func (p *Pile) Draw(screen *ebiten.Image) {
	if p.backgroundImage != nil {
		op := &ebiten.DrawImageOptions{}
		x, y := p.ScreenPosition()
		op.GeoM.Translate(float64(x), float64(y))
		// the following makes dragging a bit laggy, even with no cards
		// if x, y := ebiten.CursorPosition(); util.InRect(x, y, p.ScreenRect) {
		// 	op.ColorM.Scale(1, 1, 1, 0.5)
		// 	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && strings.HasPrefix(p.Class, "Stock") {
		// 		op.GeoM.Translate(2, 2)
		// 	}
		// }
		screen.DrawImage(p.backgroundImage, op)
	}

	if DebugMode {
		if p.scrunchSize == 0 || p.CardCount() < 4 || p.Class != "Tableau" {
			return
		}
		var maxWidth, maxHeight int
		x0, y0 := p.BaizePosition()
		switch p.Fan {
		case "", "None", "Waste", "WasteDown":
			return
		case "Down":
			maxHeight = p.scrunchSize * CardHeight
			ebitenutil.DrawRect(screen, float64(x0+TheBaize.DragOffsetX), float64(y0+TheBaize.DragOffsetY), float64(CardWidth), float64(maxHeight), color.RGBA{0, 0, 0, 0x10})
		case "Right":
			maxWidth = p.scrunchSize * CardWidth
			ebitenutil.DrawRect(screen, float64(x0+TheBaize.DragOffsetX), float64(y0+TheBaize.DragOffsetY), float64(maxWidth), float64(CardHeight), color.RGBA{0, 0, 0, 0x10})
		}
	}
}

// DrawStaticCards renders the Cards in the Pile into the screen
func (p *Pile) DrawStaticCards(screen *ebiten.Image) {
	// draw dragging/lerping cards last so they appear on top
	for _, c := range p.Cards {
		if !c.Transitioning() && !c.Flipping() && !c.Dragging() {
			c.Draw(screen)
		}
	}
}

// DrawTransitioningCards renders the Cards in the Pile into the screen
func (p *Pile) DrawTransitioningCards(screen *ebiten.Image) {
	for _, c := range p.Cards {
		if c.Transitioning() || c.Dragging() {
			c.Draw(screen)
		}
	}
}

// DrawFlippingCards renders the Cards in the Pile into the screen
func (p *Pile) DrawFlippingCards(screen *ebiten.Image) {
	for _, c := range p.Cards {
		if c.Flipping() {
			c.Draw(screen)
		}
	}
}
