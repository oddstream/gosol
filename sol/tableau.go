package sol

import (
	"reflect"

	"oddstream.games/gosol/util"
)

// Tableau is the destination for cards
type Tableau struct {
	Pile
}

// New fills in basic information
func (t *Tableau) New(info map[string]string) {

	t.x = util.GetIntFromMap(info, "x")
	t.y = util.GetIntFromMap(info, "y")
	t.fan = util.GetStringFromMap(info, "fan")
	t.deal = util.GetStringFromMap(info, "deal")
	t.accept = util.GetIntFromMap(info, "accept")

	t.createImage()
}

// Class returns the type of this Pile
func (t *Tableau) Class() string {
	return reflect.TypeOf(*t).Name() // .String() returns "sol.Stock"
}
