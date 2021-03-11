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
	Class             string
	X, Y              int
	Fan               string
	localAccept       int
	localRecycles     int
	Attributes        map[string]string
	Cards             []*Card
	Tail              []*Card
	buildRules        int
	buildFlags        int
	dragRules         int
	dragFlags         int
	scrunchPercentage int
	backgroundImage   *ebiten.Image
}

// NewPile create and fills in a Pile object
func NewPile(class string, x, y int, fan string, attribs map[string]string) *Pile {
	p := &Pile{Class: class, X: x, Y: y, Fan: fan, Attributes: attribs}

	br, ok := p.GetIntAttribute("Build")
	if !ok {
		log.Fatal("no Build rules for Pile " + p.Class)
	}
	p.buildRules = br % 100
	p.buildFlags = br / 100 // 1=rank wrap 2=power moves

	d, ok := p.GetIntAttribute("Drag")
	if !ok {
		log.Fatal("no Drag attribute for Pile " + p.Class)
	}
	p.dragRules = d % 100
	p.dragFlags = d / 100 // 1=single card only (no tail)

	p.Reset()
	return p
}

// Reset the pile
func (p *Pile) Reset() {
	p.localAccept, _ = p.GetIntAttribute("Accept")
	p.localRecycles, _ = p.GetIntAttribute("Recycles")
	p.scrunchPercentage = 100
	p.createBackgroundImage()
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
		return str == "true" || str == "True" || str == "T" || str == "1"
	}
	return false
}

func (p *Pile) createBackgroundImage() {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(colorPile)
	dc.SetLineWidth(4)
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Stroke()

	if p.localAccept > 0 && p.localAccept <= 13 {
		dc.SetFontFace(schriftbank.CardOrdinal)
		dc.DrawStringAnchored(util.OrdinalToShortString(p.localAccept), float64(CardWidth)/3.333, float64(CardHeight)/6.666, 0.5, 0.5)
		dc.SetLineWidth(1)
		// dc.DrawLine(0, float64(CardHeight)/6.666, float64(CardWidth), float64(CardHeight)/6.666)
		// dc.DrawLine(float64(CardWidth)/3.333, 0, float64(CardWidth)/3.333, float64(CardHeight))
		dc.Stroke()
	}
	if strings.HasPrefix(p.Class, "Stock") {
		dc.SetFontFace(schriftbank.CardSymbolLarge)
		if p.localRecycles == 0 {
			// anything put here either doesn't render (0x1F6AB) or looks ugly
			// dc.SetColor(BasicColors["Red"])
			// dc.SetFontFace(TheCardFonts.large)
			// dc.DrawStringAnchored("O", float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.5)
			// dc.DrawStringAnchored("/", float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.5)
		} else {
			dc.DrawStringAnchored(string(rune(0x2672)), float64(CardWidth)/2, float64(CardHeight)/2, 0.5, 0.5)
		}
		dc.Stroke()
	}
	p.backgroundImage = ebiten.NewImageFromImage(dc.Image())
}

// Position returns the x,y screen coords of this pile
func (p *Pile) Position() (int, int) {
	return (p.X * PileMarginX) + (p.X * CardWidth), TopMargin + (p.Y * PileMarginY) + (p.Y * CardHeight)
}

// Rect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0 = p.Position()
	x1 = x0 + CardWidth
	y1 = y0 + CardHeight
	return // using named return parameters
}

// FannedRect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) FannedRect() (x0 int, y0 int, x1 int, y1 int) {
	// cannot use position of top card, in case it's being dragged
	x0, y0, x1, y1 = p.Rect()
	if p.CardCount() > 1 {
		var x, y int
		if p.Tail == nil {
			x, y = p.Peek().Position()
		} else {
			// x, y = p.PushedFannedPosition() // this fudge is an approximation
			// do not include cards being dragged, stop before Tail[0]
			for i := 1; i < p.CardCount(); i++ {
				if p.Cards[i] == p.Tail[0] {
					if i > 0 {
						x, y = p.Cards[i-1].Position()
					}
					break
				}
			}
		}
		switch p.Fan {
		case "", "None":
			// do nothing
		case "Right":
			x1 = x + CardWidth
		case "Down":
			y1 = y + CardHeight
		}
	}
	return // using named return parameters
}

