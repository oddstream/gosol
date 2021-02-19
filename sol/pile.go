package sol

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	marginX int = 10
	marginY int = 10
)

var outline *ebiten.Image

func init() {
	dc := gg.NewContext(71, 96)
	dc.SetColor(colorPile)
	dc.SetLineWidth(4)
	dc.DrawRoundedRectangle(0, 0, float64(71), float64(96), 4)
	dc.Stroke()
	outline = ebiten.NewImageFromImage(dc.Image())
}

// CardOwner is an interface to objects that can own cards (Pile and Pile 'subclasses')
type CardOwner interface {
	New(map[string]string)
	Cards() []*Card
	Class() string
	Fan() string
	Position() (int, int)
	Peek() *Card
	Pop() *Card
	Push(*Card)
	Update() error
	Layout(int, int) (int, int)
	Draw(*ebiten.Image)
	DrawCards(*ebiten.Image)
}

// Pile is a generic container for cards
type Pile struct {
	cards []*Card
	x, y  int    // grid position of Pile
	fan   string // "None", "Down", "Right"
}

// Cards returns the slice of *Card
func (p *Pile) Cards() []*Card {
	return p.cards
}

// Fan returns the Fan of *Card
func (p *Pile) Fan() string {
	return p.fan
}

// Position returns the x,y screen coords of this pile
func (p *Pile) Position() (int, int) {
	return (p.x * marginX) + (p.x * 71), (p.y * marginY) + (p.y * 96)
}

// ToFront moves the Card to the top of the Pile (a stack)
func (p *Pile) ToFront(c *Card) {

}

// Peek topmost Card  of this Pile (a stack)
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
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (p *Pile) Push(c *Card) {
	p.cards = append(p.cards, c)
}

// Layout the cards in this Pile
func (p *Pile) Layout(outsideWidth, outsideHeight int) (int, int) {
	// stop if we meet a card that's transitioning
	switch p.fan {
	case "", "None":
		// do nothing
	case "Down":
		// TODO
	case "Right":
		// TODO
	}
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
	op := &ebiten.DrawImageOptions{}
	x, y := p.Position()
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(outline, op)
}

// DrawCards renders the Cards in the Pile into the screen
func (p *Pile) DrawCards(screen *ebiten.Image) {
	for _, c := range p.cards {
		c.Draw(screen)
	}
}
