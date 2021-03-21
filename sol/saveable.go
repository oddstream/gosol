package sol

import (
	"hash/crc32"
	"log"
)

// SaveableBaize is a reduced struct for converting to JSON
type SaveableBaize struct {
	Checksum uint32
	Seed     int64
	State    BaizeState
	Piles    []SaveablePile
}

// SaveablePile is a reduced struct for converting to JSON
type SaveablePile struct {
	Class    string   // for readability and sanity checks
	Accept   int      // local, mutable copy of Accept
	Recycles int      // local, mutable copy of Recycles
	Scrunch  int      // copy of scrunch percentage
	Cards    []CardID // array of Card.ID
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
	sav := SaveableBaize{Checksum: b.Checksum(), Seed: b.Seed, State: b.State}
	for _, p := range b.Piles {
		sav.Piles = append(sav.Piles, p.Saveable())
	}
	// println("Checksum", sav.Checksum)
	return sav
}

// UpdateFromSaveable updates the contents of the Piles from a saved copy of a previous state
func (b *Baize) UpdateFromSaveable(sav SaveableBaize) {

	if len(sav.Piles) == 0 {
		log.Panic("bad SaveableBaize passed to UpdateFromSaveable()")
	}

	var cardCache []*Card

	for _, p := range b.Piles {
		// S1011 – Use a single append to concatenate two slices
		// for _, c := range p.Cards {
		// 	cardCache = append(cardCache, c)
		// }
		cardCache = append(cardCache, p.Cards...) // append a slice to a slice
	}

	for i := 0; i < len(b.Piles); i++ {
		pile := b.Piles[i]
		savedPile := sav.Piles[i]
		if pile.Class != savedPile.Class {
			log.Panic("saved pile", savedPile.Class, "does not match baize pile", pile.Class)
		}
		if len(pile.Cards) != len(savedPile.Cards) {
			pile.UpdateFromSaved(cardCache, savedPile)
		}
	}

	b.Seed = sav.Seed // TODO ???
	b.State = sav.State
}

// Saveable returns a reduced Pile object for converting to JSON and saving
func (p *Pile) Saveable() SaveablePile {
	sav := SaveablePile{Class: p.Class, Accept: p.localAccept, Recycles: p.localRecycles, Scrunch: p.scrunchPercentage}
	for _, c := range p.Cards {
		sav.Cards = append(sav.Cards, c.ID)
	}
	return sav
}

// UpdateFromSaved replaces this Pile's contents
func (p *Pile) UpdateFromSaved(cardCache []*Card, sav SaveablePile) {

	findCardInCache := func(ID CardID) *Card {
		for _, c := range cardCache {
			if sameCard(c.ID, ID) {
				return c
			}
		}
		return nil
	}

	p.Cards = nil
	p.localAccept = sav.Accept
	p.localRecycles = sav.Recycles
	p.scrunchPercentage = sav.Scrunch
	for _, savedID := range sav.Cards {
		c := findCardInCache(savedID)
		if c == nil {
			log.Panic("could not find card in cache", savedID, savedID.String())
		}
		p.Push(c)
		if savedID.Prone() {
			c.FlipDown()
		} else {
			c.FlipUp()
		}
	}
}
