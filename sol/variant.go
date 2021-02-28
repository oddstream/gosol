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
	"Australian": {
		Description: "A combination of Klondike and Scorpion",
		AKA:         []string{"Australian Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Australian_Patience",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "0", "Build": "0", "Drag": "15"}},
			{"Waste", 2, 1, "Waste", map[string]string{"Build": "15", "Drag": "15"}},
			{"Foundation", 4, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 5, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Accept": "13", "Build": "22", "Drag": "15", "Deal": "uuuu"}},
		},
	},
	"Canfield": {
		Description: "Canfield has a 1 in 30 chance of winning. According to legend, it is originally a casino game, named after the casino owner who is said to have invented it.",
		AKA:         []string{"Demon"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "Build": "0", "Drag": "15", "CardsToMove": "3"}},
			{"Waste", 2, 1, "Waste", map[string]string{"Build": "15", "Drag": "15"}},
			{"Foundation", 4, 1, "None", map[string]string{"Build": "121", "Drag": "0", "AcceptFirstPush": "True", "Deal": "u"}},
			{"Foundation", 5, 1, "None", map[string]string{"Build": "121", "Drag": "0", "AcceptFirstPush": "True"}},
			{"Foundation", 6, 1, "None", map[string]string{"Build": "121", "Drag": "0", "AcceptFirstPush": "True"}},
			{"Foundation", 7, 1, "None", map[string]string{"Build": "121", "Drag": "0", "AcceptFirstPush": "True"}},
			{"Reserve", 1, 2, "Down", map[string]string{"Build": "0", "Drag": "15", "Deal": "ddddddddddddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Build": "142", "Drag": "42", "AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Build": "142", "Drag": "42", "AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Build": "142", "Drag": "42", "AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Build": "142", "Drag": "42", "AutoFillFrom": "Reserve", "Deal": "u"}},
		},
	},
	"EasyWin": {
		Description: "A game for debugging",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "Build": "0", "Drag": "15"}},
			{"Waste", 2, 1, "Waste", map[string]string{"Build": "15", "Drag": "15"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 8, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 9, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 8, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 9, 2, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},

			{"Tableau", 1, 5, "Down", map[string]string{"Build": "42", "Drag": "42"}},
			{"Tableau", 2, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 3, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 4, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 5, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 6, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 7, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "du"}},
			{"Tableau", 8, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
			{"Tableau", 9, 5, "Down", map[string]string{"Build": "42", "Drag": "42", "Deal": "ddu"}},
		},
	},
	"Freecell": {
		Description: "Popular game, unusual because almost all deals are winnable.",
		Wikipedia:   "https://en.wikipedia.org/wiki/FreeCell",
		Piles: []PileInfo{
			{"Stock", 4, -1, "None", map[string]string{"Packs": "1", "Build": "0", "Drag": "0"}},
			{"Cell", 1, 1, "None", map[string]string{"Accept": "0", "Build": "15", "Drag": "15"}},
			{"Cell", 2, 1, "None", map[string]string{"Accept": "0", "Build": "15", "Drag": "15"}},
			{"Cell", 3, 1, "None", map[string]string{"Accept": "0", "Build": "15", "Drag": "15"}},
			{"Cell", 4, 1, "None", map[string]string{"Accept": "0", "Build": "15", "Drag": "15"}},
			{"Foundation", 5, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 6, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 7, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Foundation", 8, 1, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuuu"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuuu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuuu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuuu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuu"}},
			{"Tableau", 8, 2, "Down", map[string]string{"Accept": "0", "Build": "42", "Drag": "242", "Deal": "uuuuuu"}},
		},
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
	"Klondike3": {
		Description: "The well-known solitaire variant",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 1, 1, "None", map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "Build": "0", "Drag": "15", "CardsToMove": "3"}},
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
			{"Tableau", 1, 2, "Down", map[string]string{"Build": "22", "Drag": "122", "Deal": "uuu"}},
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
	"Spider1": {
		Wikipedia: "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 1, 1, "None", map[string]string{"Packs": "8", "Build": "0", "Drag": "15", "Suits": "Spade"}},
			{"FoundationSpider", 3, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 4, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 5, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 6, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 7, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 8, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 9, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 10, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 8, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 9, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 10, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
		},
	},
	"Spider2": {
		Wikipedia: "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 1, 1, "None", map[string]string{"Packs": "4", "Build": "0", "Drag": "15", "Suits": "Spade,Heart"}},
			{"FoundationSpider", 3, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 4, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 5, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 6, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 7, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 8, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 9, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"FoundationSpider", 10, 1, "None", map[string]string{"Accept": "13", "Build": "22", "Drag": "0"}},
			{"Tableau", 1, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 2, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 3, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 4, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "dddddu"}},
			{"Tableau", 5, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 6, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 7, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 8, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 9, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
			{"Tableau", 10, 2, "Down", map[string]string{"Build": "12", "Drag": "22", "Deal": "ddddu"}},
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
