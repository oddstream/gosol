package sol

import (
	"reflect"

	"oddstream.games/gosol/util"
)

// Foundation is the destination for cards
type Foundation struct {
	Pile
}

// New fills in basic information
func (f *Foundation) New(info map[string]string) {
	f.x = util.GetIntFromMap(info, "x")
	f.y = util.GetIntFromMap(info, "y")
	f.fan = util.GetStringFromMap(info, "fan")
	f.accept = util.GetIntFromMap(info, "accept")

	f.createImage()
}

// Class returns the type of this Pile
func (f *Foundation) Class() string {
	return reflect.TypeOf(*f).Name() // .String() returns "sol.Stock"
}

// Pop here does not allow cards to be moved from Foundation
func (f *Foundation) Pop() *Card {
	return nil
}

// CanAcceptCard returns true if this Pile can accept the Card
func (f *Foundation) CanAcceptCard(c *Card) bool {
	if len(f.cards) == 0 {
		if c.ordinal == f.accept {
			return true
		}
	} else {
		// TODO build rules
		cf := f.Peek()
		if cf.suit == c.suit && cf.ordinal+1 == c.ordinal {
			return true
		}
	}
	return false
}
