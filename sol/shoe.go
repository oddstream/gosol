package sol

import "math/rand"

// Shoe is a gaming device, mainly used in casinos, to hold multiple decks of playing cards
type Shoe struct {
	cards []*Card
}

// NewShoe creates a new Shoe with the given number of packs of cards
func NewShoe(packs int) *Shoe {
	sh := &Shoe{}
	// gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range [4]string{"Club", "Diamond", "Heart", "Spade"} {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				sh.cards = append(sh.cards, c)
			}
		}
	}
	println("created", len(sh.cards), "cards")
	return sh
}

// Shuffle the cards in the Shoe
func (sh *Shoe) Shuffle() {

	// bubble sort cards in order before sorting
	// TODO use go sort package
	swapped := true
	for swapped {
		swapped = false
		for i := 1; i < len(sh.cards); i++ {
			if sh.cards[i-1].id > sh.cards[i].id {
				sh.cards[i], sh.cards[i-1] = sh.cards[i-1], sh.cards[i]
				swapped = true
			}
		}
	}

	println("-bubble-----------------")
	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }

	println("-KFY----------------")
	// Knuth Fisher-Yates shuffle
	for i := len(sh.cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		sh.cards[i], sh.cards[j] = sh.cards[j], sh.cards[i]
	}

	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }
}
