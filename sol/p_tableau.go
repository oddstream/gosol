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
	"oddstream.games/gosol/util"
)

type Tableau struct {
	pile *Pile
}

func NewTableau(slot image.Point, fanType FanType, moveType MoveType) *Pile {
	tableau := NewPile("Tableau", slot, fanType, moveType)
	tableau.vtable = &Tableau{pile: &tableau}
	TheBaize.AddPile(&tableau)
	return &tableau
}

func powerMoves(piles []*Pile, pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range piles {
		if p.Empty() {
			switch p.vtable.(type) {
			case *Cell:
				emptyCells++
			case *Tableau:
				if p.Label() == "" && p != pDraggingTo {
					// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
					emptyCols++
				}
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	// println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
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
		if ThePreferences.PowerMoves {
			moves := powerMoves(TheBaize.piles, self.pile)
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
	return TheBaize.script.TailAppendError(self.pile, tail)
}

func (self *Tableau) TailTapped(tail []*Card) {
	self.pile.DefaultTailTapped(tail)
}

func (self *Tableau) Conformant() bool {
	return TheBaize.script.UnsortedPairs(self.pile) == 0
}

func (self *Tableau) UnsortedPairs() int {
	return TheBaize.script.UnsortedPairs(self.pile)
}

func (self *Tableau) MovableTails() []*MovableTail {
	var tails []*MovableTail = []*MovableTail{}
	if self.pile.Len() > 0 {
		for _, card := range self.pile.cards {
			var tail = self.pile.MakeTail(card)
			if ok, _ := self.pile.CanMoveTail(tail); ok {
				if ok, _ := TheBaize.script.TailMoveError(tail); ok {
					var homes []*Pile = TheBaize.FindHomesForTail(tail)
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
