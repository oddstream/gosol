package sol

import (
	"log"

	"oddstream.games/gomps5/sound"
)

func AnyCardsProne(cards []*Card) bool {
	for _, c := range cards {
		if c.Prone() {
			return true
		}
	}
	return false
}

func FlipUpExposedCard(p *Pile) {
	if _, ok := (p.subtype).(*Stock); ok {
		// this Pile's subtype is *Stock
		// so don't flip an exposed card
	} else {
		if c := p.Peek(); c != nil {
			c.FlipUp()
		}
	}
}

// MoveCard is an optimized, single card version of MoveCards
func MoveCard(src *Pile, dst *Pile) {
	if c := src.Pop(); c != nil {
		sound.Play("Place")
		dst.Push(c)
		FlipUpExposedCard(src)
		src.Scrunch()
		dst.Scrunch()
	}
}

func MoveNamedCard(suit, ordinal int, dst *Pile) {

	// 1. find the card in the library
	var ID CardID = NewCardID(0, suit, ordinal)
	var c *Card
	for i := 0; i < len(TheBaize.cardLibrary); i++ {
		if SameCard(ID, TheBaize.cardLibrary[i].ID) {
			c = &TheBaize.cardLibrary[i]
		}
	}
	if c == nil {
		println("Could not find card", c.String(), "in library")
		return
	}

	// 2.find the card in it's owning pile
	var src *Pile = c.owner
	var index int = -1
	for i := 0; i < len(src.cards); i++ {
		if c == src.cards[i] {
			index = i
			break
		}
	}
	if index == -1 {
		println("Could not find card", c.String(), "in pile")
		return
	}

	// 3. extract the card from it's owning pile
	src.cards = append(src.cards[:index], src.cards[index+1:]...)

	// 4. push the card onto the dst pile
	sound.Play("Place")
	c.FlipUp()
	dst.Push(c)
	FlipUpExposedCard(src)
	src.Scrunch()
	dst.Scrunch()
}

func MoveCards(c *Card, dst *Pile) {

	src := c.Owner()
	if src == nil || !src.Valid() {
		log.Panic("invalid pile")
	}

	oldSrcLen := src.Len()

	// find the index of the first card we will move
	moveFrom := src.IndexOf(c)
	if moveFrom == -1 {
		log.Panic("")
	}

	tmp := make([]*Card, 0, oldSrcLen)

	// pop the tail off the source and push onto temp stack
	for i := oldSrcLen - 1; i >= moveFrom; i-- {
		tmp = append(tmp, src.Pop())
	}

	sound.Play("Slide")

	// pop all cards off the temp stack and onto the destination
	for i := len(tmp) - 1; i >= 0; i-- {
		dst.Push(tmp[i])
	}

	FlipUpExposedCard(src)

	if oldSrcLen == src.Len() {
		log.Println("nothing happened in MoveCards")
	}

	src.Scrunch()
	dst.Scrunch()
}
