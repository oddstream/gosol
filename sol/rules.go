package sol

import (
	"fmt"
	"strings"

	"oddstream.games/gosol/util"
)

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
		if flags&1 == 1 {
			s = s + " Aces are allowed on Kings."
		}
	case 2:
		s = s + " and down, eg a 9 goes on a 10."
		if flags&1 == 1 {
			s = s + " Kings are allowed on Aces."
		}
	case 4:
		s = s + " and either up or down."
		if flags&1 == 1 {
			s = s + " Aces and Kings are allowed on top of each other."
		}
	case 5:
		s = s + " regardless of rank."
	}
	return s
}

func (b *Baize) rulesContents() []string {

	uniquePiles := []string{}
	for _, p := range b.Piles {
		if p.X < 0 || p.Y < 0 {
			continue // don't show rules for hidden piles
		}
		if !util.Contains(uniquePiles, p.Class) {
			uniquePiles = append(uniquePiles, p.Class)
		}
	}

	rules := []string{variantDescription(b.Variant)}

	for _, pileClass := range uniquePiles {
		p := b.findPile(pileClass)
		var str strings.Builder
		switch pileClass {
		case "Stock":
			fmt.Fprintf(&str, "%s: ", "Stock")
			packs, ok := p.GetIntAttribute("Packs")
			if !ok || packs == 0 {
				packs = 1
			}
			if p.X < 0 || p.Y < 0 {
				fmt.Fprintf(&str, "The game uses %s of cards in a hidden stock. ", util.Pluralize("pack", packs))
			} else {
				fmt.Fprintf(&str, "The game uses %s of cards. ", util.Pluralize("pack", packs))
				recycles, _ := p.GetIntAttribute("Recycles")
				if recycles == 0 {
					fmt.Fprint(&str, "The stock cannot be recycled. ")
				} else if recycles > 9000 {
					fmt.Fprint(&str, "The stock can be redealt any numer of times. ")
				} else {
					fmt.Fprintf(&str, "The stock can be redealt %s. ", util.Pluralize("time", recycles))
				}
				targetClass := p.GetStringAttribute("Target")
				if targetClass != "" {
					cardsToMove, ok := p.GetIntAttribute("CardsToMove")
					if !ok {
						cardsToMove = 1
					}
					fmt.Fprintf(&str, "Clicking on the stock will transfer %s to %s.", util.Pluralize("card", cardsToMove), targetClass)
				}

			}
		case "StockSpider":
			fmt.Fprintf(&str, "%s: ", "Stock")
			fmt.Fprint(&str, "Clicking on the stock will transfer one card to each of the tableaux, if all spaces in the tableaux have been filled. ")
		case "StockScorpion":
			fmt.Fprintf(&str, "%s: ", "Stock")
			targetClass := p.GetStringAttribute("Target")
			if targetClass == "" {
				targetClass = "Tableau"
			}
			fmt.Fprintf(&str, "Clicking on the stock will transfer one card to each %s. ", targetClass)
		case "Waste":
			fmt.Fprint(&str, "Waste: Cards can be be moved from here to Cells, Tableaux or Foundations.")
		case "Foundation":
			fmt.Fprint(&str, "Foundation: Build cards ")
			fmt.Fprint(&str, englishRules(p.buildRules, p.buildFlags))
		case "FoundationSpider":
			fmt.Fprint(&str, "Foundation: Build cards ")
			fmt.Fprint(&str, englishRules(p.buildRules, p.buildFlags))
			fmt.Fprint(&str, ". Only a set of 13 cards are allowed to be moved here.")
		case "Tableau":
			if p.buildRules == p.dragRules {
				fmt.Fprint(&str, "Tableau: Build cards ")
				fmt.Fprint(&str, englishRules(p.buildRules, p.buildFlags))
			} else {
				fmt.Fprint(&str, "Tableau: Build cards ")
				fmt.Fprint(&str, englishRules(p.buildRules, p.buildFlags))
				fmt.Fprint(&str, ". Move cards ")
				fmt.Fprint(&str, englishRules(p.dragRules, p.dragFlags))
			}
			accept, ok := p.GetIntAttribute("Accept")
			if !ok {
				accept = 0
			}
			if accept == 0 {
				fmt.Fprint(&str, " Any card may be placed on an empty tableaux.")
			} else if accept > 0 && accept < 14 {
				fmt.Fprintf(&str, " Only a %s may be placed on an empty tableaux.", util.OrdinalToLongString(accept))
			} else {
				fmt.Fprint(&str, " No card may be placed on an empty tableaux.")
			}
			if p.buildFlags&2 == 2 {
				fmt.Fprint(&str, " Strictly, only the top card of each stack may be moved. However, the game automates moves of several cards, when empty tableau columns and empty cells allow.")
			} else {
				if p.dragFlags&1 == 1 {
					fmt.Fprint(&str, " Only a single card may be moved at once.")
				} else {
					fmt.Fprint(&str, " Completed sequences of cards may be moved together.")
				}
			}
			// TODO Bury, Disinter
		case "Cell":
			fmt.Fprint(&str, "Cell: Can store one card of any type.")
		case "Reserve":
			fmt.Fprint(&str, "Reserve: stores multiple cards of any type. You cannot move a card to a reserve.")
		}
		rules = append(rules, str.String())
	}

	return rules
}

func (b *Baize) ShowRules() {
	b.ui.ShowRules(b.rulesContents())
}
