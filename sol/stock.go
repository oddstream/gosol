package sol

import "math/rand"

// Stock is the home for all cards
type Stock struct {
	Pile
}

// NewStock creates a new Stock with the given number of packs of cards
func NewStock(packs, x, y int) *Stock {
	s := &Stock{Pile{X: x, Y: y}}
	s.CreateCards(packs)
	println("created", len(s.cards), "cards")
	return s
}

// Shuffle the cards in the Stock
func (s *Stock) Shuffle() {

	// bubble sort cards in order before sorting
	// TODO use go sort package
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

	println("-bubble-----------------")
	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }

	println("-KFY----------------")
	// Knuth Fisher-Yates shuffle
	for i := len(s.cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	}

	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }
}
