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
const proneFlag uint32 = 0b1000000000000
const movableFlag uint32 = 0b10000000000000

// const flagMask uint32 = 0b1111000000000000

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
	return SuitIntToString(cid.Suit())
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

// Movable returns the movable flag buried in the card id
func (cid CardID) Movable() bool {
	return uint32(cid)&movableFlag == movableFlag
}

// Movable returns true if the card is marked as movable, false if it is unmarked
func (c *Card) Movable() bool {
	return c.ID.Movable()
}

// SetProne true or false
func (c *Card) SetProne(prone bool) {
	if prone {
		c.ID = CardID(uint32(c.ID) | proneFlag)
	} else {
		c.ID = CardID(uint32(c.ID) & ^proneFlag)
	}
}

// SetMovable set the "this card is movable" flag
func (c *Card) SetMovable(movable bool) {
	if movable {
		c.ID = CardID(uint32(c.ID) | movableFlag)
	} else {
		c.ID = CardID(uint32(c.ID) & ^movableFlag)
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

// SameCard returns true if the two cards have the same ordinal and suit
func SameCard(ID1, ID2 CardID) bool {
	return uint32(ID1)&(suitMask|ordinalMask) == uint32(ID2)&(suitMask|ordinalMask)
}

// SameCardAndPack returns true if the two card IDs have the same ordinal and suit, and are from the same pack
func SameCardAndPack(ID1, ID2 CardID) bool {
	return uint32(ID1)&(packMask|suitMask|ordinalMask) == uint32(ID2)&(packMask|suitMask|ordinalMask)
}

// SuitStringToInt converts a suit string ("Heart") to an int (HEART)
func SuitStringToInt(suit string) int {
	switch suit {
	case "Club":
		return CLUB
	case "Diamond":
		return DIAMOND
	case "Heart":
		return HEART
	case "Spade":
		return SPADE
	}
	return 0
}

// SuitIntToString converts a suit int (HEART) to a string ("Heart")
func SuitIntToString(suit int) string {
	switch suit {
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
