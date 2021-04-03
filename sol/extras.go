package sol

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// create the set of cards, into a stock pile
// return the number of cards created
func createCards(stock *Pile) int {

	defer util.Duration(time.Now(), "createCards")
	packs, ok := stock.GetIntAttribute("Packs")
	if !ok || packs == 0 {
		packs = 1
	}

	var createSuitStrings []string
	attribSuits := stock.GetStringAttribute("Suits")
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
				c.owner = stock
				c.SetPosition(stock.BaizePosition())
				stock.Cards = append(stock.Cards, c)
			}
		}
	}

	return stock.CardCount()
}

func shuffleCards(stock *Pile, seed int64) {

	defer util.Duration(time.Now(), "shuffleCards")

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

	if NoShuffle {
		println("not shuffling cards")
		return
	}

	rand.Seed(seed)

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

	rand.Shuffle(len(stock.Cards), func(i, j int) { stock.Cards[i], stock.Cards[j] = stock.Cards[j], stock.Cards[i] })

	// for i, c := range sh.cards {
	// 	println(i, c.ID.String())
	// }

}

func findCard(cards []*Card, card rune) (int, bool) {
	// card should be one of 123456789ABCD
	i64, err := strconv.ParseInt(string(card), 16, 0)
	if err != nil {
		log.Panic("cannot parse", card)
	}
	ordinal := int(i64)
	for i, c := range cards {
		if c.Ordinal() == ordinal {
			return i, true
		}
	}
	return 0, false
}

func CreateScalables() {
	schriftbank.MakeCardFonts(CardWidth) // CardWidth/Height have now been set

	switch TheUserData.CardStyle {
	case "retro":
		TheCIP = NewRetroCardImageProvider()
		CardBackImage = TheCIP.BackImage(TheUserData.CardBackPattern)
	default:
		TheCIP = NewScalableCardImageProvider()
		CardBackImage = TheCIP.BackImage(TheUserData.CardBackColor)
	}
	CardShadowImage = TheCIP.ShadowImage()
	// CardMovableImage = TheCIP.MovableImage()
}
