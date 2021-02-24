package sol

import "log"

// SaveableBaize is a reduced struct for converting to JSON
type SaveableBaize struct {
	Variant string
	Seed    int64
	Piles   []SaveablePile
}

// SaveablePile is a reduced struct for converting to JSON
type SaveablePile struct {
	Class    string // for readability and sanity checks
	Accept   int    // local, mutable copy of Accept
	Recycles int    // local, mutable copy of Recycles
	Cards    []SaveableCard
}

// SaveableCard is a reduced struct for converting to JSON
type SaveableCard struct {
	ID    string
	Prone bool
}

// Saveable creates a saveable version of the current state
func (b *Baize) Saveable() SaveableBaize {
	sav := SaveableBaize{Variant: b.Variant, Seed: b.Seed}
	for _, p := range b.Piles {
		sav.Piles = append(sav.Piles, p.Saveable())
	}
	return sav
}

// UpdateFromSaveable updates the contents of the Piles from a saved copy of a previous state
func (b *Baize) UpdateFromSaveable(sav SaveableBaize) {

	var cardCache []*Card = nil

	for _, p := range b.Piles {
		for _, c := range p.Cards {
			cardCache = append(cardCache, c)
		}
	}

	for i := 0; i < len(b.Piles); i++ {
		pile := b.Piles[i]
		savedPile := sav.Piles[i]
		if len(pile.Cards) != len(savedPile.Cards) {
			pile.UpdateFromSaved(cardCache, savedPile)
		}
	}
}

// Saveable returns a reduced object for converting to JSON and saving
func (p *Pile) Saveable() SaveablePile {
	sav := SaveablePile{Class: p.Class, Accept: p.localAccept, Recycles: p.localRecycles}
	for _, c := range p.Cards {
		sav.Cards = append(sav.Cards, c.Saveable())
	}
	return sav
}

// UpdateFromSaved replaces this Pile's contents
func (p *Pile) UpdateFromSaved(cardCache []*Card, sav SaveablePile) {

	findCardInCache := func(id string) *Card {
		for _, c := range cardCache {
			if c.id == id {
				return c
			}
		}
		return nil
	}

	if sav.Class != p.Class {
		log.Fatal(p.Class, "!=", sav.Class)
	}

	p.Cards = nil
	p.localAccept = sav.Accept
	p.localRecycles = sav.Recycles
	for _, cSaved := range sav.Cards {
		c := findCardInCache(cSaved.ID)
		if cSaved.Prone != c.prone { // TODO copy this back to Opsole
			if cSaved.Prone {
				c.FlipDown()
			} else {
				c.FlipUp()
			}
		}
		p.Push(c)
	}
}

// Saveable returns a reduced object for converting to JSON and saving
func (c *Card) Saveable() SaveableCard {
	return SaveableCard{ID: c.id, Prone: c.prone}
}

// NewCardFromSaveable is a factory for Card objects
// func NewCardFromSaveable(sav SaveableCard) *Card {
// 	p, s, o := parseID(sav.ID)
// 	c := NewCard(p, s, o)
// 	c.prone = sav.Prone
// 	return c
// }
