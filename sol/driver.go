package sol

//lint:file-ignore ST1005 the error messages are toasted
// see https://staticcheck.io/docs/configuration

/*
	Driver exists as a behavioural layer between the base Pile class
	and the (sub)'classes' that implement the different types of piles in the game.
	Previously, the code was littered with 'switch Pile.Class' statements
	which started to get smelly and weren't easily maintained or extensible.
	Creating this layer puts the subtype functionality in one place, rather than
	having it scattered throughout the code.

	Pile ... Driver ... geddit?

	I crack myself up with smugness.
*/

import (
	"errors"
	"fmt"
	"strings"

	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/util"
)

// Driver implements different actions on each pile subtype
type Driver interface {
	CanAcceptTail([]*Card) (bool, error)
	CardTapped(*Card) (bool, error)
	Collect() int
	Complete() bool
	Conformant() bool
	English(*strings.Builder)
	Movable() int
	Tapped() bool
}

// Class2NewDriver links the name of a pile subclass with a factory func
var Class2NewDriver = map[string]func(*Pile) Driver{
	"Cell":             NewCell,
	"Foundation":       NewFoundation,
	"FoundationSpider": NewFoundationSpider,
	"Golf":             NewGolf,
	"Reserve":          NewReserve,
	"Stock":            NewStock,
	"StockCruel":       NewStockCruel,
	"StockScorpion":    NewStockScorpion,
	"StockSpider":      NewStockSpider,
	"Tableau":          NewTableau,
	"TableauSpider":    NewTableauSpider,
	"Waste":            NewWaste,
}

// replaced canToast parameter with ok, error return
// Errors.new("my string") or fmt.Errorf("my string", args...)
// this lower-level code has no business displaying toasts

func genericCanAcceptTail(pile *Pile, tail []*Card) (bool, error) {

	if len(tail) == 0 {
		return false, nil
	}
	c0 := tail[0]
	if c0 == nil || c0.owner == pile {
		return false, nil
	}

	if pile.Flags&DragFlagSingle == DragFlagSingle {
		if ThePreferences.PowerMoves {
			pm := powerMoves(TheBaize.Piles, pile)
			if len(tail) > pm {
				return false, fmt.Errorf("Enough free space to move %s, not %d", util.Pluralize("card", pm), len(tail))
			}
		} else {
			if len(tail) > 1 {
				return false, errors.New("You can only drag a single card")
			}
		}
	}
	if pile.Flags&DragFlagSingleOrPile == DragFlagSingleOrPile {
		if !(len(tail) == 1 || len(tail) == c0.owner.CardCount()) {
			return false, errors.New("You can only drag a single card or the whole pile")
		}
	}
	if pile.Empty() {
		if afAttrib := pile.GetStringAttribute("AcceptFrom"); afAttrib != "" {
			afList := strings.Split(afAttrib, ",")
			for _, class := range afList {
				if c0.owner.Class == class {
					return true, nil
				}
			}
			return false, fmt.Errorf("%s can only accept cards from %s", pile.Class, afAttrib)
		}
		if pile.localAccept > 0 {
			return c0.Ordinal() == pile.localAccept, nil
		}
		return true, nil
	}
	return isCardPairConformant(pile.Build, pile.Flags, pile.Peek(), c0), nil
}

func simpleCanAcceptTail(pile *Pile, tail []*Card) (bool, error) {

	if len(tail) == 0 {
		return false, nil
	}
	c0 := tail[0]
	if c0 == nil || c0.owner == pile {
		return false, nil
	}
	if pile.Empty() && pile.localAccept > 0 {
		return c0.Ordinal() == pile.localAccept, nil
	}
	return isCardPairConformant(pile.Build, pile.Flags, pile.Peek(), c0), nil
}

