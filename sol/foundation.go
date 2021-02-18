package sol

// Foundation is the destination for cards
type Foundation struct {
	Pile

	Class string
}

// FoundationInfo contains configuration for all Waste objects
type FoundationInfo struct {
	Accept      int    // ordinal of card to accept on empty pile, 0 == any (FoundationSpider has it's own rules)
	DealOrdinal int    // ordinal of card to deal here, 0 == don't (why isn't this in Deal= parameter?)
	DealSuit    string // "" == any
}

// NewFoundation creates a new Foundation
func NewFoundation(x, y int) *Foundation {
	f := &Foundation{Pile: Pile{x: x, y: y, fan: "None"}, Class: "Foundation"}
	return f
}

// Pop here does not allow cards to be moved from Foundation
func (f *Foundation) Pop() *Card {
	return nil
}
