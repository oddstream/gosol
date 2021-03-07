package sol

import (
	"hash/crc32"
	"log"
)

// SaveableBaize is a reduced struct for converting to JSON
type SaveableBaize struct {
	Checksum uint32
	Variant  string
	Seed     int64
	State    BaizeState
	Piles    []SaveablePile
}

// SaveablePile is a reduced struct for converting to JSON
type SaveablePile struct {
	Class    string // for readability and sanity checks
	Accept   int    // local, mutable copy of Accept
	Recycles int    // local, mutable copy of Recycles
	Scrunch  int    // copy of scrunch percentage
	Cards    []string
}

// Checksum creates checksum for the current state
func (b *Baize) Checksum() uint32 {
	// https://golang.org/src/hash/crc32/example_test.go
	var lens []byte
	// crc32q := crc32.MakeTable(0xD5828281)
	for _, p := range b.Piles {
		lens = append(lens, byte(p.CardCount()))
	}
	// return crc32.Checksum(lens, crc32q)
	return crc32.ChecksumIEEE(lens)
}

// Saveable creates a saveable version of the current state
func (b *Baize) Saveable() SaveableBaize {
	sav := SaveableBaize{Checksum: b.Checksum(), Variant: b.Variant, Seed: b.Seed, State: b.State}
	for _, p := range b.Piles {
		sav.Piles = append(sav.Piles, p.Saveable())
	}
	// println("Checksum", sav.Checksum)
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
		if pile.Class != savedPile.Class {
			log.Fatal("saved pile", savedPile.Class, "does not match baize pile", pile.Class)
		}
		if len(pile.Cards) != len(savedPile.Cards) {
			// println("updating pile", pile.Class)
			pile.UpdateFromSaved(cardCache, savedPile)
		}
	}

	b.State = sav.State
}

// Saveable returns a reduced object for converting to JSON and saving
func (p *Pile) Saveable() SaveablePile {
	sav := SaveablePile{Class: p.Class, Accept: p.localAccept, Recycles: p.localRecycles, Scrunch: p.scrunchPercentage}
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

	p.Cards = nil
	p.localAccept = sav.Accept
	p.localRecycles = sav.Recycles
	p.scrunchPercentage = sav.Scrunch
	for _, cSaved := range sav.Cards {
		id := cSaved[0:3] // substring operation up to, but not including, cSaved[3]
		var prone bool
		switch string(cSaved[3]) {
		case "d":
			prone = true
		case "u":
			prone = false
		default:
			log.Fatal("unexpected saved card", cSaved)
		}
		c := findCardInCache(id)
		if c == nil {
			log.Fatal("could not find card in cache", id)
		}
		p.Push(c)
		if prone != c.prone { // TODO copy this back to Opsole
			if prone {
				c.FlipDown()
			} else {
				c.FlipUp()
			}
		}
	}
}

// Saveable returns a 4-char string for converting to JSON and saving
func (c *Card) Saveable() string {
	if c.prone {
		return c.id + "d"
	}
	return c.id + "u"
}
