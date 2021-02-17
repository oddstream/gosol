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
