package sol

import (
	"log"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
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
	Class           string
	X, Y            int
	Fan             string
	localAccept     int
	localRecycles   int
	Attributes      map[string]string
	Cards           []*Card
	Tail            []*Card
	backgroundImage *ebiten.Image
}

// NewPile create and fills in a Pile object
func NewPile(class string, x, y int, fan string, attribs map[string]string) *Pile {
	p := &Pile{Class: class, X: x, Y: y, Fan: fan, Attributes: attribs}
	p.localAccept, _ = p.GetIntAttribute("Accept")
	p.localRecycles, _ = p.GetIntAttribute("Recycles")
	p.createBackgroundImage()
	return p
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
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 4)
	accept, ok := p.GetIntAttribute("Accept")
	if ok && accept > 0 && accept <= 13 {
		dc.SetFontFace(TheCardFonts.regular)
		dc.DrawString(util.OrdinalToChar(accept), float64(CardWidth)/7, float64(CardHeight)/3)
	}
	if p.localRecycles == 0 {
		// TODO red no entry, but will appear on all piles :-/
	} else if p.localRecycles < 10 {
		// TODO green recycle glyph
	}
	dc.Stroke()
	p.backgroundImage = ebiten.NewImageFromImage(dc.Image())
}

// Position returns the x,y screen coords of this pile
func (p *Pile) Position() (int, int) {
	return (p.X * PileMarginX) + (p.X * CardWidth), (p.Y * PileMarginY) + (p.Y * CardHeight)
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
	x0, y0, x1, y1 = p.Rect()
	if p.CardCount() > 1 {
		x, y := p.Peek().Position()
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
	x, y := p.PushedFannedPosition()
	p.Cards = append(p.Cards, c)
	CTQ.Add(c, x, y)
}

// CanAcceptCard returns true if this Pile can accept the Card
func (p *Pile) CanAcceptCard(c *Card) bool {

	build, ok := p.GetIntAttribute("Build")
	if !ok {
		log.Fatal("no Build rules for Pile " + p.Class)
	}

	switch p.Class {
	case "Stock":
		return false // user cannot move card to stock
	case "Waste":
		return c.owner.Class == "Stock" // user can only move card to waste from stock
	case "Foundation":
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c.ordinal == p.localAccept
			}
			return true
		}
		return isConformant0(build, p.Peek(), c)
	case "Tableau":
		if p.CardCount() == 0 {
			if p.localAccept > 0 {
				return c.ordinal == p.localAccept
			}
			return true
		}
		return isConformant0(build, p.Peek(), c)
	case "Cell":
		return p.CardCount() == 0
	}
	return false
}

// CanAcceptTail returns true if this Pile can accept the tail of Cards from another Pile
func (p *Pile) CanAcceptTail(Tail []*Card) bool {

	if Tail == nil || len(Tail) == 0 {
		log.Fatal("CanAcceptTail with empty tail")
	}

	c0 := Tail[0]

	if c0.owner == p {
		println("Cannot drag cards to yourself")
		return false
	}

	targetClass := c0.owner.GetStringAttribute("Target")
	if targetClass != "" {
		if targetClass != p.Class {
			println("Cards from", c0.owner.Class, "can only be dragged to", targetClass, "not to", p.Class)
			return false
		}
	}

	accept, ok := p.GetIntAttribute("Accept")
	if !ok {
		accept = 0 // accept any card
	}
	buildRules, ok := p.GetIntAttribute("Build")
	if !ok {
		log.Fatal("No Build attribute for Pile " + p.Class)
	}

	switch p.Class {
	case "Stock":
		return false // user cannot drag cards to stock

	case "Waste":
		return c0.owner.Class == "Stock" // user can drag a card from stock to waste

	case "FoundationSpider":
		if len(Tail) != 13 {
			return false
		}
		if p.CardCount() > 0 {
			return false
		}
		return isConformant(buildRules, Tail)

	case "Foundation":
		if len(Tail) != 1 {
			return false
		}
		if p.CardCount() == 0 {
			if accept > 0 {
				return c0.ordinal == accept
			}
			return true
		}
		return isConformant0(buildRules, p.Peek(), c0)

	case "Tableau":
		if p.CardCount() == 0 {
			if accept > 0 {
				return c0.ordinal == accept
			}
			return true
		}
		return isConformant0(buildRules, p.Peek(), c0)

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
			if c.prone {
				y = y + CardHeight/backFanFactor
			} else {
				y = y + CardHeight/faceFanFactor
			}
		}
	case "Right":
		for _, c := range p.Cards {
			if c.prone {
				x = x + CardWidth/backFanFactor
			} else {
				x = x + CardHeight/faceFanFactor
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
			// card below needs to transition from slot[2] to slot[1]
			c := p.Cards[p.CardCount()-1]
			CTQ.Add(c, x1, y0)
			// card below that needs to transition from slot[1] to slot[0]
			c = p.Cards[p.CardCount()-2]
			CTQ.Add(c, x0, y0)
			// all other cards will be at pile x,y
			for i := 0; i < p.CardCount()-2; i++ {
				c = p.Cards[i]
				c.SetPosition(x0, y0)
			}
		}
	}
	return x, y
}

