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
	"Freecell": {
		Description: "Popular game, unusual because almost all deals are winnable.",
		Wikipedia:   "https://en.wikipedia.org/wiki/FreeCell",
		Piles:       []PileInfo{},
	},
	"Klondike": {
		Description: "The well-known solitaire variant",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "Build": "0", "Drag": "15"}},
			{"Waste", 2, 1, "Waste", map[string]string{"Build": "15", "Drag": "15"}},
			{"Foundation", 4, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 5, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "u"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "dddu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "ddddu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "dddddu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Accept": "13", "Build": "42", "Drag": "42", "Deal": "ddddddu"}},
		},
	},
	"Limited": {
		Description: "A more balanced version of Forty Thieves, with a wider tableaux.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(card_game)",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "2", "Target": "Waste", "Recycles": "0", "Build": "0", "Drag": "15"}},
			{"Waste", 2, 1, "None", map[string]string{"Build": "15", "Drag": "15"}},
			{"Foundation", 5, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 8, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 9, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 10, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 11, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 12, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Build": "22", "Drag": "122"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 8, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 9, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 10, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 11, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
			{"Tableau", 12, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
		},
	},
	"Spider": {},
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
