package sol

// Tableau is the destination for cards
type Tableau struct {
	Pile

	Class string
}

// NewTableau creates a new Tableau
func NewTableau(x, y int) *Tableau {
	t := &Tableau{Pile: Pile{x: x, y: y, fan: "Down"}, Class: "Tableau"}
	return t
}
