package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

type Foundation struct {
	pile *Pile
}

func NewFoundation(slot image.Point) *Pile {
	pile := NewPile("Foundation", slot, FAN_NONE, MOVE_NONE)
	pile.vtable = &Foundation{pile: &pile}
	TheBaize.AddPile(&pile)
	return &pile
}

// CanAcceptTail does some obvious check on the tail before passing it to the script
func (self *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if self.pile.Len() == 13 {
		return false, errors.New("That Foundation already contains 13 cards")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card to a Foundation")
	}
	return TheBaize.script.TailAppendError(self.pile, tail)
}

func (*Foundation) TailTapped([]*Card) {}

func (*Foundation) Conformant() bool {
	return true
}

func (*Foundation) UnsortedPairs() int {
	// you can only put a sorted sequence into a Foundation, so this will always be zero
	return 0
}

func (*Foundation) MovableTails() []*MovableTail {
	return nil
}

func (self *Foundation) Placeholder() *ebiten.Image {
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
