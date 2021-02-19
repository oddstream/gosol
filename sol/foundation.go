package sol

import "oddstream.games/gosol/util"

// Foundation is the destination for cards
type Foundation struct {
	Pile

	class string
}

// FoundationInfo contains configuration for all Foundation objects
type FoundationInfo struct {
	Accept      int    // ordinal of card to accept on empty pile, 0 == any (FoundationSpider has it's own rules)
	DealOrdinal int    // ordinal of card to deal here, 0 == don't (why isn't this in Deal= parameter?)
	DealSuit    string // "" == any
}

// New fills in basic information
func (f *Foundation) New(info map[string]string) {
	f.class = "Foundation"
	f.x = util.GetIntFromMap(info, "x")
	f.y = util.GetIntFromMap(info, "y")
	f.fan = util.GetStringFromMap(info, "fan")
}

// Class returns the class of this Pile
func (f *Foundation) Class() string {
	return f.class
}

// Pop here does not allow cards to be moved from Foundation
func (f *Foundation) Pop() *Card {
	return nil
}
