package sol

import (
	"fmt"
	"log"
	"strings"

	"oddstream.games/gosol/util"
)

const DragFlagSingle = 1
const DragFlagSingleOrPile = 2
const BuildFlagRankWrap = 4

// retired const BuildFlagSpider = 8

var UpInOnesArray = [14]int{0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1}

var DownInOnesArray = [14]int{0, 13, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

var UpInTwosArray = [14]int{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 1, 2}

func isCardPairConformant(rules, flags int, cPrev, cThis *Card) bool {

	if cPrev == nil {
		return true
	}
	if cPrev.Prone() || cThis.Prone() {
		// println("prone cards are not conformant")
		return false
	}

	localSuit := rules / 10
	localRank := rules % 10

	switch localSuit {
	case 0: // may not build or move
		return false
	case 1: // regardless of suit
	case 2: // in suit
		if cPrev.Suit() != cThis.Suit() {
			return false
		}
	case 3: // in color
		if cPrev.Color() != cThis.Color() {
			return false
		}
	case 4: // in alternate color
		if cPrev.Color() == cThis.Color() {
			return false
		}
	case 5: // in any suit but it's own
		if cPrev.Suit() == cThis.Suit() {
			return false
		}
	}

	if flags&BuildFlagRankWrap == BuildFlagRankWrap { // rank wrap == true
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if UpInOnesArray[cPrev.Ordinal()] != cThis.Ordinal() {
				return false
			}
		case 2: // down, e.g. a 9 goes on a 10
			if DownInOnesArray[cPrev.Ordinal()] != cThis.Ordinal() {
				return false
			}
		case 4: // either up or down
			if !(UpInOnesArray[cPrev.Ordinal()] == cThis.Ordinal() || DownInOnesArray[cPrev.Ordinal()] == cThis.Ordinal()) {
				return false
			}
		case 5: // regardless of rank
		case 6: // rank up in twos (Royal Cotillion)
			if UpInTwosArray[cPrev.Ordinal()] != cThis.Ordinal() {
				return false
			}
		}
	} else { // rank wrap == false
		switch localRank {
		case 0: // may not build or move
			return false
		case 1: // up, e.g. a 10 goes on a 9
			if cThis.Ordinal() != cPrev.Ordinal()+1 {
				return false
			}
		case 2: // down, e.g. a 9 goes on a 10
			if cThis.Ordinal() != cPrev.Ordinal()-1 {
				return false
			}
		case 4: // either up or down
			if util.Abs(cThis.Ordinal()-cPrev.Ordinal()) != 1 {
				return false
			}
		case 5: // regardless of rank
		case 6: // up in twos (Royal Cotillion), useless without rankwrap
			if cThis.Ordinal() != cPrev.Ordinal()+2 {
				return false
			}
		}
	}

	// TODO localRank == 13 (Pyramid) cPrev.Ordinal() + cThis.Ordinal() == 13

	return true
}

func isTailConformant(rules, flags int, cards []*Card) bool {
	if len(cards) == 0 {
		log.Panic("isTailConformant passed empty tail")
		return false
	}
	if rules == 0 {
		return false // may not build or move, even a single card
	}
	cPrev := cards[0]
	for n := 1; n < len(cards); n++ {
		cThis := cards[n]
		if !isCardPairConformant(rules, flags, cPrev, cThis) {
			return false
		}
		cPrev = cThis
	}
	return true
}

func powerMoves(piles []*Pile, pDraggingTo *Pile) int {
	// (1 + number of empty freecells) * 2 ^ (number of empty columns)
	// see http://ezinearticles.com/?Freecell-PowerMoves-Explained&id=104608
	// and http://www.solitairecentral.com/articles/FreecellPowerMovesExplained.html
	var emptyCells, emptyCols int
	for _, p := range piles {
		switch p.Class {
		case "Cell":
			if p.Empty() {
				emptyCells++
			}
		case "Tableau":
			if !p.Empty() || p.localAccept == 99 {
				continue
			}
			// 'If you are moving into an empty column, then the column you are moving into does not count as empty column.'
			if p != pDraggingTo {
				emptyCols++
			}
		}
	}
	// 2^1 == 2, 2^0 == 1, 2^-1 == 0.5
	n := (1 + emptyCells) * util.Pow(2, emptyCols)
	// println(emptyCells, "emptyCells,", emptyCols, "emptyCols,", n, "powerMoves")
	return n
}

