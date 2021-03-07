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

// CardID holds flags (marked, prone), pack, suit, ordinal
type CardID uint32

func (cid CardID) String() string {
	return fmt.Sprintf("%d %s %d", cid.Pack(), cid.StringSuit(), cid.Ordinal())
}

/*
   Precedence    Operator
   5             *  /  %  <<  >>  &  &^
   4             +  -  |  ^
   3             ==  !=  <  <=  >  >=
   2             &&
   1             ||
*/

// Pack returns the pack number buried in the card id
func (cid CardID) Pack() int {
	return int((uint32(cid) & packMask) >> 8)
}

// Pack returns the pack number this card belongs to
func (c *Card) Pack() int {
	return c.ID.Pack()
}

// Suit returns the suit number buried in the card id
func (cid CardID) Suit() int {
	return int((uint32(cid) & suitMask) >> 4)
}

// StringSuit returns the suit number buried in the card id, expressed as a string
func (cid CardID) StringSuit() string {
	switch cid.Suit() {
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

// Suit returns 1=club, 2=diamond, 3=heart, 4=spade
func (c *Card) Suit() int {
	return c.ID.Suit()
}

// StringSuit returns the suit as a string
func (c *Card) StringSuit() string {
	return c.ID.StringSuit()
}

// Ordinal returns the ordinal number buried in the card id
func (cid CardID) Ordinal() int {
	return int(uint32(cid) & ordinalMask)
}

// Ordinal returns the face value of this card 1..13
func (c *Card) Ordinal() int {
	return c.ID.Ordinal()
}

// Prone returns the prone flag buried in the card id
func (cid CardID) Prone() bool {
	return uint32(cid)&proneFlag == proneFlag
}

// Prone returns true if the card is face down, false if it is face up
func (c *Card) Prone() bool {
	return c.ID.Prone()
}

// Marked returns the marked flag buried in the card id
func (cid CardID) Marked() bool {
	return uint32(cid)&markedFlag == markedFlag
}

// Marked returns true if the card is marked, false if it is unmarked
func (c *Card) Marked() bool {
	return c.ID.Marked()
}

// SetProne true or false
func (c *Card) SetProne(prone bool) {
	if prone {
		c.ID = CardID(uint32(c.ID) | proneFlag)
	} else {
		c.ID = CardID(uint32(c.ID) & ^proneFlag)
	}
}

// SetMarked true or false
func (c *Card) SetMarked(marked bool) {
	if marked {
		c.ID = CardID(uint32(c.ID) | markedFlag)
	} else {
		c.ID = CardID(uint32(c.ID) & ^markedFlag)
	}
}

// Color returns Red or Black
func (cid CardID) Color() color.RGBA {
	suit := cid.Suit()
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
	return c.ID.Color()
}

// NewCardID constructor
func NewCardID(pack, suit, ordinal int) CardID {
	var u uint32
	u += uint32(pack) << 8
	u += uint32(suit) << 4
	u += uint32(ordinal)
	return CardID(u)
}

func sameCard(ID1, ID2 CardID) bool {
	return uint32(ID1)&cardMask == uint32(ID2)&cardMask
}
