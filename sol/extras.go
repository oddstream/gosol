package sol

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"oddstream.games/gosol/util"
)

func createCards(stock *Pile) {

	packs, ok := stock.GetIntAttribute("Packs")
	if !ok || packs == 0 {
		packs = 1
	}

	var createSuitStrings []string
	attribSuits := stock.GetStringAttribute("Suits")
	if attribSuits != "" {
		createSuitStrings = strings.Split(attribSuits, ",")
	} else {
		createSuitStrings = []string{"Club", "Diamond", "Heart", "Spade"}
	}
	var createSuitInts []int
	for _, suit := range createSuitStrings {
		switch suit {
		case "Club":
			createSuitInts = append(createSuitInts, 1)
		case "Diamond":
			createSuitInts = append(createSuitInts, 2)
		case "Heart":
			createSuitInts = append(createSuitInts, 3)
		case "Spade":
			createSuitInts = append(createSuitInts, 4)
		}
	}

	// gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range createSuitInts {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				c.owner = stock
				c.SetPosition(stock.Position())
				stock.Cards = append(stock.Cards, c)
			}
		}
	}
}

func shuffleCards(stock *Pile, seed int64) {

	rand.Seed(seed)

	// sort cards in order before shuffle (why?)
	/*
		swapped := true
		for swapped {
			swapped = false
			for i := 1; i < len(s.cards); i++ {
				if s.cards[i-1].id > s.cards[i].id {
					s.cards[i], s.cards[i-1] = s.cards[i-1], s.cards[i]
					swapped = true
				}
			}
		}
	*/
	sort.Slice(stock.Cards, func(i, j int) bool { return stock.Cards[i].ID < stock.Cards[j].ID })

	// println("-ordered------------")
	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }

	// println("-KFY----------------")
	// Knuth Fisher-Yates shuffle
	for i := len(stock.Cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		stock.Cards[i], stock.Cards[j] = stock.Cards[j], stock.Cards[i]
	}

	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }

}

func findCard(cards []*Card, card rune) (int, bool) {
	// card should be one of 123456789ABCD
	i64, err := strconv.ParseInt(string(card), 16, 0)
	if err != nil {
		log.Fatal("cannot parse", card)
	}
	ordinal := int(i64)
	for i, c := range cards {
		if c.Ordinal() == ordinal {
			return i, true
		}
	}
	return 0, false
}

func isConformant0(rules, flags int, cPrev, cThis *Card) bool {
	if cPrev.Prone() || cThis.Prone() {
		println("prone cards are not conformant")
		return false
	}

	localSuit := rules / 10
	localRank := rules % 10

	switch localSuit {
	case 0: // may not build or move
		return false
	case 1: // regardless of suit
	case 2: // in suit
		if cPrev.Suit() != cThis.Suit() {
			return false
		}
	case 3: // in color
		if cPrev.Color() != cThis.Color() {
			return false
		}
	case 4: // in alternate color
		if cPrev.Color() == cThis.Color() {
			return false
		}
	case 5: // in any suit but it's own
		if cPrev.Suit() == cThis.Suit() {
			return false
		}
	}

	if flags&1 == 1 { // rank wrap == true
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if cPrev.Ordinal() == 13 && cThis.Ordinal() == 1 {
				// an Ace on a King
			} else {
				if cThis.Ordinal() != cPrev.Ordinal()+1 {
					return false
				}
			}
		case 2: // down, e.g. a 9 goes on a 10
			if cPrev.Ordinal() == 1 && cThis.Ordinal() == 13 {
				// a King on an Ace
			} else {
				if cThis.Ordinal() != cPrev.Ordinal()-1 {
					return false
				}
			}
		case 4: // either up or down
			if (cPrev.Ordinal() == 13 && cThis.Ordinal() == 1) || (cPrev.Ordinal() == 1 && cThis.Ordinal() == 13) {
				// a king on an ace or an ace on a king
			} else {
				if util.Abs(cPrev.Ordinal()-cThis.Ordinal()) != 1 {
					return false
				}
			}
		case 5: // regardless of rank
		}
	} else { // rank wrap == false
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if cThis.Ordinal() != cPrev.Ordinal()+1 {
				return false
			}
		case 2: // down, e.g. a 9 goes on a 10
			if cThis.Ordinal() != cPrev.Ordinal()-1 {
				return false
			}
		case 4: // either up or down
			if util.Abs(cThis.Ordinal()-cPrev.Ordinal()) != 1 {
				return false
			}
		case 5: // regardless of rank
		}
	}

	// TODO localRank == 13 (Pyramid) cPrev.Ordinal() + cThis.Ordinal() == 13

	return true
}

func isConformant(rules, flags int, cards []*Card) bool {
	if nil == cards || len(cards) == 0 {
		log.Fatal("isConformant passed empty tail")
	}
	if rules == 0 {
		return false // may not build or move, even a single card
	}
	cPrev := cards[0]
	for n := 1; n < len(cards); n++ {
		cThis := cards[n]
		if !isConformant0(rules, flags, cPrev, cThis) {
			return false
		}
		cPrev = cThis
	}
	return true
}

func powerMoves(piles []*Pile, pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range piles {
		switch p.Class {
		case "Cell":
			if 0 == p.CardCount() {
				emptyCells++
			}
		case "Tableau":
			// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
			if p == pDraggingTo && 0 == pDraggingTo.CardCount() {
				// empty column doesn't count
			} else if 0 == p.CardCount() {
				emptyCols++
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
}
