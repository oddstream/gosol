package sol

import (
	"fmt"
)

func (b *Baize) FindHomeForTail(owner Pile, tail []*Card) Pile {
	if len(tail) == 1 {
		for _, dst := range b.script.Foundations() {
			if dst == owner {
				continue
			}
			if ok, _ := dst.CanAcceptTail(tail); ok {
				return dst
			}
		}
		for _, dst := range b.script.Cells() {
			if dst == owner {
				continue
			}
			if ok, _ := dst.CanAcceptTail(tail); ok {
				return dst
			}
		}
	}
	for _, dst := range b.script.Tableaux() {
		if dst == owner {
			continue
		}
		if ok, _ := dst.CanAcceptTail(tail); ok {
			return dst
		}
	}
	return nil
}

func (b *Baize) Stuck() bool {
	var moves int

	if !b.script.Stock().Empty() {
		moves++
	}
	if wastePile := b.script.Waste(); wastePile != nil {
		if !wastePile.Empty() {
			var tail []*Card
			tail = append(tail, wastePile.Peek())
			if dst := b.FindHomeForTail(wastePile, tail); dst != nil {
				moves++
			}

			if b.script.Stock().Empty() && b.recycles > 0 {
				println("can recycle")
				moves++
			}
		}
	}

	for _, pile := range b.script.Cells() {
		for _, card := range pile.cards {
			if card.Prone() {
				continue
			}
			tail := pile.MakeTail(card)
			if ok, _ := pile.CanMoveTail(tail); !ok {
				continue
			}
			if dst := b.FindHomeForTail(pile, tail); dst != nil {
				moves++
			}
		}
	}
	for _, pile := range b.script.Reserves() {
		for _, card := range pile.cards {
			if card.Prone() {
				continue
			}
			tail := pile.MakeTail(card)
			if ok, _ := pile.CanMoveTail(tail); !ok {
				continue
			}
			if dst := b.FindHomeForTail(pile, tail); dst != nil {
				moves++
			}
		}
	}
	for _, pile := range b.script.Tableaux() {
		for _, card := range pile.cards {
			if card.Prone() {
				continue
			}
			tail := pile.MakeTail(card)
			if ok, _ := pile.CanMoveTail(tail); !ok {
				continue
			}
			if dst := b.FindHomeForTail(pile, tail); dst != nil {
				if !(dst.Empty() && len(tail) == pile.Len()) {
					if ok, _ := b.script.TailMoveError(tail); ok {
						moves++
					}
				}
			}
		}
	}
	if DebugMode {
		TheUI.SetMiddle(fmt.Sprintf("MOVES: %d", moves))
	}
	return moves == 0
}
