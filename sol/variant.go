package sol

// VariantInfo contains configuration info for a variant
type VariantInfo struct {
	Description string
	AKA         []string
	Wikipedia   string
	Piles       []PileInfo
}

// PileInfo contains the basic Pile members and a lookup table for it's attributes
type PileInfo struct {
	Class      string
	X, Y       int
	Fan        string
	Attributes map[string]string
}

// Variants contains configuration info for all the variants
var Variants = map[string]VariantInfo{
	"Klondike": {
		Description: "The well-known solitaire variant",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "TapTarget": "Waste", "Recycles": "9999"}},
			{"Waste", 2, 1, "Waste", map[string]string{}},
			{"Foundation", 4, 1, "None", map[string]string{"Accept": "1", "Build": "21"}},
			{"Foundation", 5, 1, "None", map[string]string{"Accept": "1", "Build": "21"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "u"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "du"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "ddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "dddu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "ddddu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "dddddu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Move": "42", "Deal": "ddddddu"}},
		},
	},
}

func buildVariantPiles(v string) ([]*Pile, bool) {
	var piles []*Pile
	if vi, exists := Variants[v]; exists {
		for _, pi := range vi.Piles {
			p := NewPile(pi.Class, pi.X, pi.Y, pi.Fan, pi.Attributes)
			piles = append(piles, p)
		}
	}

	return piles, len(piles) > 0
}
