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

const (
	// https://en.wikipedia.org/wiki/Miscellaneous_Symbols
	RECYCLE_RUNE   = rune(0x267B)
	NORECYCLE_RUNE = rune(0x2613)
)

/***** ARCHIVED TO REMEMBER HOW TO DO suitFilter *******************
func CreateCardLibrary(packs int, suits int, cardFilter *[14]bool, jokersPerPack int) {

	var numberOfCardsInSuit int = 0
	if cardFilter == nil {
		cardFilter = &[14]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true}
		numberOfCardsInSuit = 13
	} else {
		for i := 1; i < 14; i++ {
			if cardFilter[i] {
				numberOfCardsInSuit++
			}
		}
	}

	var cardsRequired int = packs * suits * numberOfCardsInSuit
	cardsRequired += packs * jokersPerPack
	CardLibrary = make([]Card, 0, cardsRequired)

	for pack := 0; pack < packs; pack++ {
		for suit := 0; suit < suits; suit++ {
			for ord := 1; ord < 14; ord++ {
				if cardFilter[ord] {
					// suits are numbered NOSUIT=0, CLUB=1, DIAMOND=2, HEART=3, SPADE=4
					// (i.e. not 0..3)
					// run the suits loop backwards, so spades are used first
					// (folks expect Spider One Suit to use spades)
					var c Card = NewCard(pack, cardid.SPADE-suit, ord)
					CardLibrary = append(CardLibrary, c)
				}
			}
		}
		for i := 0; i < jokersPerPack; i++ {
			var c Card = NewCard(pack, cardid.NOSUIT, 0) // NOSUIT and ordinal == 0 creates a joker
			CardLibrary = append(CardLibrary, c)
		}
	}
	log.Printf("%d packs, %d suits, %d cards created\n", packs, suits, len(CardLibrary))
}
*********/

type Stock struct {
	pile *Pile
}

func NewStock(slot image.Point, fanType FanType, packs int, suits int, cardFilter *[14]bool, jokersPerPack int) *Pile {
	pile := NewPile("Stock", slot, fanType, MOVE_ONE)
	pile.vtable = &Stock{pile: pile}
	TheGame.Baize.cardCount = pile.Fill(packs, suits)
	pile.Shuffle()
	return pile
}

func (*Stock) CanAcceptTail([]*Card) (bool, error) {
	return false, errors.New("Cannot move cards to the Stock")
}

func (*Stock) TailTapped([]*Card) {
	// do nothing, handled by script, which had first dibs
}

func (self *Stock) Conformant() bool {
	return self.pile.Empty()
}

// UnsortedPairs - cards in a stock pile are always considered to be unsorted
func (self *Stock) UnsortedPairs() int {
	if self.pile.Empty() {
		return 0
	}
	return self.pile.Len() - 1
}

func (self *Stock) MovableTails() []*MovableTail {
	var tails []*MovableTail = []*MovableTail{}
	if self.pile.Len() > 0 {
		var card *Card = self.pile.Peek()
		var tail []*Card = []*Card{card}
		var homes []*Pile = TheGame.Baize.FindHomesForTail(tail)
		for _, home := range homes {
			tails = append(tails, &MovableTail{dst: home, tail: tail})
		}
	}
	return tails
}

func (*Stock) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)

	// farted around trying to use icons for this
	// but they were 48x48 and got fuzzy when scaled
	// and were stubbornly white

	var label rune
	if TheGame.Baize.recycles == 0 {
		label = NORECYCLE_RUNE
	} else {
		label = RECYCLE_RUNE
	}
	dc.SetFontFace(schriftbank.CardSymbolHuge)
	dc.DrawStringAnchored(string(label), float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)

	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
