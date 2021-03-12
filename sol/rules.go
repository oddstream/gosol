package sol

import "oddstream.games/gosol/util"

func (b *Baize) ShowRulesForVariant(v string) {
	rules := []string{variantDescription(v)}

	rules = append(rules, "this is\na multiline\nstring")

	for _, p := range b.Piles {
		if !util.Contains(rules, p.Class) {
			rules = append(rules, p.Class)
		}
	}
	b.ui.OpenWindow(b.input, variantDisplayName(b.Variant), rules)
}

func (b *Baize) ShowRules() {
	b.ShowRulesForVariant(b.Variant)
}
