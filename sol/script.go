package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"sort"

	"oddstream.games/gosol/util"
)

type ScriptBase struct {
	cells       []*Pile
	discards    []*Pile
	foundations []*Pile
	reserves    []*Pile
	stock       *Pile
	tableaux    []*Pile
	waste       *Pile

	wikipedia  string
	cardColors int
}

// Complete - default is number of cards in Foundations == number of cards in CardLibrary.
//
// In Bisley, there may be <13 cards in a Foundation.
// This will need overriding for any variants with Discard piles.
// Could also do this by checking if any pile other than a Foundation is not empty.
func (sb ScriptBase) Complete() bool {
	var n = 0
	for _, f := range sb.foundations {
		n += len(f.cards)
	}
	return n == len(CardLibrary)
}

// SpiderComplete - used to override default Complete() in Spider varaints.
//
// Each tableau must be either empty or contain a sequence.
// Discard contents must be sequences, otherwise they wouldn't be there.
// There aren't any foundations.
func (sb ScriptBase) SpiderComplete() bool {
	for _, t := range sb.tableaux {
		switch len(t.cards) {
		case 0:
			// that's fine
		case 13: // TODO 104 cards, 8 tabs = 13 cards/tab
			if !t.vtable.Conformant() {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (sb ScriptBase) Wikipedia() string {
	if sb.wikipedia == "" { // uninitialized default
		return "https://en.wikipedia.org/wiki/Patience_(game)"
	} else {
		return sb.wikipedia
	}
}

func (sb ScriptBase) CardColors() int {
	if sb.cardColors == 0 { // uninitialized default
		return 2
	} else {
		return sb.cardColors
	}
}

func (sb ScriptBase) SafeCollect() bool {
	return sb.CardColors() == 2
}

func (sb ScriptBase) Cells() []*Pile {
	return sb.cells
}

func (sb ScriptBase) Foundations() []*Pile {
	return sb.foundations
}

func (sb ScriptBase) Discards() []*Pile {
	return sb.discards
}

func (sb ScriptBase) Reserves() []*Pile {
	return sb.reserves
}

func (sb ScriptBase) Stock() *Pile {
	return sb.stock
}

func (sb ScriptBase) Tableaux() []*Pile {
	return sb.tableaux
}

func (sb ScriptBase) Waste() *Pile {
	return sb.waste
}

func (sb ScriptBase) AfterMove() {}

// You can't use functions as keys in maps : the key type must be comparable
// so you can't do: var ExtendedColorMap = map[CardPairCompareFunc]bool{}
// type CardPairCompareFunc func(CardPair) (bool, error)

type Scripter interface {
	BuildPiles()
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(*Pile, []*Card) (bool, error)
	UnsortedPairs(*Pile) int

	TailTapped([]*Card)
	PileTapped(*Pile)

	Cells() []*Pile
	Discards() []*Pile
	Foundations() []*Pile
	Reserves() []*Pile
	Stock() *Pile
	Tableaux() []*Pile
	Waste() *Pile

	Complete() bool
	Wikipedia() string
	CardColors() int
	SafeCollect() bool
}

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
		},
	},
	"American Toad": &Toad{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
			cardColors: 4,
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
		},
		packs:    2,
		draw:     1,
		recycles: 2,
		founds:   []int{3, 4, 5, 6, 7, 8, 9, 10},    // 8
		tabs:     []int{2, 3, 4, 5, 6, 7, 8, 9, 10}, // 9
	},
	"Triple Klondike": &Klondike{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Gargantua_(card_game)",
		},
		packs:    3,
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
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
	},
	"Josephine": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
		moveType:    MOVE_ANY,
	},
	"Rank and File": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
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
		},
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Number Ten": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
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
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 3,
	},
	"Forty and Eight": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
			cardColors: 4,
		},
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 5,
		recycles:    1,
	},
	"Red and Black": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
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
		},
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 1,
	},
	"Maria": &FortyThieves{
		ScriptBase: ScriptBase{
			wikipedia: "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
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
		},
		packs:       3,
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		cardsPerTab: 5,
	},
	"Mrs Mop": &MrsMop{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
		},
	},
	"Mrs Mop Easy": &MrsMop{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Mrs._Mop",
			cardColors: 4,
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
		},
		packs: 8,
		suits: 1,
	},
	"Spider Two Suits": &Spider{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 2,
		},
		packs: 4,
		suits: 2,
	},
	"Spider Four Suits": &Spider{
		ScriptBase: ScriptBase{
			wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
			cardColors: 4,
		},
		packs: 2,
		suits: 4,
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
	"Oddstream": &Oddstream{
		ScriptBase: ScriptBase{
			wikipedia:  "https://oddstream.games",
			cardColors: 4,
		},
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
	var vnames []string = make([]string, 0, len(VariantGroups[group]))
	vnames = append(vnames, VariantGroups[group]...)
	if group == "> All by Played" {
		sort.Slice(vnames, func(i, j int) bool { return TheStatistics.Played(vnames[i]) > TheStatistics.Played(vnames[j]) })
	} else {
		sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	}
	return vnames
}

// useful generic game library of functions

func RecycleWasteToStock(waste *Pile, stock *Pile) {
	if TheBaize.Recycles() > 0 {
		for waste.Len() > 0 {
			MoveCard(waste, stock)
		}
		TheBaize.SetRecycles(TheBaize.Recycles() - 1)
		switch {
		case TheBaize.recycles == 0:
			TheUI.ToastInfo("No more recycles")
		case TheBaize.recycles == 1:
			TheUI.ToastInfo(fmt.Sprintf("%d recycle remaining", TheBaize.Recycles()))
		case TheBaize.recycles < 10:
			TheUI.ToastInfo(fmt.Sprintf("%d recycles remaining", TheBaize.Recycles()))
		}
	} else {
		TheUI.ToastInfo("No more recycles")
	}
}

func Compare_Empty(p *Pile, c *Card) (bool, error) {

	if p.Label() != "" {
		if p.Label() == "x" || p.Label() == "X" {
			return false, errors.New("Cannot move cards to that empty pile")
		}
		ord := util.OrdinalToShortString(c.Ordinal())
		if ord != p.Label() {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(p.Label()), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}

// UnsortedPairs
//
// A generic way of calculating the number of unsorted card pairs in a pile.
// Called by *Pile.vtable.UnsortedPairs()
func UnsortedPairs(pile *Pile, fn CardPairCompareFunc) int {
	if pile.Len() < 2 {
		return 0
	}
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.c1.Prone() || pair.c2.Prone() {
			unsorted++
		} else {
			if ok, _ := fn(pair); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

type CardPairCompareFunc func(CardPair) (bool, error)

type CardPair struct {
	c1, c2 *Card
}

type CardPairs []CardPair

func NewCardPairs(cards []*Card) CardPairs {
	if len(cards) < 2 {
		return []CardPair{} // always return a list, not nil
	}
	var cpairs []CardPair
	c1 := cards[0]
	for i := 1; i < len(cards); i++ {
		c2 := cards[i]
		cpairs = append(cpairs, CardPair{c1, c2})
		c1 = c2
	}
	return cpairs
}

// func (cpairs CardPairs) Print() {
// 	for _, pair := range cpairs {
// 		println(pair.c1.String(), pair.c2.String())
// 	}
// }

// little library of simple compares

func (cp CardPair) Compare_Up() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must be in ascending sequence")
}

func (cp CardPair) Compare_UpWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (cp CardPair) Compare_Down() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	return false, errors.New("Cards must be in descending sequence")
}

func (cp CardPair) Compare_DownWrap() (bool, error) {
	if cp.c1.Ordinal() == cp.c2.Ordinal()+1 {
		return true, nil
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
}

func (cp CardPair) Compare_UpOrDown() (bool, error) {
	if !(cp.c1.Ordinal()+1 == cp.c2.Ordinal() || cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_UpOrDownWrap() (bool, error) {
	if (cp.c1.Ordinal()+1 == cp.c2.Ordinal()) || (cp.c1.Ordinal() == cp.c2.Ordinal()+1) {
		return true, nil
	} else if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace On King
	} else if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	} else {
		return false, errors.New("Cards must be in ascending or descending sequence")
	}
}

func (cp CardPair) Compare_Color() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return true, nil
}

func (cp CardPair) Compare_AltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return true, nil
}

func (cp CardPair) Compare_Suit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return true, nil
}

func (cp CardPair) Compare_OtherSuit() (bool, error) {
	if cp.c1.Suit() == cp.c2.Suit() {
		return false, errors.New("Cards must not be the same suit")
	}
	return true, nil
}

// library of compare functions made from simple compares

func (cp CardPair) Compare_DownColor() (bool, error) {
	ok, err := cp.Compare_Color()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownAltColor() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

// Compare_DownColorWrap not used
func (cp CardPair) Compare_DownColorWrap() (bool, error) {
	ok, err := cp.Compare_Color()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

func (cp CardPair) Compare_DownAltColorWrap() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

// Compare_UpAltColor not used
func (cp CardPair) Compare_UpAltColor() (bool, error) {
	ok, err := cp.Compare_AltColor()
	if !ok {
		return ok, err
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_UpSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_DownSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpOrDownSuit() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpOrDown()
}

func (cp CardPair) Compare_UpOrDownSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpOrDownWrap()
}

// Compare_DownOtherSuit not used
func (cp CardPair) Compare_DownOtherSuit() (bool, error) {
	ok, err := cp.Compare_OtherSuit()
	if !ok {
		return ok, err
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_UpWrap()
}

func (cp CardPair) Compare_DownSuitWrap() (bool, error) {
	ok, err := cp.Compare_Suit()
	if !ok {
		return ok, err
	}
	return cp.Compare_DownWrap()
}

// ChainCall
//
// Call using CardPair method expressions
// eg ChainCall(CardPair.Compare_UpOrDown, CardPair.Compare_Suit)
//
// TODO think of something else for UnsortedPairs(*Pile)
func (cp CardPair) ChainCall(fns ...func(CardPair) (bool, error)) (ok bool, err error) {
	for _, fn := range fns {
		if ok, err = fn(cp); err != nil {
			break
		}
	}
	return
}

func TailConformant(tail []*Card, fn CardPairCompareFunc) (bool, error) {
	for _, pair := range NewCardPairs(tail) {
		if ok, err := fn(pair); !ok {
			return false, err
		}
	}
	return true, nil
}
