package sol

// Waste is the destination for cards
type Waste struct {
	Pile

	Class string
}

// NewWaste creates a new Waste
func NewWaste(x, y int) *Waste {
	w := &Waste{Pile: Pile{X: x, Y: y}, Class: "Waste"}
	return w
}
