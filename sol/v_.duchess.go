package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"image"
	"log"

	"oddstream.games/gomps5/util"
)

type Duchess struct {
	ScriptPiles
}

func (du *Duchess) BuildPiles() {
	du.stock = NewStock(image.Point{1, 1}, FAN_NONE, 1, 4, nil)

	for i := 0; i < 4; i++ {
		du.reserves = append(du.reserves, NewReserve(image.Point{i * 2, 0}, FAN_RIGHT))
	}

	du.waste = NewWaste(image.Point{1, 2}, FAN_DOWN3)

	for x := 3; x < 7; x++ {
		du.foundations = append(du.foundations, NewFoundation(image.Point{x, 1}, FAN_NONE))
	}

	for x := 3; x < 7; x++ {
		du.tableaux = append(du.tableaux, NewTableau(image.Point{x, 2}, FAN_DOWN, MOVE_ANY))
	}
}

func (du *Duchess) StartGame() {
	if s, ok := (du.stock.subtype).(*Stock); ok {
		s.recycles = 1
	}
	for _, pile := range du.foundations {
		pile.SetLabel("")
	}
	for _, pile := range du.reserves {
		MoveCard(du.stock, pile)
		MoveCard(du.stock, pile)
		MoveCard(du.stock, pile)
	}

	for _, pile := range du.tableaux {
		MoveCard(du.stock, pile)
	}
	du.stock.SetRune(RECYCLE_RUNE)
	TheUI.Toast("Choose a Reserve card and move it to a Foundation")
}

func (du *Duchess) AfterMove() {
}

func (*Duchess) TailMoveError(tail []*Card) (bool, error) {
	// One card can be moved at a time, but sequences can also be moved as one unit.
	var pile *Pile = tail[0].Owner()
	switch (pile.subtype).(type) {
	case *Tableau:
		var cpairs CardPairs = NewCardPairs(tail)
		// cpairs.Print()
		for _, pair := range cpairs {
			if ok, err := pair.Compare_DownAltColorWrap(); !ok {
				return false, err
			}
		}
	default:
		log.Panic("unknown pile type in TailMoveError")
	}
	return true, nil
}

func (du *Duchess) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst.subtype).(type) {
	case *Stock:
		return false, errors.New("You cannot move cards to the Stock")
	case *Waste:
		c := tail[0]
		if _, ok := (c.owner.subtype).(*Stock); !ok {
			return false, errors.New("Waste can only accept cards from the Stock")
		}
	case *Foundation:
		if dst.Empty() {
			c := tail[0]
			ord := util.OrdinalToShortString(c.Ordinal())
			if dst.label == "" {
				if _, ok := (c.owner.subtype).(*Reserve); !ok {
					return false, errors.New("The first Foundation card must come from a Reserve")
				}
				for _, pile := range du.foundations {
					pile.SetLabel(ord)
				}
			}
			if ord != dst.label {
				return false, fmt.Errorf("Foundations can only accept an %s, not a %s", dst.label, ord)
			}
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_UpSuitWrap()
		}
	case *Tableau:
		if dst.Empty() {
			var rescards int = 0
			for _, p := range du.reserves {
				rescards += p.Len()
			}
			if rescards > 0 {
				// Spaces that occur on the tableau are filled with any top card in the reserve
				c := tail[0]
				if _, ok := (c.owner.subtype).(*Reserve); !ok {
					return false, errors.New("An empty Tableau must be filled from a Reserve")
				}
			}
			return true, nil
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_DownAltColorWrap()
		}
	default:
		log.Panic("unknown pile type in TailAppendError")
	}
	return true, nil
}

func (*Duchess) UnsortedPairs(pile *Pile) int {
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := pair.Compare_DownAltColorWrap(); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

func (du *Duchess) TailTapped(tail []*Card) {
	var pile *Pile = tail[0].Owner()
	if _, ok := (pile.subtype).(*Stock); ok && len(tail) == 1 {
		MoveCard(du.stock, du.waste)
	} else {
		pile.subtype.TailTapped(tail)
	}
}

func (du *Duchess) PileTapped(pile *Pile) {
	if s, ok := (pile.subtype).(*Stock); ok {
		if s.recycles > 0 {
			for du.waste.Len() > 0 {
				MoveCard(du.waste, du.stock)
			}
			s.recycles--
			switch {
			case s.recycles == 0:
				du.stock.SetRune(NORECYCLE_RUNE)
				TheUI.Toast("No more recycles")
			case s.recycles == 1:
				TheUI.Toast(fmt.Sprintf("%d recycle remaining", s.recycles))
			case s.recycles < 10:
				TheUI.Toast(fmt.Sprintf("%d recycles remaining", s.recycles))
			}
		} else {
			TheUI.Toast("No more recycles")
		}
	}
}

func (*Duchess) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Duchess_(solitaire)"
}

func (*Duchess) Discards() []*Pile {
	return nil
}

func (du *Duchess) Foundations() []*Pile {
	return du.foundations
}

func (du *Duchess) Waste() *Pile {
	return du.waste
}
