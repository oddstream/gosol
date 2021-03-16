package sol

import (
	"fmt"
	"strings"

	"oddstream.games/gosol/util"
)

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
	rules = append(rules, "This is a very long line of text for testing the word wrap widget to see how it handles longer lines that wrap and everything")

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
		}
		rules = append(rules, str.String())
	}

	return rules
}

func (b *Baize) ShowRules() {
	b.ui.ShowRules(b.rulesContents())
}
