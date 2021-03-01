package sol

func (p *Pile) popWithoutFlip() *Card {
	if 0 == p.CardCount() {
		return nil
	}
	c := p.Cards[p.CardCount()-1]
	p.Cards = p.Cards[:p.CardCount()-1]
	return c
}

func (p *Pile) popAllWithoutFlip() []*Card {
	var tmp []*Card
	for p.CardCount() > 0 {
		c := p.popWithoutFlip()
		tmp = append(tmp, c)
	}
	return tmp
}

func (p *Pile) pushWithoutFlip(c *Card) {
	c.TransitionTo(p.PushedFannedPosition()) // do this BEFORE appending card to pile
	p.Cards = append(p.Cards, c)
}

func (p *Pile) pushAllWithoutFlip(tmp []*Card) {
	for len(tmp) > 0 {
		c := tmp[len(tmp)-1]
		tmp = tmp[:len(tmp)-1]
		p.pushWithoutFlip(c)
	}
}

func (p *Pile) fannedHeight(scrunchPercent int) int {
	var y int
	for i := 0; i < p.CardCount()-1; i++ {
		c := p.Cards[i]
		if c.prone {
			y = y + (CardHeight / backFanFactor * scrunchPercent / 100)
		} else {
			y = y + (CardHeight / faceFanFactor * scrunchPercent / 100)
		}
	}
	y = y + CardHeight
	return y
}

func (p *Pile) fannedWidth(scrunchPercent int) int {
	var x int
	for i := 0; i < p.CardCount()-1; i++ {
		c := p.Cards[i]
		if c.prone {
			x = x + (CardWidth / backFanFactor * scrunchPercent / 100)
		} else {
			x = x + (CardWidth / faceFanFactor * scrunchPercent / 100)
		}
	}
	x = x + CardWidth
	return x
}

func (p *Pile) scrunchCardsDown(s int) {
	var currHeight = p.fannedHeight(p.scrunchPercentage)
	var maxHeight = s * CardHeight

	// check scrunch if curr height > max height (need more scrunch) or percent scrunch is < 100 (may need less scrunch)
	var scrunchRequired bool = (currHeight > maxHeight) || (p.scrunchPercentage < 100)
	if !scrunchRequired {
		return
	}

	println("scrunching", p.Class, "scrunchPercentage now ", p.scrunchPercentage)
	var percent int
	for percent = 100; percent > 50; percent -= 5 {
		testHeight := p.fannedHeight(percent)
		println(percent, testHeight, maxHeight)
		if testHeight <= maxHeight {
			break
		}
	}
	if percent != p.scrunchPercentage {
		p.scrunchPercentage = percent
		tmp := p.popAllWithoutFlip()
		println("scrunchPercentage now ", p.scrunchPercentage)
		p.pushAllWithoutFlip(tmp)
	}
}

func (p *Pile) scrunchCardsRight(s int) {
	var currWidth = p.fannedWidth(p.scrunchPercentage)
	var maxWidth = s * CardWidth

	// check scrunch if curr height > max height (need more scrunch) or percent scrunch is < 100 (may need less scrunch)
	var scrunchRequired bool = (currWidth > maxWidth) || (p.scrunchPercentage < 100)
	if !scrunchRequired {
		return
	}

	println("scrunching", p.Class, "scrunchPercentage now ", p.scrunchPercentage)
	var percent int
	for percent = 100; percent > 50; percent -= 5 {
		testHeight := p.fannedWidth(percent)
		println(percent, testHeight, maxWidth)
		if testHeight <= maxWidth {
			break
		}
	}
	if percent != p.scrunchPercentage {
		p.scrunchPercentage = percent
		tmp := p.popAllWithoutFlip()
		println("scrunchPercentage now ", p.scrunchPercentage)
		p.pushAllWithoutFlip(tmp)
	}
}

// ScrunchCards alters the fan so that cards overlap more to fit in view
func (p *Pile) ScrunchCards() {

	if p.CardCount() < 3 {
		return
	}
	s, ok := p.GetIntAttribute("Scrunch")
	if !ok {
		return
	}
	switch p.Fan {
	case "Down":
		p.scrunchCardsDown(s)
	case "Right":
		p.scrunchCardsRight(s)
	}
}
