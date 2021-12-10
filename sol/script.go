package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"log"

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

func GetVariantInterface(v string) ScriptInterface {
	switch v {
	case "Clondike":
		return &Clondike{}
	case "Easy":
		return &Easy{}
	case "Freecell":
		return &Freecell{}
	case "Simple Simon":
		return &SimpleSimon{}
	case "Spider One Suit":
		return &Spider{packs: 8, suits: 1}
	case "Spider Two Suits":
		return &Spider{packs: 4, suits: 2}
	default:
		log.Panicf("Unknown variant %s", v)
	}
	return nil
}

func VariantNames() []string {
	//	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return []string{"Clondike", "Easy", "Freecell", "Simple Simon", "Spider One Suit", "Spider Two Suits"}
}

// useful generic game library of functions

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
