package dark

import "sort"

var variants = map[string]scripter{}

var variantGroups = map[string][]string{}

// init is used to assemble the "> All" alpha-sorted group of variants
func init() {
	var vnames []string = make([]string, 0, len(variants))
	for k := range variants {
		vnames = append(vnames, k)
	}
	// no need to sort here, sort gets done by func VariantNames()
	variantGroups["> All"] = vnames
	variantGroups["> All by Played"] = vnames
}

func (d *dark) ListVariantGroups() []string {
	var vnames []string = make([]string, 0, len(variantGroups))
	for k := range variantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

func (d *dark) ListVariants(group string) []string {
	var vnames []string = nil
	vnames = append(vnames, variantGroups[group]...)
	if group == "> All by Played" {
		sort.Slice(vnames, func(i, j int) bool {
			return theDark.statistics.Played(vnames[i]) > theDark.statistics.Played(vnames[j])
		})
	} else {
		sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	}
	return vnames
}
