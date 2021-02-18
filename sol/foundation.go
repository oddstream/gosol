package sol

// Foundation is the destination for cards
type Foundation struct {
	Pile

	Class string
}

// NewFoundation creates a new Foundation
func NewFoundation(x, y int) *Foundation {
	f := &Foundation{Pile: Pile{X: x, Y: y}, Class: "Foundation"}
	return f
}

// Pop here does not allow cards to be moved from Foundation
func (f *Foundation) Pop() *Card {
	return nil
}
