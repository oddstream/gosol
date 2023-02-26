package cardid

import (
	"fmt"
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
	packMask    CardID = 0b0000111100000000 // 0..15 0xf00
	suitMask    CardID = 0b0000000011110000 // 0..15 0xf0
	ordinalMask CardID = 0b0000000000001111 // 0..15 0xf
	proneFlag   CardID = 0b0001000000000000 // single bit
	jokerFlag   CardID = 0b0010000000000000 // single bit
)

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
	return int((cid & packMask) >> 8)
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
	case CLUB:
		return CLUB_RUNE
	case DIAMOND:
		return DIAMOND_RUNE
	case HEART:
		return HEART_RUNE
	case SPADE:
		return SPADE_RUNE
	default:
		return 0
	}
}

// Ordinal returns the ordinal number buried in the card id
func (cid CardID) Ordinal() int {
	return int(cid & ordinalMask)
}

// Prone returns the prone flag buried in the card id
func (cid CardID) Prone() bool {
	return cid&proneFlag == proneFlag
}

// SetProne true or false
func (cid CardID) SetProne(prone bool) CardID {
	if prone {
		cid = cid | proneFlag
	} else {
		cid = cid & (^proneFlag)
	}
	return cid
}

// Joker returns the joker flag buried in the card id
func (cid CardID) Joker() bool {
	return cid&jokerFlag == jokerFlag
}

func (cid CardID) Black() bool {
	suit := cid.Suit()
	return suit == CLUB || suit == SPADE
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

func (cid CardID) PackSuitOrdinal() CardID {
	return cid & (packMask | suitMask | ordinalMask)
}

// SameCard returns true if the two cards have the same ordinal and suit; pack is ignored
func SameCard(ID1, ID2 CardID) bool {
	return ID1&(suitMask|ordinalMask) == ID2&(suitMask|ordinalMask)
}

// SameCardAndPack returns true if the two card IDs have the same ordinal and suit, and are from the same pack
func SameCardAndPack(ID1, ID2 CardID) bool {
	return ID1.PackSuitOrdinal() == ID2.PackSuitOrdinal()
}

// SuitStringToInt converts a suit string ("Heart") to an int (HEART)
// func SuitStringToInt(suit string) int {
// 	switch suit {
// 	case "":
// 		return NOSUIT
// 	case "Club":
// 		return CLUB
// 	case "Diamond":
// 		return DIAMOND
// 	case "Heart":
// 		return HEART
// 	case "Spade":
// 		return SPADE
// 	}
// 	return 0
// }

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
	default:
		return ""
	}
}
