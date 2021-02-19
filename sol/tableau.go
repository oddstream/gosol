package sol

import "oddstream.games/gosol/util"

// Tableau is the destination for cards
type Tableau struct {
	Pile

	class string
}

// TableauInfo contains configuration for all Tableau objects
type TableauInfo struct {
	Accept int // ordinal of card to accept on empty pile, 0 == any
}

// New fills in basic information
func (t *Tableau) New(info map[string]string) {
	t.class = "Tableau"
	t.x = util.GetIntFromMap(info, "x")
	t.y = util.GetIntFromMap(info, "y")
	t.fan = util.GetStringFromMap(info, "fan")
}

// Class returns the class of this Pile
func (t *Tableau) Class() string {
	return t.class
}
