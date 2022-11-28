package sol

import (
	"fmt"
	"image/color"
)

// enum types for card suits
const (
	NOSUIT  = 0
	CLUB    = 1
	DIAMOND = 2
	HEART   = 3
	SPADE   = 4
)

const (
	CLUB_RUNE    = rune(9827) // 0x2663
	DIAMOND_RUNE = rune(9830) // 0x2666
	HEART_RUNE   = rune(9829) // 0x2665
	SPADE_RUNE   = rune(9824) // 0x2660
)

// CardID holds flags (prone &c), pack, suit, ordinal
// CardID is crammed like this to make JSON smaller
type CardID uint16

const (
	packMask    CardID = 0b0000111100000000
	suitMask    CardID = 0b0000000011110000
	ordinalMask CardID = 0b0000000000001111
	proneFlag   CardID = 0b0001000000000000
	jokerFlag   CardID = 0b0010000000000000
)

func (cid CardID) String() string {
	return fmt.Sprintf("%d %d %s", cid.Pack(), cid.Ordinal(), cid.StringSuit())
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
	return int((cid & packMask) >> 8)
}

// Pack returns the pack number this card belongs to
func (c *Card) Pack() int {
	return c.ID.Pack()
}

// Suit returns the suit number buried in the card id
func (cid CardID) Suit() int {
	return int((cid & suitMask) >> 4)
}

// StringSuit returns the suit number buried in the card id, expressed as a string
func (cid CardID) StringSuit() string {
	return SuitIntToString(cid.Suit())
}

// SuitRune returns the unicode rune/glyph/symbol for this suit
func (cid CardID) SuitRune() (r rune) {
	switch cid.Suit() {
	case NOSUIT:
		r = 0
	case CLUB:
		r = CLUB_RUNE
	case DIAMOND:
		r = DIAMOND_RUNE
	case HEART:
		r = HEART_RUNE
	case SPADE:
		r = SPADE_RUNE
	}
	return
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
	return int(cid & ordinalMask)
}

// Ordinal returns the face value of this card 1..13
func (c *Card) Ordinal() int {
	return c.ID.Ordinal()
}

// Prone returns the prone flag buried in the card id
func (cid CardID) Prone() bool {
	return cid&proneFlag == proneFlag
}

// Prone returns true if the card is face down, false if it is face up
func (c *Card) Prone() bool {
	return c.ID.Prone()
}

// SetProne true or false
func (c *Card) SetProne(prone bool) {
	if prone {
		c.ID = c.ID | proneFlag
	} else {
		c.ID = c.ID & (^proneFlag)
	}
}

// Prone returns the joker flag buried in the card id
func (cid CardID) Joker() bool {
	return cid&jokerFlag == jokerFlag
}

// Prone returns true if the card is face down, false if it is face up
func (c *Card) Joker() bool {
	return c.ID.Joker()
}

// Color returns Red or Black
func (cid CardID) Color() color.RGBA {
	suit := cid.Suit()
	if ThePreferences.ColorfulCards {
		switch TheBaize.script.CardColors() {
		case 4:
			switch suit {
			case NOSUIT:
				return BasicColors["Silver"]
			case CLUB:
				return ExtendedColors[ThePreferences.ClubColor]
			case DIAMOND:
				return ExtendedColors[ThePreferences.DiamondColor]
			case HEART:
				return ExtendedColors[ThePreferences.HeartColor]
			case SPADE:
				return ExtendedColors[ThePreferences.SpadeColor]
			}
		case 2:
			switch suit {
			case NOSUIT:
				return BasicColors["Silver"]
			case CLUB, SPADE:
				return ExtendedColors[ThePreferences.BlackColor]
			case DIAMOND, HEART:
				return ExtendedColors[ThePreferences.RedColor]
			}
		case 1:
			return ExtendedColors[ThePreferences.SpadeColor]
		}
	} else {
		switch suit {
		case NOSUIT:
			return BasicColors["Silver"]
		case CLUB, SPADE:
			return ExtendedColors[ThePreferences.BlackColor]
		case DIAMOND, HEART:
			return ExtendedColors[ThePreferences.RedColor]
		}
	}
	return BasicColors["Purple"]
}

// Color returns Red or Black
func (c *Card) Color() color.RGBA {
	return c.ID.Color()
}

func (cid CardID) Black() bool {
	suit := cid.Suit()
	return suit == CLUB || suit == SPADE
}

func (c *Card) Black() bool {
	return c.ID.Black()
}

// NewCardID constructor
func NewCardID(pack, suit, ordinal int) CardID {
	var u uint32
	u += uint32(pack) << 8
	u += uint32(suit) << 4
	u += uint32(ordinal)
	if suit == NOSUIT && ordinal == 0 {
		u += uint32(jokerFlag)
	}
	return CardID(u)
}

// SameCard returns true if the two cards have the same ordinal and suit; pack is ignored
func SameCard(ID1, ID2 CardID) bool {
	return ID1&(suitMask|ordinalMask) == ID2&(suitMask|ordinalMask)
}

// SameCardAndPack returns true if the two card IDs have the same ordinal and suit, and are from the same pack
func SameCardAndPack(ID1, ID2 CardID) bool {
	return ID1&(packMask|suitMask|ordinalMask) == ID2&(packMask|suitMask|ordinalMask)
}

// SuitStringToInt converts a suit string ("Heart") to an int (HEART)
func SuitStringToInt(suit string) int {
	switch suit {
	case "":
		return NOSUIT
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
	case NOSUIT:
		return ""
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
