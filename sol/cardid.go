package sol

import (
	"fmt"
	"image/color"
	"log"
)

// enum types for card suits
const (
	CLUB    = 1
	DIAMOND = 2
	HEART   = 3
	SPADE   = 4
)

const packMask uint32 = 0b111100000000
const suitMask uint32 = 0b000011110000
const ordinalMask uint32 = 0b1111
const cardMask uint32 = (packMask | suitMask | ordinalMask)
const proneFlag uint32 = 0b1000000000000
const markedFlag uint32 = 0b10000000000000
const flagMask uint32 = 0b1111000000000000

// type CardID uint32

/*
   Precedence    Operator
   5             *  /  %  <<  >>  &  &^
   4             +  -  |  ^
   3             ==  !=  <  <=  >  >=
   2             &&
   1             ||
*/

// func (cid CardID) pack() int {
// 	return int((uint32(cid) & packMask) >> 8)
// }

func packFromCardID(ID uint32) int {
	return int(ID & packMask >> 8)
}

// Pack returns the pack number this card belongs to
func (c *Card) Pack() int {
	return packFromCardID(c.ID)
}

func suitFromCardID(ID uint32) int {
	return int(ID & suitMask >> 4)
}

// Suit returns 1=club, 2=diamond, 3=heart, 4=spade
func (c *Card) Suit() int {
	return suitFromCardID(c.ID)
}

func ordinalFromCardID(ID uint32) int {
	return int(ID & ordinalMask)
}

// Ordinal returns the face value of this card 1..13
func (c *Card) Ordinal() int {
	return ordinalFromCardID(c.ID)
}

func proneFromCardID(ID uint32) bool {
	return ID&proneFlag == proneFlag
}

// Prone returns true if the card is face down, false if it is face up
func (c *Card) Prone() bool {
	return proneFromCardID(c.ID)
}

// Marked returns true if the card is marked, false if it is unmarked
func (c *Card) Marked() bool {
	return c.ID&markedFlag == markedFlag
}

// SetProne true or false
func (c *Card) SetProne(prone bool) {
	if prone {
		c.ID |= proneFlag
	} else {
		c.ID &= ^proneFlag
	}
}

// SetMarked true or false
func (c *Card) SetMarked(prone bool) {
	if prone {
		c.ID |= markedFlag
	} else {
		c.ID &= ^markedFlag
	}
}

// ColorFromCardID returns Red or Black
func colorFromCardID(ID uint32) color.RGBA {
	suit := int(ID & suitMask >> 4)
	switch suit {
	case HEART, DIAMOND:
		return BasicColors["Red"]
	case CLUB, SPADE:
		return BasicColors["Black"]
	default:
		log.Fatal("unknown suit in id", suit)
	}
	return BasicColors["Purple"]
}

// Color returns Red or Black
func (c *Card) Color() color.RGBA {
	return colorFromCardID(c.ID)
}

func makeCardID(pack, suit, ordinal int) uint32 {
	var u uint32
	u += uint32(pack) << 8
	u += uint32(suit) << 4
	u += uint32(ordinal)
	return u
}

// func scalableID(suit, ordinal int) uint32 {
// 	return makeCardID(0, suit, ordinal)
// }

func stringSuitFromCardID(ID uint32) string {
	switch suitFromCardID(ID) {
	case CLUB:
		return "Club"
	case DIAMOND:
		return "Diamond"
	case HEART:
		return "Heart"
	case SPADE:
		return "Spade"
	}
	return ""
}

// StringSuit returns the suit as a string
func (c *Card) StringSuit() string {
	return stringSuitFromCardID(c.ID)
}

func sameCard(ID1, ID2 uint32) bool {
	return ID1&cardMask == ID2&cardMask
}

func cardIDToString(ID uint32) string {
	return fmt.Sprintf("%d %s %d", packFromCardID(ID), stringSuitFromCardID(ID), ordinalFromCardID(ID))
}
