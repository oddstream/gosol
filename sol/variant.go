package sol

import (
	"log"
	"sort"
)

// VariantInfo contains configuration info for a variant
type VariantInfo struct {
	Description string
	AKA         []string
	Wikipedia   string
	PowerMoves  bool
	Piles       []PileInfo
}

// PileInfo contains the basic Pile members and a lookup table for it's attributes
type PileInfo struct {
	Class              string
	X, Y               PilePositionType
	Fan                string
	Build, Drag, Flags int
	Attributes         map[string]string
}

// Variants contains configuration info for all the variants
var Variants = map[string]VariantInfo{
	"American Toad": {
		Description: "This game is similar to Canfield except that the tableau builds down in suit, and a partial tableau stack cannot be moved (only the top card or entire stack can be moved).",
		Wikipedia:   "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
		Piles: []PileInfo{
			{"Foundation", 0, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True", "Deal": "u"}},
			{"Foundation", 1, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 2, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 3, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 4, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 5, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 6, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 7, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Tableau", 0, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 1, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 2, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 3, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 4, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 5, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 6, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 7, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Stock", 9, 0, "None", 0, 15, 0, map[string]string{"Packs": "2", "Target": "Waste", "Recycles": "1"}},
			{"Waste", 9, 1, "WasteDown", 15, 15, 1, nil},
			{"Reserve", 9, 3, "Down", 0, 15, 0, map[string]string{"Deal": "dddddddddddddddddddu"}},
		},
	},
	"Australian": {
		Description: "A combination of Klondike and Scorpion.",
		AKA:         []string{"Australian Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Australian_Patience",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "0"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 1, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 2, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 3, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 5, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 6, 1, "Down", 22, 15, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
		},
	},
	"Baker's Dozen": {
		Description: "The game is so called because of the 13 columns in the game, the number in a baker's dozen. Empty piles cannot be filled, so Kings are placed at the bottom of a pile during the initial dealing.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, -2, "None", 0, 15, 0, map[string]string{}},
			{"Tableau", 0, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 1, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 2, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 3, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 5, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 6, 0, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 0, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 1, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 2, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 3, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 5, 4, "Down", 12, 15, 1, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 1, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 2, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 3, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
		},
	},
	"Baker's Dozen (Relaxed)": {
		Description: "The game is so called because of the 13 columns in the game, the number in a baker's dozen. Empty piles cannot be filled, so Kings are placed at the bottom of a pile during the initial dealing.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, -2, "None", 0, 15, 0, map[string]string{}},
			{"Tableau", 0, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 1, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 2, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 3, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 5, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 6, 0, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 0, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 1, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 2, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 3, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Tableau", 5, 4, "Down", 12, 22, 0, map[string]string{"Accept": "99", "Bury": "13", "Deal": "uuuu"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 1, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 2, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 3, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
		},
	},
	"Canfield": {
		Description: "Canfield has a 1 in 30 chance of winning. According to legend, it is originally a casino game, named after the casino owner who is said to have invented it.",
		AKA:         []string{"Demon"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "CardsToMove": "3"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"AcceptFirstPush": "True", "Deal": "u"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"AcceptFirstPush": "True"}},
			{"Reserve", 0, 1, "Down", 0, 15, 0, map[string]string{"Deal": "ddddddddddddu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 1, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 4, 1, "Down", 42, 42, 1, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 5, 1, "Down", 42, 42, 1, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 6, 1, "Down", 42, 42, 1, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
		},
	},
	"EasyWin": {
		Description: "A game for debugging.",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "1"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1", "Deal": "1"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1", "Deal": "1"}},
			{"Foundation", 9, 0, "None", 21, 0, 0, map[string]string{"Accept": "1", "Deal": "1"}},
			{"Foundation", 10, 0, "None", 21, 0, 0, map[string]string{"Accept": "1", "Deal": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 7, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 8, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 9, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 10, 1, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 0, 4, "Down", 42, 42, 0, map[string]string{"Deal": "du"}},
			{"Tableau", 1, 4, "Down", 42, 42, 0, map[string]string{"Deal": "du"}},
			{"Tableau", 2, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 3, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 4, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 5, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 6, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 7, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 8, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 9, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
			{"Tableau", 10, 4, "Down", 42, 42, 0, map[string]string{"Deal": "uu"}},
		},
	},
	"Forty And Eight": {
		Description: "A variation of Forty Thieves that allows the stock to be recycled once.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(card_game)",
		PowerMoves:  true,
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "2", "Target": "Waste", "Recycles": "1"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 9, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 10, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 3, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 4, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 5, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 6, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 7, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 8, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 9, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 10, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuuuu"}},
		},
	},
	"Freecell": {
		Description: "Popular game, unusual because almost all deals are winnable.",
		Wikipedia:   "https://en.wikipedia.org/wiki/FreeCell",
		PowerMoves:  true,
		Piles: []PileInfo{
			{"Stock", -2, -2, "None", 0, 0, 0, map[string]string{"Packs": "1"}},
			{"Cell", 0, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 1, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 2, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 3, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuuu"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuuu"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuuu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuuu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuu"}},
			{"Tableau", 7, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Deal": "uuuuuu"}},
		},
	},
	"Freecell (Easy)": {
		Description: "Popular game, unusual because almost all deals are winnable.",
		Wikipedia:   "https://en.wikipedia.org/wiki/FreeCell",
		PowerMoves:  true,
		Piles: []PileInfo{
			{"Stock", -2, -2, "None", 0, 0, 0, map[string]string{"Packs": "1"}},
			{"Cell", 0, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 1, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 2, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Cell", 3, 0, "None", 15, 15, 1, map[string]string{"Accept": "0"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuuu"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuuu"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuuu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuuu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuu"}},
			{"Tableau", 7, 1, "Down", 42, 42, 0, map[string]string{"Accept": "0", "Bury": "13", "Disinter": "1", "Deal": "uuuuuu"}},
		},
	},
	"Giant": {
		Description: "An easier version of Miss Milligan, with no reserve pile to waive or weave onto.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Miss_Milligan",
		Piles: []PileInfo{
			{"StockScorpion", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "2"}},
			{"Foundation", 2, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 9, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 7, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 8, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 9, 1, "Down", 42, 42, 0, map[string]string{"Deal": "u"}},
		},
	},
	// // "King Albert": {
	// // 	Description: "Like Klondike, but with all cards face up and a seven-card reserve instead of stock and waste piles. It is the best known of the three games that are each called Idiot's Delight because of the low chance of winning the game.",
	// // 	AKA:         []string{"Idiot's Delight"},
	// // 	Wikipedia:   "https://en.wikipedia.org/wiki/King_Albert_(solitaire)",
	// // 	Piles: []PileInfo{
	// // 		{"Stock", -2, -2, "None", map[string]string{"Packs": "1", "Build": "0", "Drag": "15"}},
	// // 		{"Reserve", 0, 0, "Right", map[string]string{"Invisible": "True", "Build": "15", "Drag": "115", "Deal": "uuuuuuu"}},
	// // 		{"Foundation", 5, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
	// // 		{"Foundation", 6, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
	// // 		{"Foundation", 7, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
	// // 		{"Foundation", 8, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0"}},
	// // 		{"Tableau", 0, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "u"}},
	// // 		{"Tableau", 1, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uu"}},
	// // 		{"Tableau", 2, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuu"}},
	// // 		{"Tableau", 3, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuu"}},
	// // 		{"Tableau", 4, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuu"}},
	// // 		{"Tableau", 5, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuu"}},
	// // 		{"Tableau", 6, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuu"}},
	// // 		{"Tableau", 7, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuuu"}},
	// // 		{"Tableau", 8, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuuuu"}},
	// // 	},
	// // },
	"Klondike": {
		Description: "The well-known solitaire variant.",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "u"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "du"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "dddu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "dddddu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddddddu"}},
		},
	},
	"Klondike (Draw Three)": {
		Description: "The well-known solitaire variant.",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999", "CardsToMove": "3"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "u"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "du"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "dddu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "dddddu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "ddddddu"}},
		},
	},
	"Limited": {
		Description: "A more balanced version of Forty Thieves, with a wider tableaux.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(card_game)",
		PowerMoves:  true,
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "2", "Target": "Waste", "Recycles": "0"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 7, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 9, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 10, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 11, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 1, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 2, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 3, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 4, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 5, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 6, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 7, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 8, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 9, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 10, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
			{"Tableau", 11, 1, "Down", 22, 22, 0, map[string]string{"Deal": "uuu"}},
		},
	},
	// // "Raglan": {
	// // 	Description: "A slightly easier version of King Albert.",
	// // 	AKA:         []string{"Idiot's Delight"},
	// // 	Wikipedia:   "https://en.wikipedia.org/wiki/King_Albert_(solitaire)",
	// // 	Piles: []PileInfo{
	// // 		{"Stock", -2, -2, "None", map[string]string{"Packs": "1", "Build": "0", "Drag": "15"}},
	// // 		{"Reserve", 0, 0, "Right", map[string]string{"Invisible": "True", "Build": "15", "Drag": "115", "Deal": "uuu"}},
	// // 		{"Foundation", 5, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0", "Deal": "1"}},
	// // 		{"Foundation", 6, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0", "Deal": "1"}},
	// // 		{"Foundation", 7, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0", "Deal": "1"}},
	// // 		{"Foundation", 8, 0, "None", map[string]string{"Accept": "1", "Build": "21", "Drag": "0", "Deal": "1"}},
	// // 		{"Tableau", 0, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "u"}},
	// // 		{"Tableau", 1, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uu"}},
	// // 		{"Tableau", 2, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuu"}},
	// // 		{"Tableau", 3, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuu"}},
	// // 		{"Tableau", 4, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuu"}},
	// // 		{"Tableau", 5, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuu"}},
	// // 		{"Tableau", 6, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuu"}},
	// // 		{"Tableau", 7, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuuu"}},
	// // 		{"Tableau", 8, 1, "Down", map[string]string{"Build": "42", "Drag": "142", "Deal": "uuuuuuuuu"}},
	// // 	},
	// // },
	"Scorpion": {
		Description: "Related to Spider, with a method of game play like Yukon.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Scorpion_(solitaire)",
		Piles: []PileInfo{
			{"StockScorpion", 0, 0, "None", 0, 15, 0, nil},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "uuuuuuu"}},
			{"Tableau", 1, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "ddduuuu"}},
			{"Tableau", 2, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "ddduuuu"}},
			{"Tableau", 3, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "ddduuuu"}},
			{"Tableau", 4, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "uuuuuuu"}},
			{"Tableau", 5, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "uuuuuuu"}},
			{"Tableau", 6, 1, "Down", 22, 15, 8, map[string]string{"Accept": "13", "Deal": "uuuuuuu"}},
		},
	},
	"Simple Simon": {
		Description: "A wonderfully simple game with no stock or waste, that plays like Spider. Most games are winnable, but require skill.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)",
		Piles: []PileInfo{
			{"Stock", -2, -2, "None", 0, 0, 0, nil},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuuuuu"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuuuuu"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuuuuu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuuuu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuuu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuuu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuuu"}},
			{"Tableau", 7, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uuu"}},
			{"Tableau", 8, 1, "Down", 12, 22, 8, map[string]string{"Deal": "uu"}},
			{"Tableau", 9, 1, "Down", 12, 22, 8, map[string]string{"Deal": "u"}},
		},
	},
	"Spider": {
		Description: "The game originates in 1949, and its name comes from a spider's eight legs, referencing the eight foundation piles that must be filled to win the game.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "2"}},
			{"Foundation", 2, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 7, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 8, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 9, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 7, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 8, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 9, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
		},
	},
	"Spider (One Suit)": {
		Description: "The game originates in 1949, and its name comes from a spider's eight legs, referencing the eight foundation piles that must be filled to win the game.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "8", "Suits": "Spade"}},
			{"Foundation", 2, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 7, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 8, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 9, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 7, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 8, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 9, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
		},
	},
	"Spider (Two Suits)": {
		Description: "The game originates in 1949, and its name comes from a spider's eight legs, referencing the eight foundation piles that must be filled to win the game.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "4", "Suits": "Spade,Heart"}},
			{"Foundation", 2, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 7, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 8, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 9, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 7, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 8, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 9, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
		},
	},
	"Spiderette": {
		Description: "A compact version of Spider.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "4", "Suits": "Spade"}},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "u"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "du"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddddddu"}},
		},
	},
	"Storehouse": {
		Description: "An easier version of Canfield.",
		AKA:         []string{"Thirteen Up"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "1", "CardsToMove": "1"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 1, map[string]string{"AcceptFirstPush": "True", "Deal": "2"}},
			{"Foundation", 4, 0, "None", 21, 0, 1, map[string]string{"AcceptFirstPush": "True", "Deal": "2"}},
			{"Foundation", 5, 0, "None", 21, 0, 1, map[string]string{"AcceptFirstPush": "True", "Deal": "2"}},
			{"Foundation", 6, 0, "None", 21, 0, 1, map[string]string{"AcceptFirstPush": "True", "Deal": "2"}},
			{"Reserve", 0, 1, "Down", 0, 15, 1, map[string]string{"Deal": "ddddddddddddu"}},
			{"Tableau", 3, 1, "Down", 22, 22, 0, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 4, 1, "Down", 22, 22, 0, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 5, 1, "Down", 22, 22, 0, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
			{"Tableau", 6, 1, "Down", 22, 22, 0, map[string]string{"AutoFillFrom": "Reserve", "Deal": "u"}},
		},
	},
	"Thoughtful": {
		Description: "Klondike, but with all the cards face up.",
		AKA:         []string{"Patience", "American Patience", "Fascination", "Triangle", "Demon Patience"},
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999"}},
			{"Waste", 1, 0, "Waste", 15, 15, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "u"}},
			{"Tableau", 1, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uu"}},
			{"Tableau", 2, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uuu"}},
			{"Tableau", 3, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uuuu"}},
			{"Tableau", 4, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uuuuu"}},
			{"Tableau", 5, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uuuuuu"}},
			{"Tableau", 6, 1, "Down", 42, 42, 0, map[string]string{"Accept": "13", "Deal": "uuuuuuu"}},
		},
	},
	"The Toad": {
		Description: "As described in the 1908 by Hapgood, it is similar to Canfield except that the tableau builds down in suit, and a partial tableau stack cannot be moved (only the top card or entire stack can be moved). With two passes through the stock, it's hard to imagine how you could lose a game.",
		Wikipedia:   "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "2", "Target": "Waste", "Recycles": "1"}},
			{"Waste", 0, 1, "WasteDown", 15, 15, 1, map[string]string{}},
			{"Reserve", 0, 4, "None", 0, 15, 0, map[string]string{"Deal": "uuuuuuuuuuuuuuuuuuuu"}},
			{"Foundation", 2, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True", "Deal": "u"}},
			{"Foundation", 3, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 4, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 5, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 6, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 7, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 8, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Foundation", 9, 0, "None", 21, 0, 4, map[string]string{"AcceptFirstPush": "True"}},
			{"Tableau", 2, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 3, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 4, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 5, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 6, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 7, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 8, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
			{"Tableau", 9, 1, "Down", 22, 22, 6, map[string]string{"AutoFillFrom": "Reserve,Waste", "Deal": "u"}},
		},
	},
	"Thumb and Pouch": {
		Description: "An easy variant of Klondike, where a card in the tableau can be built upon another that is any suit other than its own (e.g. spades cannot be placed over spades) and spaces can be filled by any card or sequence.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Klondike_(solitaire)#Variations",
		Piles: []PileInfo{
			{"Stock", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "0"}},
			{"Waste", 1, 0, "Waste", 15, 0, 1, nil},
			{"Foundation", 3, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 4, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 5, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 6, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Tableau", 0, 1, "Down", 52, 52, 0, map[string]string{"Deal": "u"}},
			{"Tableau", 1, 1, "Down", 52, 52, 0, map[string]string{"Deal": "du"}},
			{"Tableau", 2, 1, "Down", 52, 52, 0, map[string]string{"Deal": "ddu"}},
			{"Tableau", 3, 1, "Down", 52, 52, 0, map[string]string{"Deal": "dddu"}},
			{"Tableau", 4, 1, "Down", 52, 52, 0, map[string]string{"Deal": "ddddu"}},
			{"Tableau", 5, 1, "Down", 52, 52, 0, map[string]string{"Deal": "dddddu"}},
			{"Tableau", 6, 1, "Down", 52, 52, 0, map[string]string{"Deal": "ddddddu"}},
		},
	},
	"Wasp": {
		Description: "An easier version of Scorpion, related to Spider, with a method of game play like Yukon.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Scorpion_(solitaire)",
		Piles: []PileInfo{
			{"StockScorpion", 0, 0, "None", 0, 15, 0, nil},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 22, 15, 8, map[string]string{"Deal": "ddduuuu"}},
			{"Tableau", 1, 1, "Down", 22, 15, 8, map[string]string{"Deal": "ddduuuu"}},
			{"Tableau", 2, 1, "Down", 22, 15, 8, map[string]string{"Deal": "ddduuuu"}},
			{"Tableau", 3, 1, "Down", 22, 15, 8, map[string]string{"Deal": "uuuuuuu"}},
			{"Tableau", 4, 1, "Down", 22, 15, 8, map[string]string{"Deal": "uuuuuuu"}},
			{"Tableau", 5, 1, "Down", 22, 15, 8, map[string]string{"Deal": "uuuuuuu"}},
			{"Tableau", 6, 1, "Down", 22, 15, 8, map[string]string{"Deal": "uuuuuuu"}},
		},
	},
	"Will o' the Wisp": {
		Description: "Invented by Geoffrey Mott-Smith, it is played the same way as Spiderette.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		Piles: []PileInfo{
			{"StockSpider", 0, 0, "None", 0, 15, 0, map[string]string{"Packs": "4", "Suits": "Spade"}},
			{"Foundation", 3, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 4, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 5, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Foundation", 6, 0, "None", 22, 0, 8, map[string]string{"Accept": "13"}},
			{"Tableau", 0, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 1, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 2, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 3, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 4, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 5, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
			{"Tableau", 6, 1, "Down", 12, 22, 8, map[string]string{"Deal": "ddu"}},
		},
	},
	"Yukon": {
		Description: "Like Klondike, but with no stock or waste piles.",
		Wikipedia:   "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		Piles: []PileInfo{
			{"Stock", -2, -2, "None", 0, 15, 0, map[string]string{"Packs": "1", "Target": "Waste", "Recycles": "9999"}},
			{"Tableau", 0, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "u"}},
			{"Tableau", 1, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "duuuuu"}},
			{"Tableau", 2, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "dduuuuu"}},
			{"Tableau", 3, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "ddduuuuu"}},
			{"Tableau", 4, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "dddduuuuu"}},
			{"Tableau", 5, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "ddddduuuuu"}},
			{"Tableau", 6, 0, "Down", 42, 15, 0, map[string]string{"Accept": "13", "Deal": "dddddduuuuu"}},
			{"Foundation", 8, 0, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 1, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 2, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
			{"Foundation", 8, 3, "None", 21, 0, 0, map[string]string{"Accept": "1"}},
		},
	},
}

