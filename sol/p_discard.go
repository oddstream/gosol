package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Discard struct {
	parent *Pile
}

func NewDiscard(slot image.Point, fanType FanType) *Pile {
	discard := NewPile("Discard", slot, FAN_NONE, MOVE_NONE)
	discard.vtable = &Discard{parent: &discard}
	TheBaize.AddPile(&discard)
	return &discard
}

func (self *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.parent.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card to a Discard")
	}
	if len(tail) != len(CardLibrary)/len(TheBaize.script.Discards()) {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	return TheBaize.script.TailMoveError(tail) // check cards are conformant
}

func (*Discard) TailTapped([]*Card) {
	// do nothing
}

func (*Discard) Conformant() bool {
	// no Baize that contains any discard piles should be Conformant,
	// because there is no use showing the collect all FAB
	// because that would do nothing
	// because cards are not collected to discard piles
	return false
}

func (*Discard) UnsortedPairs() int {
	// you can only put a sorted sequence into a Discard, so this will always be zero
	return 0
}

func (self *Discard) MovableTails() []*MovableTail {
	return nil
}

func (self *Discard) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	dc.Fill() // difference for this subpile
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
