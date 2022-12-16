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

type Cell struct {
	parent *Pile
}

func NewCell(slot image.Point) *Pile {
	cell := NewPile("Cell", slot, FAN_NONE, MOVE_ONE)
	cell.vtable = &Cell{parent: &cell}
	TheBaize.AddPile(&cell)
	return &cell
}

func (self *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	if !self.parent.Empty() {
		return false, errors.New("A Cell can only contain one card")
	}
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Cell")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot move a face down card")
	}
	return true, nil
}

func (self *Cell) TailTapped(tail []*Card) {
	self.parent.DefaultTailTapped(tail)
}

func (*Cell) Conformant() bool {
	return true
}

func (*Cell) UnsortedPairs() int {
	return 0
}

func (self *Cell) MovableTails() []*MovableTail {
	// nb same as Reserve.MovableTails
	var tails []*MovableTail = []*MovableTail{}
	if self.parent.Len() > 0 {
		var card *Card = self.parent.Peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = TheBaize.FindHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &MovableTail{dst: home, tail: tail})
		}
	}
	return tails
}

// Placeholder creates a basic outline
func (self *Cell) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
