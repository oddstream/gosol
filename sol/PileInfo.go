package sol

// PileInfo contains the basic Pile members and a lookup table for it's attributes
type PileInfo struct {
	Class              string           // "Stock"|"StockScorpion"|"StockSpider"|"Waste"|"Foundation"|"Tableau"|"Cell"|"Reserve"|"Golf"
	X, Y               PilePositionType // relative position on Baize in CardWidth/Height units (ie not screen coords)
	Fan                string           // ""|"None"|"Down"|"Right"|"Waste"|"WasteRight"|"WasteDown"
	Build, Drag, Flags int
	Attributes         M
}

// GetIntAttribute gets an integer Pile.PileInfo attribute
func (p *Pile) GetIntAttribute(key string) (int, bool) {
	// str, exists := p.Attributes[key]
	// if !exists {
	// 	return 0, false
	// }
	// i, err := strconv.Atoi(str)
	// if err != nil {
	// 	log.Fatal(str + " is not an int")
	// }
	i, exists := p.Attributes[key]
	if !exists {
		return 0, false
	}
	return i.(int), true
}

// GetStringAttribute gets a string Pile.PileInfo attribute
func (p *Pile) GetStringAttribute(key string) string {
	str, exists := p.Attributes[key]
	if exists {
		return str.(string)
	}
	return ""
}

// GetBoolAttribute gets a boolean Pile.PileInfo attribute
func (p *Pile) GetBoolAttribute(key string) bool {
	// str, exists := p.Attributes[key]
	// if exists {
	// 	value, err := strconv.ParseBool(str)
	// 	if err != nil {
	// 		log.Panic("expecting a bool, got ", str)
	// 	}
	// 	return value
	// }
	b, exists := p.Attributes[key]
	if exists {
		return b.(bool)
	}
	return false
}
