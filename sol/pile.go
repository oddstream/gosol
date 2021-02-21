package sol

import (
	"reflect"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	marginX          int = 10
	marginY          int = 10
	proneStackFactor int = 5
	cardStackFactor  int = 3
)

// CardOwner is an interface to objects that can own cards (Pile and Pile 'subclasses')
type CardOwner interface {
	New(map[string]string)
	Cards() []*Card
	Class() string
	Deal() string
	Position() (int, int)
	Rect() (int, int, int, int)
	FannedRect() (int, int, int, int)
	StartDrag(*Card)
	CancelDrag(*Card)
	StopDrag(*Card)
	Peek() *Card
	Pop() *Card
	Push(*Card)
	CanAcceptCard(*Card) bool
	Update() error
	Layout(int, int) (int, int)
	Draw(*ebiten.Image)
	DrawCards(*ebiten.Image)
	DrawMovingCards(*ebiten.Image)
}

// Pile is a generic container for cards
type Pile struct {
	cards   []*Card
	x, y    int // grid position of Pile
	outline *ebiten.Image
	fan     string // "None", "Down", "Right"
	deal    string
	accept  int // ordinal of card to accept on empty pile, 0 == any (FoundationSpider has it's own rules)
}

// New fills in a blank Pile object to satify the CardOwner interface
func (p *Pile) New(map[string]string) {
}

func (p *Pile) createImage() {
	dc := gg.NewContext(71, 96)
	dc.SetColor(colorPile)
	dc.SetLineWidth(4)
	dc.DrawRoundedRectangle(0, 0, float64(71), float64(96), 4)
	if p.accept > 0 && p.accept <= 13 {
		var acceptChars = []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
		dc.SetFontFace(TheAcmeFonts.normal)
		dc.DrawString(acceptChars[p.accept], 71/7, 96/3)
	}
	dc.Stroke()
	p.outline = ebiten.NewImageFromImage(dc.Image())
}

// Cards returns the slice of *Card
func (p *Pile) Cards() []*Card {
	return p.cards
}

// Class returns the type of this Pile
func (p *Pile) Class() string {
	// can't use this generic Class() for all subtypes; they each need their own Class()
	return reflect.TypeOf(*p).Name() // .String() returns "sol.Stock"
}

// Deal returns the Deal of *Card
func (p *Pile) Deal() string {
	return p.deal
}

// Position returns the x,y screen coords of this pile
func (p *Pile) Position() (int, int) {
	return (p.x * marginX) + (p.x * 71), (p.y * marginY) + (p.y * 96)
}

// Rect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0 = p.Position()
	x1 = x0 + 71
	y1 = y0 + 96
	return // using named return parameters
}

// FannedRect gives the x,y screen coords of the pile's top left and bottom right corners
func (p *Pile) FannedRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0, x1, y1 = p.Rect()
	if len(p.cards) > 1 {
		x, y := p.Peek().Position()
		switch p.fan {
		case "", "None":
			// do nothing
		case "Right":
			x1 = x + 71
		case "Down":
			y1 = y + 96
		}
	}
	return // using named return parameters
}

// Peek topmost Card of this Pile (a stack)
func (p *Pile) Peek() *Card {
	if 0 == len(p.cards) {
		return nil
	}
	return p.cards[len(p.cards)-1]
}

// Pop a Card off the end of this Pile (a stack)
func (p *Pile) Pop() *Card {
	if 0 == len(p.cards) {
		return nil
	}
	c := p.cards[len(p.cards)-1]
	p.cards = p.cards[:len(p.cards)-1]
	c.owner = nil
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (p *Pile) Push(c *Card) {
	c.owner = p
	x, y := p.PushedFannedPosition()
	p.cards = append(p.cards, c)
	CTQ.Add(c, x, y)
}

// CanAcceptCard returns true if this Pile can accept the Card
func (p *Pile) CanAcceptCard(*Card) bool {
	return false
}

// PushedFannedPosition returns the x,y screen coords of a Card that will be pushed onto this Pile
func (p *Pile) PushedFannedPosition() (int, int) {
	x, y := p.Position()
	switch p.fan {
	case "", "None":
		// do nothing
	case "Down":
		for _, c := range p.cards {
			if c.prone {
				y = y + 96/proneStackFactor
			} else {
				y = y + 96/cardStackFactor
			}
		}
	case "Right":
		for _, c := range p.cards {
			if c.prone {
				x = x + 71/proneStackFactor
			} else {
				x = x + 96/cardStackFactor
			}
		}
	case "Waste":
		x0, y0 := p.Position()
		x1 := x0 + 71/cardStackFactor
		x2 := x1 + 71/cardStackFactor
		switch len(p.cards) {
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
			c := p.cards[len(p.cards)-1]
			CTQ.Add(c, x1, y0)
			// card below that needs to transition from slot[1] to slot[0]
			c = p.cards[len(p.cards)-2]
			CTQ.Add(c, x0, y0)
			// all other cards will be at pile x,y
			for i := 0; i < len(p.cards)-2; i++ {
				c = p.cards[i]
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
// 				y = y + 96/proneStackFactor
// 			} else {
// 				y = y + 96/cardStackFactor
// 			}
// 		}
// 	case "right":
// 		for _, c := range p.cards {
// 			CTQ.Add(c, x, y)
// 			if c.prone {
// 				x = x + 71/proneStackFactor
// 			} else {
// 				x = x + 96/cardStackFactor
// 			}
// 		}
// 	case "waste":
// 		// TODO
// 	}
// }

// StartDrag this card and all the others after it in the stack
func (p *Pile) StartDrag(c *Card) {
	p.ApplyToTail(c, (*Card).StartDrag)
}

// StopDrag this card and all the others after it in the stack
func (p *Pile) StopDrag(c *Card) {
	p.ApplyToTail(c, (*Card).StopDrag)
}

// CancelDrag this card and all the others after it in the stack
func (p *Pile) CancelDrag(c *Card) {
	p.ApplyToTail(c, (*Card).CancelDrag)
}

// https://golang.org/ref/spec#Method_expressions
// (*Card).CancelDrag yields a function with the signature func(*Card)

// ApplyToTail applies a method func to this card and all the others after it in the stack
func (p *Pile) ApplyToTail(c *Card, fn func(*Card)) {
	marking := false
	for i := 0; i < len(p.cards); i++ {
		pci := p.cards[i]
		if !marking && pci == c {
			marking = true
		}
		if marking {
			fn(pci)
		}
	}
}

// Layout the cards in this Pile
func (p *Pile) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the Pile state (transitions, user input)
func (p *Pile) Update() error {
	for _, c := range p.cards {
		c.Update()
	}
	return nil
}

// Draw renders the Pile into the screen
func (p *Pile) Draw(screen *ebiten.Image) {
	if p.outline != nil {
		op := &ebiten.DrawImageOptions{}
		x, y := p.Position()
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(p.outline, op)
	}
}

// DrawCards renders the Cards in the Pile into the screen
func (p *Pile) DrawCards(screen *ebiten.Image) {
	// draw dragging/lerping cards last so they appear on top
	for _, c := range p.cards {
		if c.dragging == false && c.lerping == false {
			c.Draw(screen)
		}
	}
}

// DrawMovingCards renders the Cards in the Pile into the screen
func (p *Pile) DrawMovingCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		if c.dragging == true || c.lerping == true {
			// ebitenutil.DebugPrint(screen, fmt.Sprintf("dragging card %s %d,%d", c.id, c.screenX, c.screenY))
			c.Draw(screen)
		}
	}
}
