package sol

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	Peek() *Card
	Pop() *Card
	Push(*Card)
	Fan()
	Update() error
	Layout(int, int) (int, int)
	Draw(*ebiten.Image)
	DrawCards(*ebiten.Image)
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

// Class returns the class of *Pile to satify the CardOwner interface
func (p *Pile) Class() string {
	return ""
}

// Deal returns the Deal of *Card
func (p *Pile) Deal() string {
	return p.deal
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
	c.owner = nil
	return c
}

// Push a Card onto the end of this Pile (a stack)
func (p *Pile) Push(c *Card) {
	c.owner = p
	p.cards = append(p.cards, c)
	x, y := p.Position()
	// c.TransitionTo(x, y)
	c.SetPosition(x, y)
}

// Fan lays out the cards according to the Pile's fan attribute
func (p *Pile) Fan() {
	// TODO stop if we meet a card that's transitioning (or flipping)?
	x, y := p.Position()
	switch p.fan {
	case "", "none":
		for _, c := range p.cards {
			if c.lerping {
				break
			}
			c.SetPosition(x, y)
		}
	case "down":
		for _, c := range p.cards {
			if c.lerping {
				break
			}
			c.SetPosition(x, y)
			if c.prone {
				y = y + 96/proneStackFactor
			} else {
				y = y + 96/cardStackFactor
			}
		}
	case "right":
		for _, c := range p.cards {
			if c.lerping {
				break
			}
			c.SetPosition(x, y)
			if c.prone {
				x = x + 71/proneStackFactor
			} else {
				x = x + 96/cardStackFactor
			}
		}
	case "waste":
		// TODO
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
		if c.dragging == false {
			c.Draw(screen)
		}
	}
	for _, c := range p.cards {
		if c.dragging == true {
			ebitenutil.DebugPrint(screen, fmt.Sprintf("dragging card %s %d,%d", c.id, c.screenX, c.screenY))
			c.Draw(screen)
		}
	}
}
