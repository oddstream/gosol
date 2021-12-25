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
	CardLibrary = make([]Card, 0, cardsRequired)

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
					CardLibrary = append(CardLibrary, c)
				}
			}
		}
	}
	log.Printf("%d packs, %d suits, %d cards created\n", packs, suits, len(CardLibrary))
}

type Stock struct {
	pile     *Pile
	recycles int
}

func (p *Pile) FillFromLibrary() {
	if !p.Empty() {
		log.Panic("stock should be empty")
	}
	for i := 0; i < len(CardLibrary); i++ {
		var c *Card = &CardLibrary[i]
		if c == nil || !c.Valid() {
			log.Panicf("invalid card at library index %d", i)
		}
		// the following mimics Base.Push
		p.Append(c)
		c.SetOwner(p)
		c.pos = p.BaizePos() // start at the Stock pile position
		c.src = image.Point{0, 0}
		c.dst = c.pos
		c.SetProne(true)
		// s.Push(c)
	}
}

func (p *Pile) Shuffle() {

	if p == nil || !p.Valid() {
		log.Fatal("invalid stock")
	}
	if NoShuffle {
		log.Println("not shuffling cards")
		return
	}
	seed := time.Now().UnixNano() & 0xFFFFFFFF
	if DebugMode {
		log.Println("shuffle with seed", seed)
	}
	rand.Seed(seed)
	for i := 0; i < 6; i++ {
		rand.Shuffle(p.Len(), p.Swap)
	}
}

func NewStock(slot image.Point, fanType FanType, packs int, suits int, cardFilter *[14]bool) *Pile {

	p := &Pile{}
	p.Ctor(&Stock{pile: p, recycles: 32767}, "Stock", slot, fanType, MOVE_ONE)
	CreateCardLibrary(packs, suits, cardFilter)
	p.FillFromLibrary()
	p.Shuffle()

	return p
}

func (s *Stock) CanAcceptCard(*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (s *Stock) CanAcceptTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) TailTapped([]*Card) {
	// do nothing, handled by script, which had first dibs
}

func (s *Stock) Collect() {
	// never collect from the stock
	// over-ride base collect to do nothing
}

func (s *Stock) Conformant() bool {
	return s.pile.Empty()
}

func (s *Stock) Complete() bool {
	return s.pile.Empty()
}

func (s *Stock) UnsortedPairs() int {
	// Stock is always considered unsorted
	if s.pile.Empty() {
		return 0
	}
	return s.pile.Len() - 1
}
