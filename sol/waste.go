package sol

// Waste is the destination for cards
type Waste struct {
	Pile

	Class string
}

// WasteInfo contains configuration for all Waste objects
type WasteInfo struct {
	// No additional members; X, Y, Fan will do it
}

// NewWaste creates a new Waste
func NewWaste(x, y int) *Waste {
	w := &Waste{Pile: Pile{x: x, y: y, fan: "Right"}, Class: "Waste"}
	return w
}
