package sol

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"oddstream.games/gosol/util"
)

// CreateStock finds the stock pile and creates the stock cards
// sets Baize.totalCards
func (b *Baize) CreateStock() {

	// defer util.Duration(time.Now(), "CreateStock")

	packs, ok := b.stock.GetIntAttribute("Packs")
	if !ok || packs == 0 {
		packs = 1
	}

	var createSuitStrings []string
	attribSuits := b.stock.GetStringAttribute("Suits")
	if attribSuits == "" {
		createSuitStrings = []string{"Club", "Diamond", "Heart", "Spade"}
	} else {
		createSuitStrings = strings.Split(attribSuits, ",")
	}
	var createSuitInts []int
	for _, suit := range createSuitStrings {
		createSuitInts = append(createSuitInts, SuitStringToInt(suit))
	}

	// golang gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range createSuitInts {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				c.owner = b.stock
				c.SetPosition(b.stock.BaizePosition())
				b.stock.Cards = append(b.stock.Cards, c)
			}
		}
	}

	// if DebugMode {
	// 	TestShuffle(b.stock)
	// }

	b.totalCards = b.stock.CardCount()
}

func (b *Baize) Shuffle() {

	// defer util.Duration(time.Now(), "ShuffleStock")

	cards := b.stock.Cards

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

	/*
		used to restart the same game by reusing the random seed
		but no longer do that (now we unwind the undo stack)
		so we no longer need to sort cards into order before shuffle

		sort.Slice(cards, func(i, j int) bool { return cards[i].ID < cards[j].ID })
	*/

	if NoShuffle {
		println("not shuffling cards")
		return
	}

	// println("-ordered------------")
	// for i, c := range sh.cards {
	// 	println(i, c.ID.String())
	// }

	// println("-KFY----------------")
	// Knuth Fisher-Yates shuffle
	// for i := len(stock.Cards) - 1; i > 0; i-- {
	// 	j := rand.Intn(i + 1)
	// 	stock.Cards[i], stock.Cards[j] = stock.Cards[j], stock.Cards[i]
	// }

	// tmp := make([]*Card, len(cards), cap(cards))
	// copy(tmp, cards)

	// TestShuffle shows that the 7 can have a consistently lower distribution; shuffling twice corrects this
	seed := time.Now().UnixNano()
	if DebugMode {
		log.Println("seed", seed)
	}
	rand.Seed(seed)
	for range []int{1, 2, 3, 4, 5, 6} {
		rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	}

	// var notShuffled int
	// for i := 0; i < len(tmp); i++ {
	// 	if tmp[i] == cards[i] {
	// 		println("not shuffled at", i)
	// 		notShuffled++
	// 	}
	// }
	// if notShuffled > 0 {
	// 	println(notShuffled, "cards not shuffled")
	// } else {
	// 	println("all cards shuffled")
	// }
}

func TestShuffle(stock *Pile) {
	const cycles int = 500000
	dist := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < cycles; i++ {
		// sort.Slice(stock.Cards, func(i, j int) bool { return stock.Cards[i].ID < stock.Cards[j].ID })
		rand.Shuffle(len(stock.Cards), func(i, j int) { stock.Cards[i], stock.Cards[j] = stock.Cards[j], stock.Cards[i] })
		for j, c := range stock.Cards {
			// 1 % 13 = 1
			// 13 % 13 = 0
			if c.Ordinal()%13 == j {
				dist[j]++
			}
		}
	}
	for i := 0; i < len(dist); i++ {
		println(i, dist[i])
	}
}

// func findHexCard(cards []*Card, card rune) (int, bool) {
// 	// card should be one of 123456789ABCD
// 	i64, err := strconv.ParseInt(string(card), 16, 0)
// 	if err != nil {
// 		log.Panic("cannot parse", card)
// 	}
// 	ordinal := int(i64)
// 	for i, c := range cards {
// 		if c.Ordinal() == ordinal {
// 			return i, true
// 		}
// 	}
// 	return 0, false
// }

func parseCardsFromDeal(stock *Pile, deal string) []*Card {
	var cards []*Card
	var runes = []rune(deal)
	for i := 0; i < len(runes); i++ {
		// Note that since the rune type is an alias for int32, we must use %c instead of the usual %v in the Printf statement,
		// or we will see the integer representation of the Unicode code point
		var c *Card
		if runes[i] == 'u' || runes[i] == 'd' {
			c = stock.Pop()
		} else {
			ord := util.RuneToOrdinal(runes[i])
			i++
			suit := util.RuneToSuit(runes[i])
			i++
			cid := NewCardID(0, suit, ord)
			for idx, cs := range stock.Cards {
				if SameCard(cid, cs.ID) { // ignores pack
					c = stock.Extract(idx)
					break
				}
			}
		}
		if c == nil {
			log.Fatal("out of cards during deal from ", deal)
			break
		}
		if runes[i] == 'd' {
			c.FlipDown()
		}
		cards = append(cards, c)
	}
	return cards
}
