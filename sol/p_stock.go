package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
	"log"
	"math/rand"
	"time"
)

func CreateCardLibrary(packs int, suits int, cardFilter *[14]bool, jokersPerPack int) {

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
	cardsRequired += packs * jokersPerPack
	CardLibrary = make([]Card, 0, cardsRequired)

	for pack := 0; pack < packs; pack++ {
		for ord := 1; ord < 14; ord++ {
			for suit := 0; suit < suits; suit++ {
				if cardFilter[ord] {
					/*
						suits are numbered NOSUIT=0, CLUB=1, DIAMOND=2, HEART=3, SPADE=4
						(i.e. not 0..3)
						run the suits loop backwards, so spades are used first
						(folks expect Spider One Suit to use spades)
					*/
					var c Card = NewCard(pack, SPADE-suit, ord)
					CardLibrary = append(CardLibrary, c)
				}
			}
		}
		for i := 0; i < jokersPerPack; i++ {
			var c Card = NewCard(pack, NOSUIT, 0) // NOSUIT and ordinal == 0 creates a joker
			CardLibrary = append(CardLibrary, c)
		}
	}
	log.Printf("%d packs, %d suits, %d cards created\n", packs, suits, len(CardLibrary))
}

type Stock struct {
	Core
}

func (self *Stock) FillFromLibrary() {
	if !self.Empty() {
		log.Panic("stock should be empty")
	}
	for i := 0; i < len(CardLibrary); i++ {
		var c *Card = &CardLibrary[i]
		if !c.Valid() {
			log.Panicf("invalid card at library index %d", i)
		}
		// the following mimics Base.Push
		self.Append(c)
		c.SetOwner(self)
		// don't set Card.pos here
		// so that a new deal makes the spinning cards fall into place
		// without going back to the CardStartPoint
		c.SetProne(true)
	}
}

func (self *Stock) Shuffle() {

	if !self.Valid() {
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
		// it doesn't make sense, but testing shows that you need to do this
		// more than once to get a randomly distributed shuffle
		rand.Shuffle(self.Len(), self.Swap)
	}
}

/*
	Push a card onto the Stock pile.
	Specific version of Push that overrides Core.Push,
	to try to get rid of IsStock() and FindOwner() in Core.Push,
	and also since Stock piles are never fanned.
	Trouble is, this func is only called by Stock Recycle
	(UpdateFromSaveable uses Core, not Pile).
*/
func (self *Stock) Push(card *Card) {
	self.Append(card)
	card.TransitionTo(card.BaizePos())
	card.FlipDown()
	card.SetOwner(self)
	TheBaize.setFlag(dirtyCardPositions)
}

func NewStock(slot image.Point, fanType FanType, packs int, suits int, cardFilter *[14]bool, jokersPerPack int) *Stock {
	CreateCardLibrary(packs, suits, cardFilter, jokersPerPack)
	stock := &Stock{Core: NewCore("Stock", slot, fanType, MOVE_ONE)}
	stock.FillFromLibrary()
	stock.Shuffle()
	TheBaize.AddPile(stock)
	return stock
}

func (*Stock) CanAcceptCard(*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) CanAcceptTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) TailTapped([]*Card) {
	// do nothing, handled by script, which had first dibs
}

func (*Stock) Collect() {
	// never collect from the stock
	// over-ride Core collect to do nothing
}

func (self *Stock) Conformant() bool {
	return self.Empty()
}

func (self *Stock) Complete() bool {
	return self.Empty()
}

func (self *Stock) UnsortedPairs() int {
	// Stock is always considered unsorted
	if self.Empty() {
		return 0
	}
	return self.Len() - 1
}
