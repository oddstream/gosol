package dark

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"time"

	"oddstream.games/gosol/cardid"
	"oddstream.games/gosol/sol"
)

// pileVtabler interface for each subpile type, implements the behaviours
// specific to each subtype
type pileVtabler interface {
	CanAcceptTail([]*Card) (bool, error)
	TailTapped([]*Card)
	Conformant() bool
	UnsortedPairs() int
	MovableTails() []*movableTail
}

// movableTail is used for collecting tap destinations
type movableTail struct {
	dst  *Pile
	tail []*Card
}

// Pile holds the state of the piles and cards therein.
// Pile is exported from this package because it's used to pass between light and dark.
// LIGHT should see a Pile object as immutable, hence the unexported fields and getters.
type Pile struct {
	category string // needed by LIGHT when creating Pile Placeholder (switch)
	label    string // needed by LIGHT when creating Pile Placeholder
	moveType sol.MoveType
	fanType  sol.FanType
	cards    []*Card
	vtable   pileVtabler
	slot     image.Point
}

// Public functions

func (p *Pile) Category() string {
	return p.category
}

func (p *Pile) Label() string {
	return p.label
}

func (p *Pile) Cards() []*Card {
	return p.cards
}

func (p *Pile) Slot() image.Point {
	return p.slot
}

func (p *Pile) FanType() sol.FanType {
	return p.fanType
}

// moveType is not published

// Len returns the number of cards in this pile.
// Len satisfies the sort.Interface interface.
func (self *Pile) Len() int {
	return len(self.cards)
}

// Less satisfies the sort.Interface interface
func (self *Pile) Less(i, j int) bool {
	c1 := self.cards[i]
	c2 := self.cards[j]
	return c1.Suit() < c2.Suit() && c1.Ordinal() < c2.Ordinal()
}

// Swap satisfies the sort.Interface interface
func (self *Pile) Swap(i, j int) {
	self.cards[i], self.cards[j] = self.cards[j], self.cards[i]
}

// Private functions

func newPile(category string, slot image.Point, fanType sol.FanType, moveType sol.MoveType) Pile {
	var self Pile = Pile{
		category: category,
		fanType:  fanType,
		moveType: moveType,
		slot:     slot,
	}
	return self
}

func (self *Pile) reset() {
	self.cards = self.cards[:0]
}

func (self *Pile) isCell() bool {
	_, ok := self.vtable.(*Cell)
	return ok
}

func (self *Pile) isStock() bool {
	_, ok := self.vtable.(*Stock)
	return ok
}

func (self *Pile) shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(self.Len(), self.Swap)
	log.Printf("Shuffled %d cards", self.Len())
}

// delete a *Card from this pile
func (self *Pile) delete(index int) {
	self.cards = append(self.cards[:index], self.cards[index+1:]...)
}

// extract a specific *Card from this pile
func (self *Pile) extract(pack, ordinal, suit int) *Card {
	var ID cardid.CardID = cardid.NewCardID(pack, suit, ordinal)
	for i, c := range self.cards {
		if cardid.SameCardAndPack(ID, c.id) {
			self.delete(i)
			c.flipUp()
			return c
		}
	}
	log.Printf("Could not find card %d %d in %s", suit, ordinal, self.category)
	return nil
}

// peek topmost Card of this Pile (a stack)
func (self *Pile) peek() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	return self.cards[len(self.cards)-1]
}

// pop a Card off the end of this Pile (a stack)
func (self *Pile) pop() *Card {
	if len(self.cards) == 0 {
		return nil
	}
	c := self.cards[len(self.cards)-1]
	self.cards = self.cards[:len(self.cards)-1]
	c.flipUp()
	c.setOwner(nil)
	return c
}

// push a Card onto the end of this Pile (a stack)
func (self *Pile) push(c *Card) {
	self.cards = append(self.cards, c)
	if self.isStock() {
		c.flipDown()
	}
	c.setOwner(self)
}

func (self *Pile) flipUpExposedCard() {
	if !self.isStock() {
		if c := self.peek(); c != nil {
			c.flipUp()
		}
	}
}

func (self *Pile) reverseCards() {
	for i, j := 0, len(self.cards)-1; i < j; i, j = i+1, j-1 {
		self.cards[i], self.cards[j] = self.cards[j], self.cards[i]
	}
}

// buryCards moves cards with the specified ordinal to the beginning of the pile
func (self *Pile) BuryCards(ordinal int) {
	tmp := make([]*Card, 0, cap(self.cards))
	for _, c := range self.cards {
		if c.Ordinal() == ordinal {
			tmp = append(tmp, c)
		}
	}
	for _, c := range self.cards {
		if c.Ordinal() != ordinal {
			tmp = append(tmp, c)
		}
	}
	self.reset()
	for i := 0; i < len(tmp); i++ {
		self.push(tmp[i])
	}
}

// canMoveTail filters out cases where a tail can be moved from a given pile type
// eg if only one card can be moved at a time
func (self *Pile) canMoveTail(tail []*Card) (bool, error) {
	if !self.isStock() {
		if anyCardsProne(tail) {
			return false, errors.New("Cannot move a face down card")
		}
	}
	switch self.moveType {
	case sol.MOVE_NONE:
		// eg Discard, Foundation
		return false, fmt.Errorf("Cannot move a card from a %s", self.category)
	case sol.MOVE_ANY:
		// well, that was easy
	case sol.MOVE_ONE:
		// eg Cell, Reserve, Stock, Waste
		if len(tail) > 1 {
			return false, fmt.Errorf("Can only move one card from a %s", self.category)
		}
	case sol.MOVE_ONE_PLUS:
		// don't (yet) know destination, so we allow this as MOVE_ANY
		// and do power moves check later, in Tableau CanAcceptTail
	case sol.MOVE_ONE_OR_ALL:
		// Canfield, Toad
		if len(tail) == 1 {
			// that's okay
		} else if len(tail) == self.Len() {
			// that's okay too
		} else {
			return false, errors.New("Can only move one card, or the whole pile")
		}
	}
	return true, nil
}

func (self *Pile) makeTail(c *Card) []*Card {
	if c.owner() != self {
		log.Panic("Pile.MakeTail called with a card that is not of this pile")
	}
	if c == self.peek() {
		return []*Card{c}
	}
	for i, pc := range self.cards {
		if pc == c {
			return self.cards[i:]
		}
	}
	log.Panic("Pile.MakeTail made an empty tail")
	return nil
}

func (self *Pile) defaultTailTapped(tail []*Card) {
	card := tail[0]
	if card.tapDestination != nil {
		csrc := card.owner()
		ctail := csrc.makeTail(card)
		if len(ctail) == 1 {
			MoveCard(csrc, card.tapDestination)
		} else {
			MoveTail(card, card.tapDestination)
		}
	}
}
