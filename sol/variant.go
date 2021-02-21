package sol

// VariantInfo contains configuration info for a variant
type VariantInfo struct {
	Description string
	AKA         []string
	Wikipedia   string
	Piles       []PileInfo
}

// PileInfo contains the CardOwner object and a lookup table for it's attributes
type PileInfo struct {
	pile CardOwner
	info map[string]string
}

// Variants contains configuration info for all the variants
var Variants = map[string]VariantInfo{
	"Klondike": {
		Description: "The well-known solitaire variant",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{&Stock{}, map[string]string{"x": "1", "y": "1", "fan": "None", "Packs": "1", "TapTarget": "Waste", "Recycles": "9999"}},
			{&Waste{}, map[string]string{"x": "2", "y": "1", "fan": "Waste"}},
			{&Foundation{}, map[string]string{"x": "4", "y": "1", "fan": "None", "accept": "1"}},
			{&Foundation{}, map[string]string{"x": "5", "y": "1", "fan": "None", "accept": "1"}},
			{&Foundation{}, map[string]string{"x": "6", "y": "1", "fan": "None", "accept": "1"}},
			{&Foundation{}, map[string]string{"x": "7", "y": "1", "fan": "None", "accept": "1"}},
			{&Tableau{}, map[string]string{"x": "1", "y": "2", "fan": "Down", "accept": "13", "deal": "u"}},
			{&Tableau{}, map[string]string{"x": "2", "y": "2", "fan": "Down", "accept": "13", "deal": "du"}},
			{&Tableau{}, map[string]string{"x": "3", "y": "2", "fan": "Down", "accept": "13", "deal": "ddu"}},
			{&Tableau{}, map[string]string{"x": "4", "y": "2", "fan": "Down", "accept": "13", "deal": "dddu"}},
			{&Tableau{}, map[string]string{"x": "5", "y": "2", "fan": "Down", "accept": "13", "deal": "ddddu"}},
			{&Tableau{}, map[string]string{"x": "6", "y": "2", "fan": "Down", "accept": "13", "deal": "dddddu"}},
			{&Tableau{}, map[string]string{"x": "7", "y": "2", "fan": "Down", "accept": "13", "deal": "ddddddu"}},
		},
	},
}

func buildVariant(v string) ([]CardOwner, bool) {
	var owners []CardOwner
	if vi, exists := Variants[v]; exists {
		for _, pi := range vi.Piles {
			o := pi.pile
			o.New(pi.info)
			owners = append(owners, o)
		}
	}

	return owners, len(owners) > 0
}
