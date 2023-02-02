package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"image"
	"log"
)

type Antares struct {
	ScriptBase
}

/*
	As in FreeCell, there are four cells and four foundation piles.
	At most one card is allowed in each cell.
	The foundation piles are to be built up in suit from Ace to King.
	The game in won when all of the cards are moved here.

	The tableau is divided into two parts: the four left side piles are the Scorpion piles
	and the four right side piles are the FreeCell piles.

	The Scorpion piles are built down in suit.
	Groups of cards (regardless of any sequence) can be moved.
	Fill spaces with Kings or groups of cards headed by a King.
	The FreeCell piles are built down by alternate color.
	Move groups of cards if they are in sequence down by alternate color
	and if there are enough free cells that the cards could be moved individually.
	Spaces can be filled by any card or legal group of cards.

	Groups of cards may be moved from one tableau pile to another,
	if they form a legal sequence in their current pile.
	For example, a sequence down by alternate color
	(if there are a sufficient number of empty cells available to store the cards individually)
	may be moved from the FreeCell piles to the Scorpion piles,
	since that is the legal sequence for the FreeCell piles.
	Also, any group of cards may be moved from the Scorpion piles to the FreeCell piles
	(assuming the head card of the group can be moved there).
	These kinds of groups transfers are called "shifts".
	Shifts of groups down by alternate color are allowed from FreeCell piles to Scorpion piles,
	while shifts of any group of cards are allowed from Scorpion piles to FreeCell piles.
	Once a shift has been made from Scorpion piles to FreeCell piles,
	the cards in the groups cannot again be moved as a group
	unless they are of sequence down by alternate color.
	In short, the Scorpion piles have the same building rules as in the game Scorpion,
	and the FreeCell piles have the same building rules as in the game FreeCell.

	Note that groups of cards in the FreeCell piles of the tableau can only be moved as a group
	if there are a sufficient number of empty cells available to store the cards individually.
	The ability to move cards as a group is only a shortcut to moving the group one card at a time.
	In the Scorpion piles, groups of cards may be moved regardless of any sequence.

	Antares was invented by Thomas Warfield.
*/

func (self *Antares) BuildPiles() {

	self.stock = NewStock(image.Point{5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.cells = nil
	for x := 0; x < 4; x++ {
		self.cells = append(self.cells, NewCell(image.Point{x, 0}))
	}

	self.foundations = nil
	for x := 5; x < 9; x++ {
		f := NewFoundation(image.Point{x, 0})
		self.foundations = append(self.foundations, f)
		f.SetLabel("A")
	}

	self.tableaux = nil
	for x := 0; x < 4; x++ {
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ONE_PLUS))
	}
	for x := 5; x < 9; x++ {
		self.tableaux = append(self.tableaux, NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY))
	}
}

func (self *Antares) StartGame() {

	for i := 0; i < 4; i++ {
		for j := 0; j < 6; j++ {
			MoveCard(self.stock, self.tableaux[i])
		}
	}
	for i := 4; i < 8; i++ {
		for j := 0; j < 7; j++ {
			MoveCard(self.stock, self.tableaux[i])
		}
	}

	TheBaize.SetRecycles(0)

	if DebugMode && self.stock.Len() > 0 {
		log.Println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (self *Antares) inFirstFour(tab *Pile) bool {
	for i := 0; i < 4; i++ {
		if tab == self.tableaux[i] {
			return true
		}
	}
	return false
}

func (self *Antares) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			ok, err := TailConformant(tail, CardPair.Compare_DownAltColor)
			if !ok {
				return ok, err
			}
		}
		// else Scorpion rules - move anything anywhere
	}
	return true, nil
}

func (self *Antares) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	card := tail[0]
	src := card.Owner()
	switch dst.vtable.(type) {
	case *Foundation:
		if dst.Empty() {
			return Compare_Empty(dst, card)
		} else {
			return CardPair{dst.Peek(), card}.Compare_UpSuit()
		}
	case *Tableau:
		if self.inFirstFour(src) {
			ok, err := TailConformant(tail, CardPair.Compare_DownAltColor)
			if !ok {
				return ok, err
			}
		}
		// else Scorpion rules - move anything anywhere
		if dst.Empty() {
			return Compare_Empty(dst, card)
		} else {
			if self.inFirstFour(dst) {
				return CardPair{dst.Peek(), card}.Compare_DownAltColor()
			} else {
				return CardPair{dst.Peek(), card}.Compare_DownSuit()
			}
		}
	}
	return true, nil
}

func (self *Antares) UnsortedPairs(pile *Pile) int {
	switch pile.vtable.(type) {
	case *Tableau:
		if self.inFirstFour(pile) {
			return UnsortedPairs(pile, CardPair.Compare_DownAltColor)
		} else {
			return UnsortedPairs(pile, CardPair.Compare_DownSuit)
		}
	default:
		log.Println("*** eh?", pile.category)
	}
	return 0
}

func (*Antares) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (self *Antares) PileTapped(pile *Pile) {}
