package sol

import "github.com/hajimehoshi/ebiten/v2"

// Baize object describes the baize
type Baize struct {
}

// NewBaize is the factory func for Baize object
func NewBaize() *Baize {
	b := &Baize{}
	return b
}

// Layout implements ebiten.Game's Layout.
func (g *Baize) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (g *Baize) Update() error {
	return nil
}

// Draw renders the baize into the screen
func (g *Baize) Draw(screen *ebiten.Image) {

	screen.Fill(colorBaize)

	{
		c := NewCard(0, "Spade", 1)
		c.screenX, c.screenY = 100, 100
		c.Draw(screen)
		c.prone = false
		c.screenX, c.screenY = 200, 100
		c.Draw(screen)
	}
}
