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
	Cards    []SaveableCard
}

// SaveableCard is a reduced struct for converting to JSON
type SaveableCard struct {
	ID    string
	Prone bool
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
			println("updating pile", pile.Class)
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
		c := findCardInCache(cSaved.ID)
		p.Push(c)
		if cSaved.Prone != c.prone { // TODO copy this back to Opsole
			if cSaved.Prone {
				c.FlipDown()
			} else {
				c.FlipUp()
			}
		}
	}
}

// Saveable returns a reduced object for converting to JSON and saving
func (c *Card) Saveable() SaveableCard {
	return SaveableCard{ID: c.id, Prone: c.prone}
}

// Bits returns the compact form of a card id + prone flag
// func (c *Card) Bits() uint16 {
// 	// 1111000011110000
// 	var ui uint16
// 	ui = uint16(c.pack) << 12
// 	switch c.suit {
// 	case "Club":
// 		ui |= 0b00010000
// 	case "Diamond":
// 		ui |= 0b00100000
// 	case "Heart":
// 		ui |= 0b00110000
// 	case "Spade":
// 		ui |= 0b01000000
// 	}
// 	ui |= uint16(c.ordinal) << 4 // 1=0b0001, 13=0b1101
// 	if c.prone {
// 		ui |= 1
// 	}
// 	return ui
// }

// ParseBits unpacks the compact form of a card id+prone flag
// func (c *Card) ParseBits(ui uint16) (id string, prone bool) {
// 	var pack, ordinal int
// 	var suit string
// 	pack = int(ui >> 12 & 0b1111)
// 	switch ui & 0b11110000 {
// 	case 0b00010000:
// 		suit = "Club"
// 	case 0b00100000:
// 		suit = "Diamond"
// 	case 0b00110000:
// 		suit = "Heart"
// 	case 0b01000000:
// 		suit = "Spade"
// 	}
// 	ordinal = int(ui >> 4 & 0b1111)
// 	prone = ui&1 == 1
// 	id = fmt.Sprintf("%d%c%02d", pack, suit[0], ordinal)
// 	return
// }

// NewCardFromSaveable is a factory for Card objects
// func NewCardFromSaveable(sav SaveableCard) *Card {
// 	p, s, o := parseID(sav.ID)
// 	c := NewCard(p, s, o)
// 	c.prone = sav.Prone
// 	return c
// }
