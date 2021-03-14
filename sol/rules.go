package sol

import "oddstream.games/gosol/util"

func (b *Baize) ShowRulesForVariant(v string) {
	rules := []string{variantDescription(v)}

	rules = append(rules, "this is\na multiline\nstring")

	for _, p := range b.Piles {
		if p.X < 0 || p.Y < 0 {
			continue // don't show rules for hidden piles
		}
		if !util.Contains(rules, p.Class) {
			rules = append(rules, p.Class)
			x, y := p.Position()
			b.ui.Toast(x, y, p.Class)
		}
	}
	b.ui.OpenWindow(variantDisplayName(b.Variant), rules)
}

func (b *Baize) ShowRules() {
	b.ShowRulesForVariant(b.Variant)
}
