package sol

import (
	"reflect"

	"oddstream.games/gosol/util"
)

// Waste is the destination for cards
type Waste struct {
	Pile
}

// New fills in basic information
func (w *Waste) New(info map[string]string) {
	w.x = util.GetIntFromMap(info, "x")
	w.y = util.GetIntFromMap(info, "y")
	w.fan = util.GetStringFromMap(info, "fan")

	w.createImage()
}

// Class returns the type of this Pile
func (w *Waste) Class() string {
	return reflect.TypeOf(*w).Name() // .String() returns "sol.Stock"
}

// CanAcceptCard returns true if this Pile can accept the Card
func (w *Waste) CanAcceptCard(*Card) bool {
	return true
}
