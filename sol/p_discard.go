package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Discard struct {
	Base
}

func NewDiscard(slot image.Point, fanType FanType) *Discard {
	d := &Discard{}
	d.Ctor(d, "Discard", slot, FAN_NONE)
	return d
}

func (*Discard) CanMoveTail(tail []*Card) (bool, error) {
	return false, errors.New("Cannot move cards from a Discard")
}

func (*Discard) CanAcceptCard(card *Card) (bool, error) {
	return false, errors.New("Cannot move a single card to a Discard")
}

func (d *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !d.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	if len(tail) != len(TheBaize.cardLibrary)/len(TheBaize.discards) {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	return TheBaize.script.TailMoveError(tail) // check cards are conformant
}

func (d *Discard) Conformant() bool {
	if d.Len() > 1 {
		return TheBaize.script.UnsortedPairs(d) == 0
	}
	return true
}

func (d *Discard) Complete() bool {
	if d.Empty() {
		return true
	}
	if d.Len() == len(TheBaize.cardLibrary)/len(TheBaize.discards) {
		return true
	}
	return false
}

func (d *Discard) UnsortedPairs() int {
	if d.Len() > 1 {
		return TheBaize.script.UnsortedPairs(d)
	} else {
		return 0
	}
}

func (d *Discard) CreateBackgroundImage() {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	dc.Fill()
	dc.Stroke()
	d.img = ebiten.NewImageFromImage(dc.Image())
}