// SetAccept updates the Accept for this pile and updates the background image
func (p *Pile) SetAccept(ord int) {
	p.localAccept = ord
	p.createBackgroundImage()
}

// SetRecycles updates the Recycles for this pile and updates the background image
func (p *Pile) SetRecycles(n int) {
	p.localRecycles = n
	p.createBackgroundImage()
}

// CardCount returns the number of cards in this Pile
func (p *Pile) CardCount() int {
	return len(p.Cards)
}

// Peek topmost Card of this Pile (a stack)
func (p *Pile) Peek() *Card {
	if 0 == p.CardCount() {
		return nil
	}
	return p.Cards[p.CardCount()-1]
}

// Pop a Card off the end of this Pile (a stack)
func (p *Pile) Pop() *Card {
	if 0 == p.CardCount() {
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
	if strings.HasPrefix(p.Class, "Stock") {
		c.FlipDown()
	}
	c.owner = p
	c.TransitionTo(p.PushedFannedPosition()) // do this BEFORE appending card to pile
	p.Cards = append(p.Cards, c)
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

	switch p.Class {
	case "Stock":
		return false // user cannot move card to stock
	case "Waste":
		return c.owner.Class == "Stock" // user can only move card to waste from stock
	case "Foundation":
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c.Ordinal() == p.localAccept
			}
			return true
		}
		return isConformant0(p.buildRules, p.buildFlags, p.Peek(), c)
	case "Tableau":
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c.Ordinal() == p.localAccept
			}
			return true
		}
		return isConformant0(p.buildRules, p.buildFlags, p.Peek(), c)
	case "Cell":
		return p.CardCount() == 0
	}
	return false
}

// CanAcceptTail returns true if this Pile can accept the tail of Cards from another Pile
func (p *Pile) CanAcceptTail(piles []*Pile, Tail []*Card) bool {

	if Tail == nil || len(Tail) == 0 {
		log.Fatal("empty tail passed to CanAcceptTail")
	}

	c0 := Tail[0]

	if c0.owner == p {
		println("cannot drag cards to yourself")
		return false
	}

	targetClass := c0.owner.GetStringAttribute("Target")
	if targetClass != "" {
		if targetClass != p.Class {
			TheBaize.ui.Toast("Cards from " + c0.owner.Class + " can only be dragged to " + targetClass + " not to " + p.Class)
			return false
		}
	}

	switch p.Class {
	case "Stock":
		return false // user cannot drag cards to stock

	case "Waste":
		if c0.owner.Class == "Stock" { // user can drag a single card from stock to waste
			ctm := c0.owner.GetStringAttribute("CardsToMove")
			if ctm == "" || ctm == "1" {
				return true
			}
		}
		return false

	case "FoundationSpider":
		if len(Tail) != 13 {
			return false
		}
		if p.CardCount() > 0 {
			return false
		}
		return isConformant(p.buildRules, p.buildFlags, Tail)

	case "Foundation":
		if len(Tail) != 1 {
			return false
		}
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c0.Ordinal() == p.localAccept
			}
			return true
		}
		return isConformant0(p.buildRules, p.buildFlags, p.Peek(), c0)

	case "Tableau":
		if p.buildFlags&2 == 2 {
			pm := powerMoves(piles, p)
			if len(Tail) > pm {
				TheBaize.ui.Toast(fmt.Sprintf("Not enough free space to drag %d cards", len(Tail)))
				return false
			}
			println("can drag", len(Tail), "cards")
		}
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c0.Ordinal() == p.localAccept
			}
			return true
		}
		return isConformant0(p.buildRules, p.buildFlags, p.Peek(), c0)

	case "Cell":
		return len(Tail) == 1 && p.CardCount() == 0
	}
	return false
}

// PushedFannedPosition returns the x,y screen coords of a Card that will be pushed onto this Pile
func (p *Pile) PushedFannedPosition() (int, int) {
	x, y := p.Position()
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
		x0, y0 := p.Position()
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
	}
	return x, y
}

func (p *Pile) makeTail(c *Card) []*Card {
	var tail []*Card // append works on a nil slice, yay
	marking := false
	for _, pc := range p.Cards {
		if !marking && pc == c {
			marking = true
		}
		if marking {
			tail = append(tail, pc)
		}
	}
	return tail
}

