package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"

	"oddstream.games/gosol/util"
)

type CardPair struct {
	c1, c2 *Card
}

type CardPairs []CardPair

type CardPairCompareFunc func(CardPair) (bool, error)

func TailConformant(tail []*Card, fn CardPairCompareFunc) (bool, error) {
	for _, pair := range NewCardPairs(tail) {
		if ok, err := fn(pair); !ok {
			return false, err
		}
	}
	return true, nil
}

// UnsortedPairs
//
// A generic way of calculating the number of unsorted card pairs in a pile.
// Called by *Pile.vtable.UnsortedPairs()
func UnsortedPairs(pile *Pile, fn CardPairCompareFunc) int {
	if pile.Len() < 2 {
		return 0
	}
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.c1.Prone() || pair.c2.Prone() {
			unsorted++
		} else {
			if ok, _ := fn(pair); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func NewCardPairs(cards []*Card) CardPairs {
	if len(cards) < 2 {
		return []CardPair{} // always return a list, not nil
	}
	var cpairs []CardPair
	c1 := cards[0]
	for i := 1; i < len(cards); i++ {
		c2 := cards[i]
		cpairs = append(cpairs, CardPair{c1, c2})
		c1 = c2
	}
	return cpairs
}

// func (cpairs CardPairs) Print() {
// 	for _, pair := range cpairs {
// 		println(pair.c1.String(), pair.c2.String())
// 	}
// }

func Compare_Empty(p *Pile, c *Card) (bool, error) {
	if p.Label() != "" {
		if p.Label() == "x" || p.Label() == "X" {
			return false, errors.New("Cannot move cards to that empty pile")
		}
		ord := util.OrdinalToShortString(c.Ordinal())
		if ord != p.Label() {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(p.Label()), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}

// little library of simple compares

func (cp CardPair) Compare_Up() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must be in ascending sequence")
}

func (cp CardPair) Compare_UpWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (cp CardPair) Compare_Down() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	return false, errors.New("Cards must be in descending sequence")
}

func (cp CardPair) Compare_DownWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
}

func (cp CardPair) Compare_UpOrDown() (bool, error) {
	if !(cp.c1.Ordinal()+1 == cp.c2.Ordinal() || cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_UpOrDownWrap() (bool, error) {
	if (cp.c1.Ordinal()+1 == cp.c2.Ordinal()) || (cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return true, nil
	} else if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace On King
	} else if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	} else {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
}

func (cp CardPair) Compare_Color() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return true, nil
}

func (cp CardPair) Compare_AltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return true, nil
}

func (cp CardPair) Compare_Suit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return true, nil
}

func (cp CardPair) Compare_OtherSuit() (bool, error) {
	if cp.c1.Suit() == cp.c2.Suit() {
		return false, errors.New("Cards must not be the same suit")
	}
	return true, nil
}

// library of compare functions made from simple compares

func (cp CardPair) Compare_DownColor() (bool, error) {
	ok, err := cp.Compare_Color()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownAltColor() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

// Compare_DownColorWrap not used
func (cp CardPair) Compare_DownColorWrap() (bool, error) {
	ok, err := cp.Compare_Color()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

func (cp CardPair) Compare_DownAltColorWrap() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

// Compare_UpAltColor not used
func (cp CardPair) Compare_UpAltColor() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_UpSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_DownSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpOrDownSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpOrDown()
}

func (cp CardPair) Compare_UpOrDownSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpOrDownWrap()
}

// Compare_DownOtherSuit not used
func (cp CardPair) Compare_DownOtherSuit() (bool, error) {
	ok, err := cp.Compare_OtherSuit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpWrap()
}

func (cp CardPair) Compare_DownSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

// ChainCall
//
// Call using CardPair method expressions
// eg ChainCall(CardPair.Compare_UpOrDown, CardPair.Compare_Suit)
//
// TODO think of something else for UnsortedPairs(*Pile)
func (cp CardPair) ChainCall(fns ...func(CardPair) (bool, error)) (ok bool, err error) {
	for _, fn := range fns {
		if ok, err = fn(cp); err != nil {
			break
		}
	}
	return
}
