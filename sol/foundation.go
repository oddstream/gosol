package sol

// Foundation is the destination for cards
type Foundation struct {
	Pile
}

// NewFoundation creates a new Foundation
func NewFoundation(x, y int) *Foundation {
	f := &Foundation{Pile{X: x, Y: y}}
	return f
}

// Pop here does not allow cards to be moved from Foundation
func (f *Foundation) Pop() *Card {
	return nil
}