func englishRules(rules, flags int) string {
	suit := rules / 10
	rank := rules % 10

	var s string
	switch suit {
	case 0:
		s = "not allowed."
	case 1:
		s = "regardless of suit"
	case 2:
		s = "in suit"
	case 3:
		s = "in color"
	case 4:
		s = "in alternate colors"
	case 5:
		s = "in any other suit"
	}
	switch rank {
	case 0:
	case 1:
		s = s + " and up, eg a 10 goes on a 9."
		if flags&BuildFlagRankWrap == BuildFlagRankWrap {
			s = s + " Aces are allowed on Kings."
		}
	case 2:
		s = s + " and down, eg a 9 goes on a 10."
		if flags&BuildFlagRankWrap == BuildFlagRankWrap {
			s = s + " Kings are allowed on Aces."
		}
	case 4:
		s = s + " and either up or down."
		if flags&BuildFlagRankWrap == BuildFlagRankWrap {
			s = s + " Aces and Kings are allowed on top of each other."
		}
	case 5:
		s = s + " regardless of rank."
	case 6:
		s = s + " and up in twos, eg a J goes on a 9."
		if flags&BuildFlagRankWrap == BuildFlagRankWrap {
			s = s + " Twos are allowed on Kings."
		}
	}
	return s
}

func (x *Cell) English(str *strings.Builder) {
	fmt.Fprint(str, "Cell: Can store one card of any type.")
}
func (x *Foundation) English(str *strings.Builder) {
	p := x.parent
	fmt.Fprint(str, "Foundation: Build cards ")
	fmt.Fprint(str, englishRules(p.Build, p.Flags))
}
func (x *FoundationSpider) English(str *strings.Builder) {
	p := x.parent
	fmt.Fprint(str, "Foundation: Build cards ")
	fmt.Fprint(str, englishRules(p.Build, p.Flags))
	fmt.Fprint(str, " Only a completed set of 13 cards can be moved here.")
}
func (x *Golf) English(str *strings.Builder) {
	fmt.Fprint(str, "Golf: Like a Waste pile, but cards can also be moved here.")
}
func (x *Reserve) English(str *strings.Builder) {
	fmt.Fprint(str, "Reserve: stores multiple cards of any type. You cannot move a card to a reserve.")
}
func (x *Stock) English(str *strings.Builder) {
	p := x.parent
	fmt.Fprintf(str, "%s: ", "Stock")
	packs, ok := p.GetIntAttribute("Packs")
	if !ok || packs == 0 {
		packs = 1
	}
	if p.Hidden() {
		fmt.Fprintf(str, "The game uses %s of cards in a hidden stock. ", util.Pluralize("pack", packs))
	} else {
		fmt.Fprintf(str, "The game uses %s of cards. ", util.Pluralize("pack", packs))
		recycles, _ := p.GetIntAttribute("Recycles")
		if recycles == 0 {
			fmt.Fprint(str, "The stock cannot be recycled. ")
		} else if recycles > 9000 {
			fmt.Fprint(str, "The stock can be redealt any numer of times. ")
		} else {
			fmt.Fprintf(str, "The stock can be redealt %s. ", util.Pluralize("time", recycles))
		}
		targetClass := p.GetStringAttribute("Target")
		if targetClass != "" {
			cardsToMove, ok := p.GetIntAttribute("CardsToMove")
			if !ok {
				cardsToMove = 1
			}
			fmt.Fprintf(str, "Clicking on the stock will transfer %s to %s.", util.Pluralize("card", cardsToMove), targetClass)
		}

	}
}
func (x *StockCruel) English(str *strings.Builder) {
	fmt.Fprint(str, "Clicking on the stock will collect and then redeal the tableaux stacks. ")
}
func (x *StockScorpion) English(str *strings.Builder) {
	p := x.parent
	fmt.Fprintf(str, "%s: ", "Stock")
	targetClass := p.GetStringAttribute("Target")
	if targetClass == "" {
		targetClass = "Tableau"
	}
	fmt.Fprintf(str, "Clicking on the stock will transfer one card to each %s. ", targetClass)
}
func (x *StockSpider) English(str *strings.Builder) {
	fmt.Fprintf(str, "%s: ", "Stock")
	fmt.Fprint(str, "Clicking on the stock will transfer one card to each of the tableaux, if all spaces in the tableaux have been filled. ")
}
func (x *Tableau) English(str *strings.Builder) {
	p := x.parent
	if p.Build == p.Drag {
		fmt.Fprint(str, "Tableau: Build cards ")
		fmt.Fprint(str, englishRules(p.Build, p.Flags))
	} else {
		fmt.Fprint(str, "Tableau: Build cards ")
		fmt.Fprint(str, englishRules(p.Build, p.Flags))
		fmt.Fprint(str, " Move cards ")
		fmt.Fprint(str, englishRules(p.Drag, p.Flags))
	}
	accept, ok := p.GetIntAttribute("Accept")
	if !ok {
		accept = 0
	}
	if accept == 0 {
		fmt.Fprint(str, " Any card may be placed on an empty tableaux.")
	} else if accept > 0 && accept < 14 {
		fmt.Fprintf(str, " Only a %s may be placed on an empty tableaux.", util.OrdinalToLongString(accept))
	} else {
		fmt.Fprint(str, " No card may be placed on an empty tableaux.")
	}
	if p.Flags&DragFlagSingle == DragFlagSingle {
		fmt.Fprint(str, " Only a single card may be moved at once, unless Power Moves is enabled, when the game automates moves of several cards, when empty tableau columns and empty cells allow.")
	} else {
		fmt.Fprint(str, " Completed sequences of cards may be moved together.")
	}
	if bury, ok := p.GetIntAttribute("Bury"); ok {
		buryStr := util.OrdinalToLongString(bury)
		fmt.Fprintf(str, " Any %ss are moved to the bottom of the tableaux when dealing.", buryStr)
	}
	if disinter, ok := p.GetIntAttribute("Disinter"); ok {
		disinterStr := util.OrdinalToLongString(disinter)
		fmt.Fprintf(str, " Any %ss are moved to the top of the tableaux when dealing.", disinterStr)
	}
}
func (x *TableauSpider) English(str *strings.Builder) {
	// tab := (*Tableau)x
	// tab.English(str)
}
func (x *Waste) English(str *strings.Builder) {
	fmt.Fprint(str, "Waste: Cards can be be moved from here to Cells, Tableaux or Foundations.")
}

func (b *Baize) createRulesContents() []string {

	uniquePiles := []string{}
	for _, p := range b.Piles {
		if p.Hidden() {
			continue // don't show rules for hidden piles
		}
		if !util.Contains(uniquePiles, p.Class) {
			uniquePiles = append(uniquePiles, p.Class)
		}
	}

	vi, ok := Variants[b.Variant]
	if !ok {
		log.Fatal("rulesContents unknown variant", b.Variant)
	}
	rules := []string{vi.Description}
	if len(vi.AKA) > 0 {
		rules = append(rules, "AKA: "+strings.Join(vi.AKA, ", "))
	}
	if len(vi.Related) > 0 {
		rules = append(rules, "Related: "+strings.Join(vi.Related, ", "))
	}

	var str strings.Builder
	for _, pileClass := range uniquePiles {
		str.Reset()
		p := b.findPile(pileClass)
		p.driver.English(&str) // TODO just return a string
		rules = append(rules, str.String())
	}

	if len(vi.Wikipedia) > 0 {
		rules = append(rules, vi.Wikipedia)
	}

	// for _, rule := range rules {
	// 	println(rule)
	// }

	return rules
}

func (b *Baize) ShowRules() {
	TheUI.ShowTextDrawer(b.createRulesContents())
}