// Fan lays out the cards according to the Pile's fan attribute
// func (p *Pile) Fan() {
// 	x, y := p.Position()
// 	switch p.fan {
// 	case "", "none":
// 		for _, c := range p.cards {
// 			CTQ.Add(c, x, y)
// 		}
// 	case "down":
// 		for _, c := range p.cards {
// 			CTQ.Add(c, x, y)
// 			if c.prone {
// 				y = y + CardHeight/proneStackFactor
// 			} else {
// 				y = y + CardHeight/cardStackFactor
// 			}
// 		}
// 	case "right":
// 		for _, c := range p.cards {
// 			CTQ.Add(c, x, y)
// 			if c.prone {
// 				x = x + CardWidth/proneStackFactor
// 			} else {
// 				x = x + CardHeight/cardStackFactor
// 			}
// 		}
// 	case "waste":
// 		// TODO
// 	}
// }

// StartDrag this card and all the others after it in the stack
func (p *Pile) StartDrag(piles []*Pile, c *Card) bool {

	// no need for this with Foundation Drag=0
	// if strings.HasPrefix(c.owner.Class, "Foundation") {
	// 	return false // cannot take cards off foundation
	// }

	p.Tail = nil // append works on a nil slice, yay
	marking := false
	for i := 0; i < p.CardCount(); i++ {
		pci := p.Cards[i]
		if !marking && pci == c {
			marking = true
		}
		if marking {
			p.Tail = append(p.Tail, pci)
		}
	}
	d, ok := p.GetIntAttribute("Drag")
	if !ok {
		log.Fatal("No Drag attribute for Pile " + p.Class)
	}
	dragRules := d % 100
	dragFlags := d / 100 // 1=single card only (no tail) 2=powerMoves
	if dragFlags&1 == 1 && len(p.Tail) > 1 {
		println(p.Class, "can only drag a single card")
		p.ApplyToTail((*Card).Shake)
		p.Tail = nil
		return false
	}
	if dragFlags&2 == 2 {
		if len(p.Tail) <= powerMoves(piles, p) && isConformant(dragRules, p.Tail) {
			// that's ok
		} else {
			println("non-conformant powerMoves drag")
			p.ApplyToTail((*Card).Shake)
			p.Tail = nil
			return false
		}
	}
	if !isConformant(dragRules, p.Tail) {
		println("non-conformant drag")
		p.ApplyToTail((*Card).Shake)
		p.Tail = nil
		return false
	}
	p.ApplyToTail((*Card).StartDrag)
	return true
}

// StopDrag this card and all the others after it in the stack
func (p *Pile) StopDrag(c *Card) {
	p.ApplyToTail((*Card).StopDrag)
	p.Tail = nil
}

// CancelDrag this card and all the others after it in the stack
func (p *Pile) CancelDrag(c *Card) {
	p.ApplyToTail((*Card).CancelDrag)
	p.Tail = nil
}

// https://golang.org/ref/spec#Method_expressions
// (*Card).CancelDrag yields a function with the signature func(*Card)

// ApplyToTail applies a method func to this card and all the others after it in the stack
func (p *Pile) ApplyToTail(fn func(*Card)) {
	for _, tc := range p.Tail {
		fn(tc)
	}
	// marking := false
	// for i := 0; i < p.CardCount(); i++ {
	// 	pci := p.Cards[i]
	// 	if !marking && pci == c {
	// 		marking = true
	// 	}
	// 	if marking {
	// 		fn(pci)
	// 	}
	// }
}

// DragTailBy repositions all the cards in the tail (from c inclusive)
func (p *Pile) DragTailBy(dx, dy int) {
	// would have used https://golang.org/ref/spec#Method_expressions
	// but couldn't figure out the syntax
	// so using a standalone loop instead
	for _, tc := range p.Tail {
		tc.DragBy(dx, dy)
	}

	// marking := false
	// for i := 0; i < p.CardCount(); i++ {
	// 	ci := p.Cards[i]
	// 	if !marking && ci == c {
	// 		marking = true
	// 	}
	// 	if marking {
	// 		ci.DragBy(dx, dy)
	// 	}
	// }
}

// IsComplete returns true if this Pile is complete
func (p *Pile) IsComplete() bool {
	// a game is complete when all piles except foundations are empty

	cw, ok := p.GetIntAttribute("CompleteWhen")
	if ok {
		return p.CardCount() == cw
	}

	switch p.Class {
	case "Foundation":
	default:
		return p.CardCount() == 0
	}
	return true
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
			// ebitenutil.DebugPrint(screen, fmt.Sprintf("dragging card %s %d,%d", c.id, c.screenX, c.screenY))
			c.Draw(screen)
		}
	}
}
