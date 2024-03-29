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
	pile *Pile
}

func NewDiscard(slot image.Point, fanType FanType) *Pile {
	pile := NewPile("Discard", slot, FAN_NONE, MOVE_NONE)
	pile.vtable = &Discard{pile: pile}
	return pile
}

func (self *Discard) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.pile.Empty() {
		return false, errors.New("Can only move cards to an empty Discard")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card to a Discard")
	}
	if len(tail) != TheGame.Baize.cardCount/len(TheGame.Baize.script.Discards()) {
		return false, errors.New("Can only move a full set of cards to a Discard")
	}
	if ok, err := TailConformant(tail, CardPair.Compare_DownSuit); !ok {
		return false, err
	}
	// Scorpion tails can always be moved, but Mrs Mop/Simple Simon tails
	// must be conformant, so ...
	return TheGame.Baize.script.TailMoveError(tail)
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
	// you can only put a sequence into a Discard, so this will always be zero
	return 0
}

func (*Discard) MovableTails() []*MovableTail {
	return nil
}

func (*Discard) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	dc.Fill() // difference for this subpile
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
