package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"log"
	"sort"

	"oddstream.games/gomps5/util"
)

type ScriptInterface interface {
	BuildPiles()
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(Pile, []*Card) (bool, error)
	UnsortedPairs(Pile) int

	TailTapped([]*Card)
	PileTapped(Pile)

	PercentComplete() int
	Wikipedia() string
}

var Variants = map[string]ScriptInterface{
	"Klondike":            &Klondike{draw: 1, recycles: 2},
	"Klondike Draw Three": &Klondike{draw: 3, recycles: 9},
	"Easy":                &Easy{},
	"Freecell":            &Freecell{},
	"Simple Simon":        &SimpleSimon{},
	"Spider One Suit":     &Spider{packs: 8, suits: 1},
	"Spider Two Suits":    &Spider{packs: 4, suits: 2},
	"Spider":              &Spider{packs: 2, suits: 4},
}

func GetVariantInterface(v string) ScriptInterface {
	si, ok := Variants[v]
	if !ok {
		log.Panicf("Unknown variant %s", v)
	}
	return si
}

func VariantNames() []string {
	var vnames []string = make([]string, 0, len(Variants))
	for k := range Variants {
		vnames = append(vnames, k)
	}
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

func Script_PercentComplete() int {
	var pairs, unsorted, percent int
	for _, p := range TheBaize.piles {
		if p.Len() > 1 {
			pairs += p.Len() - 1
		}
		unsorted += p.UnsortedPairs()
	}
	percent = (int)(100.0 - util.MapValue(float64(unsorted), 0, float64(pairs), 0.0, 100.0))
	return percent
}

func CardCompare_Up(c1, c2 *Card) (bool, error) {
	if c1.Ordinal()+1 != c2.Ordinal() {
		return false, errors.New("Cards must be in ascending order")
	}
	return true, nil
}

func CardCompare_Down(c1, c2 *Card) (bool, error) {
	if c1.Ordinal() != c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending order")
	}
	return true, nil
}

func CardCompare_DownAltColor(c1, c2 *Card) (bool, error) {
	if c1.Black() == c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return CardCompare_Down(c1, c2)
}

func CardCompare_DownAltColorWrap(c1, c2 *Card) (bool, error) {
	if c1.Black() == c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	if c1.Ordinal() == 1 && c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	return CardCompare_Down(c1, c2)
}

func CardCompare_UpAltColor(c1, c2 *Card) (bool, error) {
	if c1.Black() == c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return CardCompare_Up(c1, c2)
}

func CardCompare_UpSuit(c1, c2 *Card) (bool, error) {
	if c1.Suit() != c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return CardCompare_Up(c1, c2)
}

func CardCompare_DownSuit(c1, c2 *Card) (bool, error) {
	if c1.Suit() != c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return CardCompare_Down(c1, c2)
}

func CardCompare_UpSuitWrap(c1, c2 *Card) (bool, error) {
	if c1.Suit() != c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if c1.Ordinal() == 13 && c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	if c1.Ordinal() == c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func CardCompare_DownSuitWrap(c1, c2 *Card) (bool, error) {
	if c1.Suit() != c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if c1.Ordinal() == 1 && c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if c1.Ordinal()-1 == c2.Ordinal() {
		return true, nil
	}
	return false, errors.New("Cards must go down in rank (Kings on Aces allowed)")
}
