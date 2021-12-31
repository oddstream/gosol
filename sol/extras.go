package sol

import (
	"log"

	"oddstream.games/gomps5/sound"
)

func FindCardOwner(card *Card) Pile {
	for _, pile := range TheBaize.piles {
		for _, c := range pile.Cards() {
			if c == card {
				return pile
			}
		}
	}
	log.Panic("Cannot find card")
	return nil
}

func AnyCardsProne(cards []*Card) bool {
	for _, c := range cards {
		if c.Prone() {
			return true
		}
	}
	return false
}

func FlipUpExposedCard(p Pile) {
	if !p.IsStock() {
		if c := p.Peek(); c != nil {
			c.FlipUp()
		}
	}
}

// MoveCard is an optimized, single card version of MoveCards
func MoveCard(src Pile, dst Pile) *Card {
	if c := src.Pop(); c != nil {
		sound.Play("Place")
		dst.Push(c)
		FlipUpExposedCard(src)
		src.Scrunch()
		src.Refan()
		dst.Scrunch()
		dst.Refan()
		return c
	}
	return nil
}

func MoveNamedCard(suit, ordinal int, dst Pile) {

	// 1. find the card in the library
	var ID CardID = NewCardID(0, suit, ordinal)
	var c *Card
	for i := 0; i < len(CardLibrary); i++ {
		if SameCard(ID, CardLibrary[i].ID) {
			c = &CardLibrary[i]
		}
	}
	if c == nil {
		println("Could not find card", c.String(), "in library")
		return
	}

	// 2.find the card in it's owning pile
	var src Pile = c.Owner()
	var index int = src.IndexOf(c)
	if index == -1 {
		println("Could not find card", c.String(), "in pile")
		return
	}

	// 3. extract the card from it's owning pile
	src.Delete(index)

	// 4. push the card onto the dst pile
	sound.Play("Place")
	c.FlipUp()
	dst.Push(c)
	FlipUpExposedCard(src)
	src.Scrunch()
	src.Refan()
	dst.Scrunch()
	dst.Refan()
}

func MoveCards(src Pile, moveFromIndex int, dst Pile) {

	oldSrcLen := src.Len()

	tmp := make([]*Card, 0, oldSrcLen)

	// pop the tail off the source and push onto temp stack
	for i := oldSrcLen - 1; i >= moveFromIndex; i-- {
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
	src.Refan()
	dst.Scrunch()
	dst.Refan()
}

func GenericTailTapped(self Pile, tail []*Card) {
	if len(tail) != 1 {
		return
	}
	c := tail[0]
	for _, fp := range TheBaize.script.Foundations() {
		if ok, _ := fp.CanAcceptCard(c); ok {
			MoveCard(self, fp)
			break
		}
	}
}

func GenericCollect(self Pile) {
	for _, fp := range TheBaize.script.Foundations() {
		for {
			// loop to get as many cards as possible from this core
			if self.Empty() {
				return
			}
			if ok, _ := fp.CanAcceptCard(self.Peek()); !ok {
				// this foundation doesn't want this card; onto the next one
				break
			}
			MoveCard(self, fp)
		}
	}
}
