package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

type Tableau struct {
	pile *Pile
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Pile {
	pile := NewPile("Tableau", slot, fanType, moveType)
	pile.vtable = &Tableau{pile: pile}
	return pile
}

func (self *Tableau) CanAcceptTail(tail []*Card) (bool, error) {
	// AnyCardsProne check done by pile.CanMoveTail
	// checking at this level probably isn't needed
	// if AnyCardsProne(tail) {
	// 	return false, errors.New("Cannot add a face down card")
	// }

	// kludge
	// we couldn't check MOVE_PLUS_ONE in pile.CanMoveTail
	// because we didn't then know the destination pile
	// which we need to know to calculate power moves
	if self.pile.moveType == MOVE_ONE_PLUS {
		if TheGame.Settings.PowerMoves {
			moves := TheGame.Baize.powerMoves(self.pile)
			if len(tail) > moves {
				if moves == 1 {
					return false, fmt.Errorf("Space to move 1 card, not %d", len(tail))
				} else {
					return false, fmt.Errorf("Space to move %d cards, not %d", moves, len(tail))
				}
			}
		} else {
			if len(tail) > 1 {
				return false, errors.New("Cannot add more than one card")
			}
		}
	}
	return TheGame.Baize.script.TailAppendError(self.pile, tail)
}

func (self *Tableau) TailTapped(tail []*Card) {
	self.pile.DefaultTailTapped(tail)
}

func (self *Tableau) Conformant() bool {
	// return TheGame.Baize.script.UnsortedPairs(self.pile) == 0
	return self.UnsortedPairs() == 0
}

func (self *Tableau) UnsortedPairs() int {
	return TheGame.Baize.script.UnsortedPairs(self.pile)
}

func (self *Tableau) MovableTails() []*MovableTail {
	var tails []*MovableTail = []*MovableTail{}
	if self.pile.Len() > 0 {
		for _, card := range self.pile.cards {
			var tail = self.pile.MakeTail(card)
			if ok, _ := self.pile.CanMoveTail(tail); ok {
				if ok, _ := TheGame.Baize.script.TailMoveError(tail); ok {
					var homes []*Pile = TheGame.Baize.FindHomesForTail(tail)
					for _, home := range homes {
						tails = append(tails, &MovableTail{dst: home, tail: tail})
					}
				}
			}
		}
	}
	return tails
}

func (self *Tableau) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	if self.pile.label != "" {
		dc.SetFontFace(schriftbank.CardOrdinalLarge)
		dc.DrawStringAnchored(self.pile.label, float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)
	}
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
