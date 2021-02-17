package sol

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Pile is a generic container for cards
type Pile struct {
	cards []*Card
	X, Y  int // grid position of Pile
}

// CreateCards fills the pile with packs*52 new cards
func (p *Pile) CreateCards(packs int) *Pile {
	// gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range [4]string{"Club", "Diamond", "Heart", "Spade"} {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				c.owner = p
				c.PositionTo(p.X*71, p.Y*96)
				p.cards = append(p.cards, c)
			}
		}
	}
	// println("created", len(p.cards), "cards")
	return p
}

// Position returns the x,y screen coords of this pile
func (p *Pile) Position() (int, int) {
	return p.X * 71, p.Y * 96
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

// Update the baize state (transitions, user input)
func (p *Pile) Update() error {
	for _, c := range p.cards {
		c.Update()
	}
	return nil
}

// Draw renders the card into the screen
func (p *Pile) Draw(screen *ebiten.Image) {
	for _, c := range p.cards {
		c.Draw(screen)
	}
}
