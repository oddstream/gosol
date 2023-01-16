package sol

import (
	"log"

	"oddstream.games/gosol/sound"
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
	if !p.IsStock() {
		if c := p.Peek(); c != nil {
			c.FlipUp()
		}
	}
}

// MoveCard moves the top card from src to dst
func MoveCard(src *Pile, dst *Pile) *Card {
	if c := src.Pop(); c != nil {
		dst.Push(c)
		FlipUpExposedCard(src)
		sound.Play("Place")
		return c
	}
	return nil
}

// MoveTail moves all the cards from card downwards onto dst
func MoveTail(card *Card, dst *Pile) {
	var src *Pile = card.Owner()
	tmp := make([]*Card, 0, len(src.cards))
	// pop cards from src upto and including the head of the tail
	for {
		var c *Card = src.Pop()
		if c == nil {
			log.Panicf("MoveTail could not find %s", card)
		}
		tmp = append(tmp, c)
		if c == card {
			break
		}
	}
	// pop cards from the tmp stack and push onto dst
	if len(tmp) > 0 {
		for len(tmp) > 0 {
			var c *Card = tmp[len(tmp)-1]
			tmp = tmp[:len(tmp)-1]
			dst.Push(c)
		}
		FlipUpExposedCard(src)
		sound.Play("Place")
	}
}

/*
MoveTail2 - too low-level to be of any use

var src []int = []int{1, 2, 3, 4, 5}
var dst []int = []int{7, 8, 9}

	func main() {
		fmt.Println(src)
		fmt.Println(dst)

		var idx int = 4

		if idx == 0 {
			dst = append(dst, src...)
			src = []int{}
		} else {
			var tail []int = src[len(src)-(idx-1):]
			fmt.Println(tail)
			src = src[:idx]
			dst = append(dst, tail...)
		}
		fmt.Println(src)
		fmt.Println(dst)
	}
*/
// func MoveTail2(card *Card, dst *Pile) {
// 	var src *Pile = card.owner
// 	if src.Peek() == card {
// 		MoveCard(src, dst)
// 		return
// 	}
// 	var idx int
// 	if idx = src.IndexOf(card); idx == -1 {
// 		log.Panicf("Card %s not found", card.String())
// 	}
// 	if idx == 0 {
// 		dst.cards = append(dst.cards, src.cards...)
// 		src.cards = src.cards[:0] //[]*Card{}
// 	} else {
// 		var tail []*Card = src.cards[len(src.cards)-(idx-1):]
// 		src.cards = src.cards[:idx]
// 		dst.cards = append(dst.cards, tail...)
// 	}
// 	for _, c := range dst.cards {
// 		c.owner = dst
// 	}
// 	TheBaize.setFlag(dirtyCardPositions)
// }

// func MoveAllCards(src *Pile, dst *Pile) {
// 	if src.Empty() {
// 		return
// 	}
// 	for i := 0; i < src.Len(); i++ {
// 		dst.Push(src.Get(i))
// 	}
// 	src.Reset()
// 	TheBaize.setFlag(dirtyCardPositions)
// }