func genericCardTapped(card *Card) (bool, error) {

	var sourcePile *Pile = card.owner
	var tail []*Card = sourcePile.makeTail(card)

	// single top card, try to move to a Foundation
	if card == sourcePile.Peek() {
		for _, fp := range TheBaize.foundations {
			if ok, _ := fp.driver.CanAcceptTail([]*Card{card}); ok {
				fp.MoveCards(card)
				return true, nil
			}
		}
	}

	// else try to move card to the longest Tableau or Golf
	var pLongest *Pile
	for _, p := range TheBaize.Piles {
		if p == card.owner {
			continue
		}
		if !(strings.HasPrefix(p.Class, "Tableau") || p.Class == "Golf") {
			continue
		}
		if ThePreferences.PowerMoves && sourcePile.Flags&DragFlagSingle == DragFlagSingle {
			pm := powerMoves(TheBaize.Piles, p)
			if len(tail) > pm {
				continue
			}
		}
		if ok, _ := p.driver.CanAcceptTail(tail); ok {
			if pLongest == nil || p.CardCount() > pLongest.CardCount() {
				pLongest = p
			}
		}
	}

	if pLongest == nil {
		return false, nil
	}

	pLongest.MoveCards(card)
	return true, nil
}

//

type Cell struct{ parent *Pile }

func NewCell(pile *Pile) Driver {
	return &Cell{parent: pile}
}

func (c *Cell) CanAcceptTail(tail []*Card) (bool, error) {
	// TODO could check localAccept here, but currently no Cells have Accept attribute
	if len(tail) == 1 && c.parent.Empty() {
		return true, nil
	}
	//lint:ignore ST1005 the error message is toasted
	return false, errors.New("You can only have one card in a Cell")
}

func (c *Cell) CardTapped(card *Card) (bool, error) {
	return genericCardTapped(card)
}

func (c *Cell) Complete() bool {
	return len(c.parent.Cards) == 0
}

func (c *Cell) Conformant() bool {
	return true
}

func (c *Cell) Tapped() bool { return false }

//

type Golf struct{ parent *Pile }

func NewGolf(pile *Pile) Driver {
	return &Golf{parent: pile}
}

func (g *Golf) CanAcceptTail(tail []*Card) (bool, error) {
	return simpleCanAcceptTail(g.parent, tail)
}

func (g *Golf) CardTapped(card *Card) (bool, error) {
	return genericCardTapped(card)
}

func (g *Golf) Complete() bool {
	return len(g.parent.Cards) == 0
}

func (g *Golf) Conformant() bool {
	return len(g.parent.Cards) == 0
}

func (g *Golf) Tapped() bool { return false }

//

type Foundation struct{ parent *Pile }

func NewFoundation(pile *Pile) Driver {
	return &Foundation{parent: pile}
}

func (f *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) != 1 {
		return false, nil
	}
	if len(f.parent.Cards) == 13 {
		return false, nil
	}
	return genericCanAcceptTail(f.parent, tail)
}

func (f *Foundation) CardTapped(card *Card) (bool, error) {
	card.Shake()
	return false, errors.New("You cannot move cards from a Foundation")
}

func (f *Foundation) Complete() bool {
	return len(f.parent.Cards) == 13
}

func (f *Foundation) Conformant() bool {
	return true
}

func (f *Foundation) Tapped() bool { return false }

//

type FoundationSpider struct{ Foundation }

func NewFoundationSpider(pile *Pile) Driver {
	return &FoundationSpider{Foundation: Foundation{parent: pile}}
}

func (f *FoundationSpider) CanAcceptTail(tail []*Card) (bool, error) {
	if len(f.parent.Cards) != 0 {
		return false, nil
	}
	if len(tail) != 13 {
		return false, nil
	}
	return isTailConformant(f.parent.Build, f.parent.Flags, tail), nil
}

func (f *FoundationSpider) Complete() bool {
	return len(f.parent.Cards) == 13
}

func (f *FoundationSpider) Conformant() bool {
	return len(f.parent.Cards) == 0 || len(f.parent.Cards) == 13
}

//

type Reserve struct{ parent *Pile }

func NewReserve(pile *Pile) Driver {
	return &Reserve{parent: pile}
}

func (r *Reserve) CanAcceptTail(tail []*Card) (bool, error) {
	//lint:ignore ST1005 the error message is toasted
	return false, errors.New("You cannot move a card to a Reserve pile")
}

func (r *Reserve) CardTapped(card *Card) (bool, error) {
	return genericCardTapped(card)
}

func (r *Reserve) Complete() bool {
	return len(r.parent.Cards) == 0
}

func (r *Reserve) Conformant() bool {
	if len(r.parent.Cards) == 0 {
		return true
	}
	return isTailConformant(r.parent.Build, r.parent.Flags, r.parent.Cards)
}

