package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"fmt"
	"log"

	"oddstream.games/gosol/sound"
)

type ScriptBase struct {
	cells       []*Pile
	discards    []*Pile
	foundations []*Pile
	reserves    []*Pile
	stock       *Pile
	tableaux    []*Pile
	waste       *Pile

	wikipedia    string
	cardColors   int
	packs, suits int
}

type Scripter interface {
	BuildPiles()
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(*Pile, []*Card) (bool, error)
	UnsortedPairs(*Pile) int

	TailTapped([]*Card)
	PileTapped(*Pile)

	Cells() []*Pile
	Discards() []*Pile
	Foundations() []*Pile
	Reserves() []*Pile
	Stock() *Pile
	Tableaux() []*Pile
	Waste() *Pile

	Complete() bool
	Wikipedia() string
	CardColors() int
	SafeCollect() bool
	Packs() int
	Suits() int
}

// fallback/default functions for ScriptBase+Scripter /////////////////////////

// no default/fallback for BuildPiles
// no default/fallback for StartGame

func (sb ScriptBase) AfterMove() {}

// no default/fallback for TailMoveError
// no default/fallback for TailAppendError
// no default/fallback for UnsortedPiles

// no default/fallback for TailTapped

func (sb ScriptBase) PileTapped(pile *Pile) {}

func (sb ScriptBase) Cells() []*Pile {
	return sb.cells
}

func (sb ScriptBase) Discards() []*Pile {
	return sb.discards
}

func (sb ScriptBase) Foundations() []*Pile {
	return sb.foundations
}

func (sb ScriptBase) Reserves() []*Pile {
	return sb.reserves
}

func (sb ScriptBase) Stock() *Pile {
	return sb.stock
}

func (sb ScriptBase) Tableaux() []*Pile {
	return sb.tableaux
}

func (sb ScriptBase) Waste() *Pile {
	return sb.waste
}

// Complete - default is number of cards in Foundations == number of cards in CardLibrary.
//
// In Bisley, there may be <13 cards in a Foundation.
// This will need overriding for any variants with Discard piles.
// Could also do this by checking if any pile other than a Foundation is not empty.
func (sb ScriptBase) Complete() bool {
	var n = 0
	for _, f := range sb.foundations {
		n += len(f.cards)
	}
	return n == TheGame.Baize.cardCount
}

// SpiderComplete - used to override default Complete() in Spider varaints.
//
// Each tableau must be either empty or contain a sequence.
// Discard contents must be sequences, otherwise they wouldn't be there.
// There aren't any foundations.
func (sb ScriptBase) SpiderComplete() bool {
	for _, t := range sb.tableaux {
		switch len(t.cards) {
		case 0:
			// that's fine
		case 13: // TODO 104 cards, 8 tabs = 13 cards/tab
			if !t.vtable.Conformant() {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (sb ScriptBase) Wikipedia() string {
	if sb.wikipedia == "" { // uninitialized default
		return "https://en.wikipedia.org/wiki/Patience_(game)"
	} else {
		return sb.wikipedia
	}
}

func (sb ScriptBase) CardColors() int {
	if sb.cardColors == 0 { // uninitialized default
		return 2
	} else {
		return sb.cardColors
	}
}

func (sb ScriptBase) SafeCollect() bool {
	return sb.CardColors() == 2
}

func (sb ScriptBase) Packs() int {
	if sb.packs == 0 {
		return 1
	}
	return sb.packs
}

func (sb ScriptBase) Suits() int {
	if sb.suits == 0 {
		return 4
	}
	return sb.suits
}

// You can't use functions as keys in maps : the key type must be comparable
// so you can't do: var ExtendedColorMap = map[CardPairCompareFunc]bool{}
// type CardPairCompareFunc func(CardPair) (bool, error)

// useful generic game library of functions ///////////////////////////////////

func AnyCardsProne(cards []*Card) bool {
	for _, c := range cards {
		if c.Prone() {
			return true
		}
	}
	return false
}

// MoveCard moves the top card from src to dst
func MoveCard(src *Pile, dst *Pile) *Card {
	if c := src.Pop(); c != nil {
		dst.Push(c)
		src.FlipUpExposedCard()
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
		src.FlipUpExposedCard()
		sound.Play("Place")
	}
}

func RecycleWasteToStock(waste *Pile, stock *Pile) {
	if TheGame.Baize.Recycles() > 0 {
		for waste.Len() > 0 {
			MoveCard(waste, stock)
		}
		TheGame.Baize.SetRecycles(TheGame.Baize.Recycles() - 1)
		switch {
		case TheGame.Baize.recycles == 0:
			TheGame.UI.ToastInfo("No more recycles")
		case TheGame.Baize.recycles == 1:
			TheGame.UI.ToastInfo(fmt.Sprintf("%d recycle remaining", TheGame.Baize.Recycles()))
		case TheGame.Baize.recycles < 10:
			TheGame.UI.ToastInfo(fmt.Sprintf("%d recycles remaining", TheGame.Baize.Recycles()))
		}
	} else {
		TheGame.UI.ToastInfo("No more recycles")
	}
}