// StartDrag this card and all the others after it in the stack
func (p *Pile) StartDrag(c *Card) bool {

	// no need for this with Foundation Drag=0
	// if strings.HasPrefix(c.owner.Class, "Foundation") {
	// 	return false // cannot take cards off foundation
	// }
	if c.Animating() {
		println("unwise to drag an animating card")
		return false
	}

	p.Tail = p.makeTail(c)

	if p.dragFlags&1 == 1 && len(p.Tail) > 1 {
		TheBaize.ui.Toast(p.Class + " can only drag a single card")
		p.ApplyToTail((*Card).Shake)
		p.Tail = nil
		return false
	}
	if !isConformant(p.dragRules, p.dragFlags, p.Tail) {
		println("non-conformant drag")
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

// ApplyToTail applies a method func to this card and all the others after it in the stack
func (p *Pile) ApplyToTail(fn func(*Card)) {
	// https://golang.org/ref/spec#Method_expressions
	// (*Card).CancelDrag yields a function with the signature func(*Card)
	// fn passed as a method expression so add the receiver explicitly
	for _, tc := range p.Tail {
		fn(tc)
	}
}

// DragTailBy repositions all the cards in the tail (from c inclusive)
func (p *Pile) DragTailBy(dx, dy int) {
	// would have used https://golang.org/ref/spec#Method_expressions
	// but couldn't figure out the syntax
	// so using a standalone loop instead
	for _, tc := range p.Tail {
		tc.DragBy(dx, dy)
	}
}

// Complete returns true if this Pile is complete
func (p *Pile) Complete() bool {
	// a game is complete when all piles except foundations are empty

	cw, ok := p.GetIntAttribute("CompleteWhen")
	if ok {
		return p.CardCount() == cw
	}

	if !strings.HasPrefix(p.Class, "Foundation") {
		return p.CardCount() == 0
	}

	return true
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
func (p *Pile) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

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
		x, y := p.Position()
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(p.backgroundImage, op)
	}
	if DebugMode {
		// x1, y1, x2, y2 := p.FannedRect()
		// ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(x2-x1), float64(y2-y1), color.RGBA{0, 0, 0, 0x40})
		if p.CardCount() < 3 {
			return
		}
		s, ok := p.GetIntAttribute("Scrunch")
		if !ok {
			return
		}
		// var currWidth, currHeight int
		var maxWidth, maxHeight int
		switch p.Fan {
		case "", "None", "Waste":
			return
		case "Down":
			// currHeight = p.fannedHeight(p.scrunchPercentage)
			maxHeight = s * CardHeight
		case "Right":
			// currWidth = p.fannedWidth(p.scrunchPercentage)
			maxWidth = s * CardWidth
		}
		x0, y0 := p.Position()
		if maxWidth > 0 {
			ebitenutil.DrawRect(screen, float64(x0), float64(y0), float64(maxWidth), float64(CardHeight), color.RGBA{0, 0, 0, 0x10})
		}
		if maxHeight > 0 {
			ebitenutil.DrawRect(screen, float64(x0), float64(y0), float64(CardWidth), float64(maxHeight), color.RGBA{0, 0, 0, 0x10})
		}
		// if currWidth > 0 {
		// 	ebitenutil.DrawRect(screen, float64(x0), float64(y0), float64(currWidth), float64(CardHeight), color.RGBA{0, 0, 0, 0x40})
		// }
		// if currHeight > 0 {
		// 	ebitenutil.DrawRect(screen, float64(x0), float64(y0), float64(CardWidth), float64(currHeight), color.RGBA{0, 0, 0, 0x40})
		// }
	}
}

// DrawCards renders the Cards in the Pile into the screen
func (p *Pile) DrawCards(screen *ebiten.Image) {
	// draw dragging/lerping cards last so they appear on top
	for _, c := range p.Cards {
		if !c.Animating() {
			c.Draw(screen)
		}
	}
}

// DrawAnimatingCards renders the Cards in the Pile into the screen
func (p *Pile) DrawAnimatingCards(screen *ebiten.Image) {
	for _, c := range p.Cards {
		if c.Animating() {
			// ebitenutil.DebugPrint(screen, fmt.Sprintf("dragging card %s %d,%d", c.ID.String(), c.screenX, c.screenY))
			c.Draw(screen)
		}
	}
}
