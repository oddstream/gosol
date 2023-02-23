package dark

import (
	"log"
	"sort"
)

type scriptBase struct {
	cells       []*Pile
	discards    []*Pile
	foundations []*Pile
	reserves    []*Pile
	stock       *Pile
	tableaux    []*Pile
	waste       *Pile
	wikipedia   string
	cardColors  int
}

// complete - default is number of cards in Foundations == total number of cards.
//
// In Bisley, there may be <13 cards in a Foundation.
// This will need overriding for any variants with Discard piles.
// Could also do this by checking if any pile other than a Foundation is not empty.
func (sb scriptBase) complete() bool {
	var n = 0
	for _, f := range sb.foundations {
		n += len(f.cards)
	}
	return n == len(baize.pack)
}

func (sb scriptBase) spiderComplete() bool {
	for _, tab := range sb.tableaux {
		switch len(tab.cards) {
		case 0:
			// that's fine
		case 13:
			if !tab.vtable.conformant() {
				return false
			}
		default:
			return false
		}
	}
	return true
}

type scripter interface {
	buildPiles()
	startGame()
	afterMove()

	tailMoveError([]*Card) (bool, error)
	aailAppendError(*Pile, []*Card) (bool, error)
	unsortedPairs(*Pile) int

	tailTapped([]*Card)
	pileTapped(*Pile)

	cells() []*Pile
	discards() []*Pile
	foundations() []*Pile
	reserves() []*Pile
	stock() *Pile
	tableaux() []*Pile
	waste() *Pile

	complete() bool
	wikipedia() string
	cardColors() int
	safeCollect() bool
}

var variants = map[string]scripter{}

var variantGroups = map[string][]string{}

// init is used to assemble the "> All" alpha-sorted group of variants
func init() {
	var vnames []string = make([]string, 0, len(variants))
	for k := range variants {
		vnames = append(vnames, k)
	}
	// no need to sort here, sort gets done by func VariantNames()
	variantGroups["> All"] = vnames
	variantGroups["> All by Played"] = vnames
}

func (d *dark) ListVariantGroups() []string {
	var vnames []string = make([]string, 0, len(variantGroups))
	for k := range variantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

func (d *dark) ListVariants(group string) []string {
	var vnames []string = nil
	vnames = append(vnames, variantGroups[group]...)
	if group == "> All by Played" {
		sort.Slice(vnames, func(i, j int) bool { return TheStatistics.Played(vnames[i]) > TheStatistics.Played(vnames[j]) })
	} else {
		sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	}
	return vnames
}

// useful generic game library of functions ///////////////////////////////////

func anyCardsProne(cards []*Card) bool {
	for _, c := range cards {
		if c.Prone() {
			return true
		}
	}
	return false
}

// moveCard moves the top card from src to dst
func moveCard(src *Pile, dst *Pile) *Card {
	if c := src.pop(); c != nil {
		dst.push(c)
		src.flipUpExposedCard()
		return c
	}
	return nil
}

// moveTail moves all the cards from card downwards onto dst
func moveTail(card *Card, dst *Pile) {
	var src *Pile = card.owner()
	tmp := make([]*Card, 0, len(src.cards))
	// pop cards from src upto and including the head of the tail
	for {
		var c *Card = src.pop()
		if c == nil {
			log.Panicf("MoveTail could not find %s", card)
		}
		tmp = append(tmp, c)
		if c == card {
			break
		}
	}
	// pop cards from the tmp stack and push onto dst
	if len(tmp) > 0 {
		for len(tmp) > 0 {
			var c *Card = tmp[len(tmp)-1]
			tmp = tmp[:len(tmp)-1]
			dst.push(c)
		}
		src.flipUpExposedCard()
	}
}
