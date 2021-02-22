package sol

import (
	"log"
	"math/rand"
	"sort"
	"time"

	"oddstream.games/gosol/util"
)

func createCards(stock *Pile) {

	packs, ok := stock.GetIntAttribute("Packs")
	if !ok || packs == 0 {
		packs = 1
	}
	// gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range [4]string{"Club", "Diamond", "Heart", "Spade"} {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				c.owner = stock
				x, y := stock.Position()
				c.SetPosition(x, y)
				stock.Cards = append(stock.Cards, c)
			}
		}
	}
}

func shuffleCards(stock *Pile) {

	rand.Seed(time.Now().UnixNano())

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
	sort.Slice(stock.Cards, func(i, j int) bool { return stock.Cards[i].id < stock.Cards[j].id })

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

func isConformant0(rules int, cPrev, cThis *Card) bool {
	if cPrev.prone || cThis.prone {
		println("prone cards are not conformant")
		return false
	}

	localSuit := rules / 10

	switch localSuit {
	case 0: // may not build or move
		return false
	case 1: // regardless of suit
	case 2: // in suit
		if cPrev.suit != cThis.suit {
			return false
		}
	case 3: // in color
		if cPrev.color != cThis.color {
			return false
		}
	case 4: // in alternate color
		if cPrev.color == cThis.color {
			return false
		}
	case 5: // in any suit but it's own
		if cPrev.suit == cThis.suit {
			return false
		}
	}

	localRank := rules % 10
	rankWrap := false // TODO

	if rankWrap {
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if cPrev.ordinal == 13 && cThis.ordinal == 1 {
				// an Ace on a King
			} else {
				if cThis.ordinal != cPrev.ordinal+1 {
					return false
				}
			}
		case 2: // down, e.g. a 9 goes on a 10
			if cPrev.ordinal == 1 && cThis.ordinal == 13 {
				// a King on an Ace
			} else {
				if cThis.ordinal != cPrev.ordinal-1 {
					return false
				}
			}
		case 4: // either up or down
			if (cPrev.ordinal == 13 && cThis.ordinal == 1) || (cPrev.ordinal == 1 && cThis.ordinal == 13) {
				// a king on an ace or an ace on a king
			} else {
				if util.Abs(cPrev.ordinal-cThis.ordinal) != 1 {
					return false
				}
			}
		case 5: // regardless of rank
		}
	} else { // rankWrap == false
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if cThis.ordinal != cPrev.ordinal+1 {
				return false
			}
		case 2: // down, e.g. a 9 goes on a 10
			if cThis.ordinal != cPrev.ordinal-1 {
				return false
			}
		case 4: // either up or down
			if util.Abs(cThis.ordinal-cPrev.ordinal) != 1 {
				return false
			}
		case 5: // regardless of rank
		}
	}

	// TODO localRank == 13 (Pyramid) cPrev.ordinal + cThis.ordinal == 13

	return true
}

func isConformant(rules int, cards []*Card) bool {
	if nil == cards || len(cards) == 0 {
		log.Fatal("isConformant passed empty tail")
	}
	cPrev := cards[0]
	for n := 1; n < len(cards); n++ {
		cThis := cards[n]
		if !isConformant0(rules, cPrev, cThis) {
			return false
		}
		cPrev = cThis
	}
	return true
}
