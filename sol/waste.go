package sol

import "oddstream.games/gosol/util"

// Waste is the destination for cards
type Waste struct {
	Pile

	class string
}

// New fills in basic information
func (w *Waste) New(info map[string]string) {
	w.class = "Waste"
	w.x = util.GetIntFromMap(info, "x")
	w.y = util.GetIntFromMap(info, "y")
	w.fan = util.GetStringFromMap(info, "fan")
}

// Class returns the class of this Pile
func (w *Waste) Class() string {
	return w.class
}
