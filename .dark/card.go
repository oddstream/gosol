package dark

import (
	"oddstream.games/gosol/cardid"
)

// Card holds the state of the cards.
// Card is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Card object as immutable, hence the unexported fields and getters.
type Card struct {
	id             cardid.CardID
	owner          *Pile
	tapDestination *Pile
	tapWeight      int
}

// Public functions

func (c *Card) Suit() int {
	return c.id.Suit()
}

func (c *Card) Ordinal() int {
	return c.id.Ordinal()
}

func (c *Card) ID() cardid.CardID {
	return c.id
}

func (c *Card) Prone() bool {
	return c.id.Prone()
}

func (c *Card) SetProne(prone bool) {
	c.id = c.id.SetProne(prone)
}

func (c *Card) TapWeight() int {
	return c.tapWeight
}

// Private functions

func (c *Card) owner() *Pile {
	return c.owner
}

func (c *Card) flipUp() {
	if c.Prone() {
		c.SetProne(false)
	}
}

func (c *Card) flipDown() {
	if !c.Prone() {
		c.SetProne(true)
	}
}
