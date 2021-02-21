package sol

import (
	"math/rand"
	"sort"
	"time"
)

func createCards(stock *Pile) {

	packs := stock.GetIntAttribute("Packs")
	if packs == 0 {
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
