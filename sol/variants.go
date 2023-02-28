package sol

import "sort"

var Variants = map[string]Scripter{
	"Agnes Bernauer": &Agnes{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Agnes_(solitaire)",
			cardColors: 2,
		},
	},
	"Alhambra": &Alhambra{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Alhambra_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"American Toad": &Toad{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Antares": &Antares{
		ScriptBase: ScriptBase{
			wikipedia: "https://www.goodsol.com/games/antares.html",
		},
	},
	"Australian": &Australian{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Australian_Patience",
			cardColors: 4,
		},
	},
	"Baker's Dozen": &BakersDozen{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
			cardColors: 1,
		},
	},
	"Baker's Game": &Freecell{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Baker%27s_Game",
			cardColors: 4,
		},
		tabCompareFunc: CardPair.Compare_DownSuit,
	},
	"Bisley": &Bisley{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Bisley_(card_game)",
			cardColors: 4,
		},
	},
	"Blind Freecell": &Freecell{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/FreeCell",
			cardColors: 2,
		},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		blind:          true,
	},
	"Blockade": &Blockade{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Blockade_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
	},
	"Canfield": &Canfield{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		},
		draw:           3,
		recycles:       32767,
		tabCompareFunc: CardPair.Compare_DownAltColorWrap,
	},
	"Storehouse": &Canfield{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
			cardColors: 4,
		},
		draw:           1,
		recycles:       2,
		tabCompareFunc: CardPair.Compare_DownSuitWrap,
		variant:        "storehouse",
	},
	"Duchess": &Duchess{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Duchess_(solitaire)",
		},
	},
	"Demons and Thieves": &CanThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://www.goodsol.com/pgshelp/index.html?demons_and_thieves.htm",
			packs:     2,
		},
	},
	"Klondike": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:     1,
		recycles: 2,
	},
	"Klondike Draw Three": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:     3,
		recycles: 2,
	},
	"Thoughtful": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
		draw:       1,
		recycles:   2,
		thoughtful: true,
	},
	"Gargantua": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Gargantua_(card_game)",
			packs:     2,
		},
		draw:     1,
		recycles: 2,
		founds:   []int{3, 4, 5, 6, 7, 8, 9, 10},    // 8
		tabs:     []int{2, 3, 4, 5, 6, 7, 8, 9, 10}, // 9
	},
	"Triple Klondike": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Gargantua_(card_game)",
			packs:     3,
		},
		draw:     1,
		recycles: 2,
		founds:   []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},             // 12
		tabs:     []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, // 16
	},
	"Eight Off": &EightOff{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Eight_Off",
			cardColors: 4,
		},
	},
	"Freecell": &Freecell{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/FreeCell",
		},
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Freecell Easy": &Freecell{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/FreeCell",
		},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		easy:           true,
	},
	"Forty Thieves": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
	},
	"Josephine": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
		moveType:    MOVE_ANY,
	},
	"Rank and File": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1, 2},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Indian": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    3,
		proneRows:      []int{0},
		tabCompareFunc: CardPair.Compare_DownOtherSuit,
	},
	"Streets": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Number Ten": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Limited": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 3,
	},
	"Forty and Eight": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 5,
		recycles:    1,
	},
	"Red and Black": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Lucas": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{5, 6, 7, 8, 9, 10, 11, 12},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		cardsPerTab: 3,
		dealAces:    true,
	},
	"Busy Aces": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      2,
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 1,
	},
	"Maria": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			packs:     2,
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Sixty Thieves": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
			packs:      3,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		cardsPerTab: 5,
	},
	"Mrs Mop": &MrsMop{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
			packs:      2,
		},
	},
	"Mrs Mop Easy": &MrsMop{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
			packs:      2,
		},
		easy: true,
	},
	"Penguin": &Penguin{
		ScriptBase: ScriptBase{
			wikipedia:  "https://www.parlettgames.uk/patience/penguin.html",
			cardColors: 4,
		},
	},
	"Scorpion": &Scorpion{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Scorpion_(solitaire)",
			cardColors: 4,
		},
	},
	"Seahaven Towers": &Seahaven{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Seahaven_Towers",
			cardColors: 4,
		},
	},
	"Simple Simon": &SimpleSimon{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)",
			cardColors: 4,
		},
	},
	"Spider One Suit": &Spider{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 1,
			packs:      8,
			suits:      1,
		},
	},
	"Spider Two Suits": &Spider{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 2,
			packs:      4,
			suits:      2,
		},
	},
	"Spider Four Suits": &Spider{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 4,
			packs:      2,
			suits:      4,
		},
	},
	"Spiderette": &Spiderette{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)#Variants",
			cardColors: 4,
		},
	},
	"Classic Westcliff": &Westcliff{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "Classic",
	},
	"American Westcliff": &Westcliff{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "American",
	},
	"Easthaven": &Westcliff{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Westcliff_(card_game)",
		},
		variant: "Easthaven",
	},
	"Whitehead": &Whitehead{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
		},
	},
	"Usk": &Usk{
		ScriptBase: ScriptBase{
			wikipedia: "https://politaire.com/help/usk",
		},
		tableauLabel: "K",
	},
	"Usk Relaxed": &Usk{
		ScriptBase: ScriptBase{
			wikipedia: "https://politaire.com/help/usk",
		},
	},
	"Yukon": &Yukon{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		},
	},
	"Yukon Cells": &Yukon{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		},
		extraCells: 2,
	},
}

