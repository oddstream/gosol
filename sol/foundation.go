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