func (r *Reserve) Tapped() bool { return false }

//

type Stock struct{ parent *Pile }

func NewStock(pile *Pile) Driver {
	return &Stock{parent: pile}
}

func (s *Stock) CanAcceptTail(tail []*Card) (bool, error) {
	//lint:ignore ST1005 the error message is toasted
	return false, errors.New("You cannot move cards to the Stock pile")
}

func (s *Stock) CardTapped(card *Card) (bool, error) {

	sourcePile := card.owner
	var targetClass string = sourcePile.GetStringAttribute("Target")
	var anyCardsMoved bool

	// Tap on a Stock card to send one or more cards to Waste, Golf
	if targetClass == "" {
		targetClass = "Waste"
	}
	cardsToMove, ok := sourcePile.GetIntAttribute("CardsToMove")
	if !ok || cardsToMove == 0 {
		cardsToMove = 1
	}
	// only send cards to one pile, the first targetClass
	// Waste, Golf should always accept a card from Stock,
	// don't need to check if targetClass can accept Stock card
	if targetPile := TheBaize.findPile(targetClass); targetPile != nil {
		for cardsToMove > 0 && card != nil {
			cardsToMove--
			targetPile.MoveCards(card)
			card = sourcePile.Peek()
			anyCardsMoved = true
		}
	} else {
		fmt.Println("Could not find pile", targetPile)
	}

	return anyCardsMoved, nil
}

func (s *Stock) Complete() bool {
	return len(s.parent.Cards) == 0
}

func (s *Stock) Conformant() bool {
	return len(s.parent.Cards) == 0
}

func (pTapped *Stock) Tapped() bool {
	var cardsMoved bool
	if pTapped.parent.localRecycles > 0 {
		waste := TheBaize.findPile("Waste") // TODO varidaic findPile()?
		if waste == nil {
			waste = TheBaize.findPile("Golf")
		}
		if waste == nil || len(waste.Cards) == 0 {
			return false
		}
		// Pop/Push don't play a sound, only MoveCards
		sound.Play("Slide")
		for len(waste.Cards) > 0 {
			c := waste.Pop()
			pTapped.parent.Push(c) // this will flip card down
			cardsMoved = true
		}
		pTapped.parent.SetRecycles(pTapped.parent.localRecycles - 1)
	}
	return cardsMoved
}

//

type StockCruel struct {
	Stock
}

func NewStockCruel(pile *Pile) Driver {
	return &StockCruel{Stock: Stock{parent: pile}}
}

func (pTapped *StockCruel) Tapped() bool {
	/*
	   https://politaire.com/help/cruel

	   The redeal procedure begins by picking up all cards on the tableau.
	   The cards from the tableau are collected, one column at a time, starting with the left-most column,
	   picking up the cards in each column in top to bottom order.
	   Then, without shuffling, the cards are dealt out again, starting with the first card picked up,
	   and dealing the cards in the same order as they were picked up.
	*/

	if pTapped.parent.localRecycles == 0 {
		return false
	}
	tmp := make([]*Card, 0, 52)

	for _, pTab := range TheBaize.Piles {
		if pTab.Class == "Tableau" {
			tmp = append(tmp, pTab.Cards...)
			pTab.Cards = pTab.Cards[:0]
		}
	}
	var cardsMoved bool
	sound.Play("Slide")
	for _, pTab := range TheBaize.Piles {
		if pTab.Class == "Tableau" {
			deal := pTab.GetStringAttribute("Deal")
			for i := 0; i < len(deal); i++ {
				var c *Card
				if len(tmp) > 0 {
					c, tmp = tmp[0], tmp[1:]
				} else {
					goto FinishedDealing
				}
				pTab.Push(c)
				cardsMoved = true
			}
		}
	}
FinishedDealing:
	pTapped.parent.SetRecycles(pTapped.parent.localRecycles - 1)
	return cardsMoved
}

//

type StockScorpion struct {
	Stock
}

func NewStockScorpion(pile *Pile) Driver {
	return &StockScorpion{Stock: Stock{parent: pile}}
}

