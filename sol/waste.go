package sol

// Waste is the destination for cards
type Waste struct {
	Pile
}

// NewWaste creates a new Waste
func NewWaste(x, y int) *Waste {
	w := &Waste{Pile{X: x, Y: y}}
	return w
}
