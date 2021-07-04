package sol

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/sound"
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
	driver            Driver
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

	{
		if fn, ok := Class2NewDriver[pi.Class]; !ok {
			log.Fatal("No NewDriver factory func for Class", pi.Class)
		} else {
			p.driver = fn(p)
		}
	}

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

// Spider returns true if this pile is Spideresque
func (p *Pile) Spider() bool {
	return p.Flags&BuildFlagSpider == BuildFlagSpider
}

// Empty returns true if this pile is empty
func (p *Pile) Empty() bool {
	return len(p.Cards) == 0
}

func (p *Pile) CreateBackgroundImage() {

	invisible := p.GetBoolAttribute("Invisible")
	if invisible || p.localAccept == 99 || p.Class == "Reserve" {
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
		case "Right", "Waste", "WasteRight":
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
	if p.Empty() {
		return nil
	}
	return p.Cards[p.CardCount()-1]
}

// Pop a Card off the end of this Pile (a stack)
func (p *Pile) Pop() *Card {
	if p.Empty() {
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
// func (p *Pile) CanAcceptCard(c *Card) bool {
// 	ok, err := p.driver.CanAcceptTail([]*Card{c})
// 	if err != nil {
// 		TheBaize.ui.Toast(err.Error())
// 	}
// 	return ok
// }

// PushedFannedPosition returns the x,y screen coords of a Card that will be pushed onto this Pile
func (p *Pile) PushedFannedPosition() (int, int) {
	x, y := p.BaizePosition()
	switch p.Fan {
	case "", "None":
		// do nothing
	case "Down":
		backDelta := CardHeight / backFanFactor * p.scrunchPercentage / 100
		faceDelta := CardHeight / faceFanFactor * p.scrunchPercentage / 100
		for _, c := range p.Cards {
			if c.Prone() {
				y = y + backDelta
			} else {
				y = y + faceDelta
			}
		}
	case "Right":
		backDelta := CardWidth / backFanFactor * p.scrunchPercentage / 100
		faceDelta := CardWidth / faceFanFactor * p.scrunchPercentage / 100
		for _, c := range p.Cards {
			if c.Prone() {
				x = x + backDelta
			} else {
				x = x + faceDelta
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
	if !p.Empty() {
		for i, pc := range p.Cards {
			if pc == c {
				tail = p.Cards[i:]
				break
			}
		}
	}
	if len(tail) == 0 {
		log.Panic("Pile.makeTail made an empty tail")
	}
	return tail
}

func (p *Pile) indexOf(card *Card) (int, error) {
	for i, c := range p.Cards {
		if c == card {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%s not found in %s", card.String(), p.Class)
}

func (p *Pile) MakeConformantTail(c *Card) []*Card {
	var tail []*Card

	if p != c.owner {
		log.Panic("Incorrect call to Pile.MakeConformantTail")
	}

	switch len(p.Cards) {
	case 0:
		log.Panic("Pile.MakeConformantTail called on a empty pile(!?)")
	case 1:
		tail = []*Card{c}
	default:
		idx, err := p.indexOf(c)
		if err != nil {
			log.Panic(err.Error())
		}
		tail = p.Cards[idx:]
		if !isTailConformant(p.Build, p.Flags, tail) {
			tail = nil
		}
	}
	return tail
}

func (p *Pile) CountSortedAndUnsorted(sorted, unsorted int) (int, int) {
	if strings.HasPrefix(p.Class, "Foundation") {
		sorted += len(p.Cards)
	} else {
		for i := 0; i < len(p.Cards)-1; i++ {
			c1 := p.Cards[i]
			c2 := p.Cards[i+1]
			if isCardPairConformant(p.Build, p.Flags, c1, c2) {
				sorted++
			} else {
				unsorted++
			}
		}
	}
	return sorted, unsorted
}

// MoveCards from one pile to another, always from card downwards (inclusive)
func (dst *Pile) MoveCards(c *Card) {

	src := c.owner
	oldSrcLen := len(src.Cards)

	// find the index of the first card we will move
	moveFrom, err := src.indexOf(c)
	if err != nil {
		log.Panic(err.Error())
	}

	tmp := make([]*Card, 0, cap(src.Cards))

	// pop the tail off the source and push onto temp stack
	for i := len(src.Cards) - 1; i >= moveFrom; i-- {
		tmp = append(tmp, src.Pop())
	}

	// pop all cards off the temp stack and onto the destination
	for len(tmp) > 0 {
		dc := tmp[len(tmp)-1]
		tmp = tmp[:len(tmp)-1]
		dst.Push(dc)
	}

	if oldSrcLen == len(src.Cards) {
		log.Println("nothing happened in MoveCards")
		return
	}

	switch dst.Class {
	case "Foundation", "FoundationSpider", "Waste", "Golf":
		sound.Play("Slide")
	default:
		sound.Play("Place")
	}

	// flip up an exposed source card
	if !strings.HasPrefix(src.Class, "Stock") {
		if tc := src.Peek(); tc != nil {
			tc.FlipUp()
		}
	}
	// special case: waste may need refanning if we took a card from it
	if strings.HasPrefix(src.Fan, "Waste") {
		src.RepushAllCards()
	}
	src.ScrunchCards()
	dst.ScrunchCards()

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
		// println("unwise to drag an animating or flipping card")
		return false
	}

	p.Tail = p.makeTail(c)

	if p.Flags&DragFlagSingle == DragFlagSingle {
		if ThePreferences.PowerMoves && p.Class == "Tableau" {
			pm := powerMoves(TheBaize.Piles, p)
			if len(p.Tail) > pm {
				TheUI.Toast(fmt.Sprintf("Enough free space to move %s, not %d", util.Pluralize("card", pm), len(p.Tail)))
				p.ApplyToTail((*Card).Shake)
				p.Tail = nil
				return false
			}
		} else {
			if len(p.Tail) > 1 {
				TheUI.Toast(p.Class + " can only drag a single card")
				p.ApplyToTail((*Card).Shake)
				p.Tail = nil
				return false
			}
		}
	}
	if p.Flags&DragFlagSingleOrPile == DragFlagSingleOrPile {
		if !(len(p.Tail) == 1 || len(p.Tail) == len(p.Cards)) {
			TheUI.Toast("You can only drag a single card or the entire pile")
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

// BuryCards moves cards with the specified ordinal to the beginning of the pile
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

// DisinterCards moves cards with the specified ordinal to the end of the pile
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
		if p.scrunchSize == 0 || p.CardCount() < 4 {
			return
		}
		var maxWidth, maxHeight int
		x0, y0 := p.BaizePosition()
		switch p.Fan {
		case "Down":
			maxHeight = p.scrunchSize * CardHeight
			ebitenutil.DrawRect(screen, float64(x0+TheBaize.DragOffsetX), float64(y0+TheBaize.DragOffsetY), float64(CardWidth), float64(maxHeight), color.RGBA{0, 0, 0, 0x20})
		case "Right":
			maxWidth = p.scrunchSize * CardWidth
			ebitenutil.DrawRect(screen, float64(x0+TheBaize.DragOffsetX), float64(y0+TheBaize.DragOffsetY), float64(maxWidth), float64(CardHeight), color.RGBA{0, 0, 0, 0x20})
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
