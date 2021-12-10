package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
	"log"
	"math/rand"
	"time"
)

func CreateCardLibrary(packs int, suits int, cardFilter *[14]bool) {

	var numberOfCardsInSuit int = 0
	if cardFilter == nil {
		cardFilter = &[14]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true}
		numberOfCardsInSuit = 13
	} else {
		for i := 1; i < 14; i++ {
			if cardFilter[i] {
				numberOfCardsInSuit++
			}
		}
	}
	var cardsRequired int = packs * suits * numberOfCardsInSuit
	TheBaize.cardLibrary = make([]Card, 0, cardsRequired)

	for pack := 0; pack < packs; pack++ {
		for ord := 1; ord < 14; ord++ {
			for suit := 0; suit < suits; suit++ {
				if cardFilter[ord] {
					/*
						suits are numbered CLUB=1, DIAMOND=2, HEART=3, SPADE=4
						(i.e. not 0..3)
						run the suits loop backwards, so spades are used first
						(folks expect Spider One Suit to use spades)
					*/
					var c Card = NewCard(pack, SPADE-suit, ord)
					TheBaize.cardLibrary = append(TheBaize.cardLibrary, c)
				}
			}
		}
	}
	log.Printf("%d packs, %d suits, %d cards created\n", packs, suits, len(TheBaize.cardLibrary))
}

type Stock struct {
	Base
	recycles int
}

func (s *Stock) FillFromLibrary() {
	if !s.Empty() {
		log.Panic("stock should be empty")
	}
	for i := 0; i < len(TheBaize.cardLibrary); i++ {
		var c *Card = &TheBaize.cardLibrary[i]
		if c == nil || !c.Valid() {
			log.Panicf("invalid card at library index %d", i)
		}
		// the following mimics Base.Push
		s.Append(c)
		c.TransitionTo(s.BaizePos())
		c.SetOwner(s)
		c.SetProne(true)
		// s.Push(c)
	}
}

func (s *Stock) Shuffle() {

	if s == nil || !s.Valid() {
		log.Fatal("invalid stock")
	}
	if NoShuffle {
		log.Println("not shuffling cards")
		return
	}
	/*
		used to restart the same game by reusing the random seed
		but no longer do that (now we unwind the undo stack)
		so we no longer need to sort cards into order before shuffle

		sort.Slice(p.Cards, func(i, j int) bool { return p.Cards[i].ID < p.Cards[j].ID })
	*/
	seed := time.Now().UnixNano()
	if DebugMode {
		log.Println("shuffle with seed", seed)
	}
	rand.Seed(seed)
	// for range []int{1, 2, 3, 4, 5, 6} {
	// rand.Shuffle(len(base.cards), func(i, j int) { base.cards[i], base.cards[j] = base.cards[j], base.cards[i] })
	rand.Shuffle(s.Len(), s.Swap)
	// }
}

func NewStock(slot image.Point, fanType FanType, packs int, suits int, cardFilter *[14]bool) *Stock {

	s := &Stock{}
	s.Ctor(s, "Stock", slot, fanType)
	s.recycles = 32767

	CreateCardLibrary(packs, suits, cardFilter)
	s.FillFromLibrary()
	s.Shuffle()

	return s
}

func (s *Stock) CanMoveTail(tail []*Card) (bool, error) {
	if len(tail) != 1 {
		return false, errors.New("Can only move a single Stock card")
	}
	return true, nil
}

func (s *Stock) CanAcceptCard(*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (s *Stock) CanAcceptTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (s *Stock) TailTapped(tail []*Card) {
	TheBaize.script.TailTapped(tail)
}

func (s *Stock) Collect() {
	// never collect from the stock
	// over-ride base collect to do nothing
}

func (s *Stock) Conformant() bool {
	return s.Empty()
}

func (s *Stock) Complete() bool {
	return s.Empty()
}

func (s *Stock) UnsortedPairs() int {
	if s.Empty() {
		return 0
	}
	return s.Len() - 1
}

// Reset overside Base Reset to make fill the stock with shuffled cards
func (s *Stock) Reset() {
	s.cards = s.cards[:0]
	s.FillFromLibrary()
	s.Shuffle()
}
