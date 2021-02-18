package sol

import (
	"math/rand"
	"sort"
	"time"
)

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

// CreateCards fills the pile with packs*52 new cards
func (s *Stock) CreateCards(packs int) {
	// gotcha don't use make([]*Card, packs*52) as it makes a lot of nil entries
	for pack := 0; pack < packs; pack++ {
		for _, suit := range [4]string{"Club", "Diamond", "Heart", "Spade"} {
			for ord := 1; ord < 14; ord++ {
				c := NewCard(pack, suit, ord)
				c.owner = s // Stock implements the CardOwner interface
				x, y := s.Position()
				c.PositionTo(x, y)
				s.cards = append(s.cards, c)
			}
		}
	}
	// println("created", len(p.cards), "cards")
}

// Shuffle the cards in the Stock
func (s *Stock) Shuffle() {

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
	sort.Slice(s.cards, func(i, j int) bool { return s.cards[i].id < s.cards[j].id })

	// println("-ordered------------")
	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }

	// println("-KFY----------------")
	// Knuth Fisher-Yates shuffle
	for i := len(s.cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	}

	// for i, c := range sh.cards {
	// 	println(i, c.id)
	// }
}
