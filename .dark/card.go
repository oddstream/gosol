package dark

import (
	"oddstream.games/gosol/cardid"
)

// Card holds the state of the cards.
// Card is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Card object as immutable, hence the unexported fields and getters.
type Card struct {
	id             cardid.CardID
	owningPile     *Pile
	tapDestination *Pile
	tapWeight      int
}

func NewCard(pack, suit, ordinal int) Card {
	c := Card{id: cardid.NewCardID(pack, suit, ordinal)}
	c.setProne(true)
	return c
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

func (c *Card) TapWeight() int {
	return c.tapWeight
}

// Private functions

func (c *Card) owner() *Pile {
	return c.owningPile
}

func (c *Card) setOwner(p *Pile) *Pile {
	c.owningPile = p
}

func (c *Card) setProne(prone bool) {
	c.id = c.id.SetProne(prone)
}

func (c *Card) flipUp() {
	if c.Prone() {
		c.setProne(false)
	}
}

func (c *Card) flipDown() {
	if !c.Prone() {
		c.setProne(true)
	}
}