var VariantGroups = map[string][]string{
	// "> All" added dynamically by func init()
	// don't have any group that comes alphabetically before "> All"
	"> Canfields":     {"Canfield", "Storehouse", "Duchess", "American Toad"},
	"> Easier":        {"American Toad", "American Westcliff", "Blockade", "Classic Westcliff", "Lucas", "Spider One Suit", "Usk Relaxed"},
	"> Harder":        {"Baker's Dozen", "Easthaven", "Forty Thieves", "Spider Four Suits", "Usk"},
	"> Forty Thieves": {"Forty Thieves", "Number Ten", "Red and Black", "Indian", "Rank and File", "Sixty Thieves", "Josephine", "Limited", "Forty and Eight", "Lucas", "Busy Aces", "Maria", "Streets"},
	"> Freecells":     {"Baker's Game", "Blind Freecell", "Freecell", "Freecell Easy", "Eight Off", "Seahaven Towers"},
	"> Klondikes":     {"Gargantua", "Triple Klondike", "Klondike", "Klondike Draw Three", "Thoughtful", "Whitehead"},
	"> People":        {"Agnes Bernauer", "Duchess", "Josephine", "Maria", "Simple Simon", "Baker's Game"},
	"> Places":        {"Australian", "Bisley", "Yukon", "Klondike", "Usk", "Usk Relaxed"},
	"> Puzzlers":      {"Antares", "Demons and Thieves", "Bisley", "Usk", "Mrs Mop", "Penguin", "Simple Simon", "Baker's Dozen"},
	"> Spiders":       {"Spider One Suit", "Spider Two Suits", "Spider Four Suits", "Scorpion", "Spiderette"},
	"> Yukons":        {"Yukon", "Yukon Cells"},
}

// init is used to assemble the "> All" alpha-sorted group of variants for the picker menu
func init() {
	var vnames []string = make([]string, 0, len(Variants))
	for k := range Variants {
		vnames = append(vnames, k)
	}
	// no need to sort here, sort gets done by func VariantNames()
	VariantGroups["> All"] = vnames
	VariantGroups["> All by Played"] = vnames
}

// VariantGroupNames returns an alpha-sorted []string of the variant group names
func VariantGroupNames() []string {
	var vnames []string = make([]string, 0, len(VariantGroups))
	for k := range VariantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

// VariantNames returns an alpha-sorted []string of the variants in a group
func VariantNames(group string) []string {
	var vnames []string = nil
	vnames = append(vnames, VariantGroups[group]...)
	if group == "> All by Played" {
		sort.Slice(vnames, func(i, j int) bool {
			return TheGame.Statistics.Played(vnames[i]) > TheGame.Statistics.Played(vnames[j])
		})
	} else {
		sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	}
	return vnames
}
