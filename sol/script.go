package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"sort"

	"oddstream.games/gomps5/util"
)

type ScriptBase struct {
	stock, waste                                     *Pile
	cells, discards, foundations, reserves, tableaux []*Pile
	easy                                             bool
}

type VariantInfo struct {
	windowShape string
	wikipedia   string
}

// You can't use functions as keys in maps : the key type must be comparable
// so you can't do: var ExtendedColorMap = map[CardPairCompareFunc]bool{}
// type CardPairCompareFunc func(CardPair) (bool, error)

type ScriptInterface interface {
	BuildPiles() VariantInfo
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(*Pile, []*Card) (bool, error)
	UnsortedPairs(*Pile) int

	TailTapped([]*Card)
	PileTapped(*Pile)

	Discards() []*Pile
	Foundations() []*Pile
	Stock() *Pile
	Waste() *Pile
}

var Variants = map[string]ScriptInterface{
	"Agnes Bernauer":      &Agnes{},
	"American Toad":       &Toad{},
	"Australian":          &Australian{},
	"Baker's Dozen":       &BakersDozen{},
	"Duchess":             &Duchess{},
	"Klondike":            &Klondike{draw: 1, recycles: 2},
	"Klondike Draw Three": &Klondike{draw: 3, recycles: 9},
	"Thoughtful":          &Klondike{draw: 1, recycles: 32767, thoughtful: true},
	"Easy":                &Easy{ScriptBase: ScriptBase{easy: true}},
	"Freecell":            &Freecell{},
	"Forty Thieves": &FortyThieves{
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4},
	"Limited": &FortyThieves{
		founds:      []int{5, 6, 7, 8, 9, 10, 11, 12},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		cardsPerTab: 3},
	"Penguin":           &Penguin{},
	"Scorpion":          &Scorpion{},
	"Simple Simon":      &SimpleSimon{},
	"Spider One Suit":   &Spider{packs: 8, suits: 1},
	"Spider Two Suits":  &Spider{packs: 4, suits: 2},
	"Spider Four Suits": &Spider{packs: 2, suits: 4},
	"Whitehead":         &Whitehead{},
	"Yukon":             &Yukon{},
	"Yukon Cells":       &Yukon{extraCells: 2},
}

var VariantGroups = map[string][]string{
	// "All" added dynamically by func init()
	// don't have Agnes here (as a group) because it would come before All
	// and Agnes Sorel might be retired because it's just too hard
	"> Klondike":      {"Klondike", "Klondike Draw Three", "Thoughtful", "Whitehead"},
	"> Forty Thieves": {"Forty Thieves", "Josephine", "Limited"},
	"> Spider":        {"Spider One Suit", "Spider Two Suits", "Spider Four Suits"},
	"> Scorpion":      {"Scorpion", "Wasp"},
	"> Canfield":      {"Canfield", "Acme", "Storehouse"},
	"> Freecell":      {"Freecell", "Eight Off"},
	"> Yukon":         {"Yukon", "Yukon Cells", "Alaska"},
	"> Puzzlers":      {"Penguin", "Simple Simon", "Baker's Dozen", "Freecell"},
}

func init() {
	var vnames []string
	for k := range Variants {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	VariantGroups["> All"] = vnames
}

func VariantGroupNames() []string {
	var vnames []string = make([]string, 0, len(VariantGroups))
	for k := range VariantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

func VariantNames(group string) []string {

	var vnames []string = make([]string, 0, len(VariantGroups[group]))
	vnames = append(vnames, VariantGroups[group]...)
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

// useful generic game library of functions

type CardPair struct {
	c1, c2 *Card
}

func (cp CardPair) EitherProne() bool {
	return cp.c1.Prone() || cp.c2.Prone()
}

type CardPairs []CardPair

func NewCardPairs(cards []*Card) []CardPair {
	if len(cards) < 2 {
		return []CardPair{}
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

func (cpairs CardPairs) Print() {
	for _, pair := range cpairs {
		println(pair.c1.String(), pair.c2.String())
	}
}

func (cp CardPair) Compare_Up() (bool, error) {
	if cp.c1.Ordinal()+1 != cp.c2.Ordinal() {
		return false, errors.New("Cards must be in ascending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_Down() (bool, error) {
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_DownColor() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownAltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownColorWrap() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
	}
	return true, nil
}

func (cp CardPair) Compare_DownAltColorWrap() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
	}
	return true, nil
}

func (cp CardPair) Compare_UpAltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_UpSuit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_DownSuit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpSuitWrap() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (cp CardPair) Compare_DownSuitWrap() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal()-1 == cp.c2.Ordinal() {
		return true, nil
	}
	return false, errors.New("Cards must go down in rank (Kings on Aces allowed)")
}

func Compare_Empty(p *Pile, c *Card) (bool, error) {

	if p.label != "" {
		if p.label == "x" {
			return false, errors.New("Cannot move cards there")
		}
		ord := util.OrdinalToShortString(c.Ordinal())
		if ord != p.label {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(p.label), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}