func (s *StockScorpion) CardTapped(card *Card) (bool, error) {

	var sourcePile *Pile = card.owner
	var anyCardsMoved bool
	for _, p := range TheBaize.Piles {
		if strings.HasPrefix(p.Class, "Tableau") {
			card.SetProne(false)
			p.MoveCards(card)
			anyCardsMoved = true
			card = sourcePile.Peek()
		}
		if card == nil {
			break
		}
	}

	return anyCardsMoved, nil
}

//

type StockSpider struct {
	Stock
}

func NewStockSpider(pile *Pile) Driver {
	return &StockSpider{Stock: Stock{parent: pile}}
}

func (s *StockSpider) CardTapped(card *Card) (bool, error) {

	var tabCards, tabPiles, emptyTabPiles int
	for _, p := range TheBaize.Piles {
		if strings.HasPrefix(p.Class, "Tableau") {
			tabPiles++
			if p.Empty() {
				emptyTabPiles++
			} else {
				tabCards += p.CardCount()
			}
		}
	}
	if tabCards >= tabPiles && emptyTabPiles > 0 {
		if TheBaize.stock.CardCount() > emptyTabPiles {
			return false, errors.New("All empty tableaux must be filled before dealing a new row")
		}
	}

	var sourcePile *Pile = card.owner
	var anyCardsMoved bool
	for _, p := range TheBaize.Piles {
		if strings.HasPrefix(p.Class, "Tableau") {
			card.SetProne(false)
			p.MoveCards(card)
			anyCardsMoved = true
			card = sourcePile.Peek()
		}
		if card == nil {
			break
		}
	}

	return anyCardsMoved, nil
}

//

type Tableau struct{ parent *Pile }

func NewTableau(pile *Pile) Driver {
	return &Tableau{parent: pile}
}

func (t *Tableau) CanAcceptTail(tail []*Card) (bool, error) {
	return genericCanAcceptTail(t.parent, tail)
}

func (t *Tableau) CardTapped(card *Card) (bool, error) {
	return genericCardTapped(card)
}

func (t *Tableau) Complete() bool {
	return len(t.parent.Cards) == 0
}

func (t *Tableau) Conformant() bool {
	if len(t.parent.Cards) == 0 {
		return true
	}
	return isTailConformant(t.parent.Build, t.parent.Flags, t.parent.Cards)
}

func (t *Tableau) Tapped() bool { return false }

//

type TableauSpider struct{ Tableau }

func NewTableauSpider(pile *Pile) Driver {
	return &TableauSpider{Tableau: Tableau{parent: pile}}
}

func (t *TableauSpider) CardTapped(card *Card) (bool, error) {

	var sourcePile *Pile = card.owner
	var tail []*Card = sourcePile.makeTail(card)

	// special case, tapping on the head of a run of 13 cards, see if they'll go to a Foundation
	if len(tail) == 13 && isTailConformant(sourcePile.Build, sourcePile.Flags, tail) {
		for _, p := range TheBaize.Piles {
			if p.Class == "FoundationSpider" {
				if ok, _ := p.driver.CanAcceptTail(tail); ok {
					p.MoveCards(card)
					return true, nil
				}
			}
		}
	}

	return genericCardTapped(card)
}

func (t TableauSpider) Conformant() bool {
	// similar to Tableau.Conformant(), except need 13 conformant cards
	if len(t.parent.Cards) == 0 {
		return true
	}
	return len(t.parent.Cards) == 13 && isTailConformant(t.parent.Build, t.parent.Flags, t.parent.Cards)
}

//

type Waste struct{ parent *Pile }

func NewWaste(pile *Pile) Driver {
	return &Waste{parent: pile}
}

func (w *Waste) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) != 1 {
		//lint:ignore ST1005 the error message is toasted
		return false, errors.New("Only a single card can be put on a Waste pile")
	}
	c0 := tail[0]
	if c0.owner == w.parent {
		return false, nil
	}
	return c0.owner.Class == "Stock", errors.New("Waste can only accept cards from the Stock pile")
}

func (w *Waste) CardTapped(card *Card) (bool, error) {
	return genericCardTapped(card)
}

func (w *Waste) Complete() bool {
	return len(w.parent.Cards) == 0
}

func (w *Waste) Conformant() bool {
	return len(w.parent.Cards) == 0
}

func (w *Waste) Tapped() bool { return false }