// BuildVariant creates the Piles and loads any game-wide flags
// we already know that the variant exists
func (b *Baize) BuildVariant(v string) {
	b.Piles = nil
	if vi, exists := Variants[v]; exists {
		for _, pi := range vi.Piles {
			p := NewPile(pi.Class, pi.X, pi.Y, pi.Fan, pi.Build, pi.Drag, pi.Flags, pi.Attributes)
			b.Piles = append(b.Piles, p)
		}
	} else {
		log.Fatal("unknown variant ", v)
	}
	b.PowerMoves = Variants[v].PowerMoves
	for _, p1 := range b.Piles {
		if p1.Class != "Tableau" {
			continue
		}
		var maxX, maxY PilePositionType
		for _, p2 := range b.Piles {
			if p2.Class != "Tableau" {
				continue
			}
			if p1 == p2 {
				continue
			}
			if p2.Y == p1.Y && p2.X > p1.X {
				maxX = p2.X
			}
			if p2.X == p1.X && p2.Y > p1.Y {
				maxY = p2.Y
			}
		}
		switch p1.Fan {
		case "Right":
			if maxX == 0 {
				p1.scrunchSize = 6
			} else {
				p1.scrunchSize = int(maxX - p1.X)
			}
			// println("Right", p1.scrunchSize)
		case "Down":
			if maxY == 0 {
				p1.scrunchSize = 6
			} else {
				p1.scrunchSize = int(maxY - p1.Y)
			}
			// println("Down", p1.scrunchSize)
		}
	}
}

func variantDescription(v string) string {
	if vi, exists := Variants[v]; exists {
		return vi.Description
	}
	return ""
}

func (b *Baize) ShowVariantPicker() {
	var vnames []string
	for key := range Variants {
		vnames = append(vnames, key)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	b.ui.ShowVariantPicker(vnames)
}
