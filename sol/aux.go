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
		// this Pile's concrete value is *Stock
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

func MoveCards(c *Card, dst *Pile) {

	src := c.Owner()
	if src == nil || !src.Valid() {
		log.Panic("invalid pile")
	}

	oldSrcLen := src.Len()

	// find the index of the first card we will move
	moveFrom, err := src.IndexOf(c)
	if err != nil {
		log.Panic(err.Error())
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
