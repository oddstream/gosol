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
	tmp := make([]*Card, 0, cap(p.Cards))
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

func (p *Pile) scrunch(maxSize int, fnCalcSize func(int) int) {
	// https://stackoverflow.com/questions/38897529/pass-method-argument-to-function
	// https://golang.org/ref/spec#Method_values
	// https://play.google.com/books/reader?id=SJHvCgAAQBAJ&pg=GBS.PT217
	// fnCalcSize has implicit receiver value of p *Pile, via a kind of closure, because it was passed as p.fannedWidth/Height
	// could make it explicit with a method expression
	var currSize = fnCalcSize(p.scrunchPercentage)

	// check scrunch if curr height > max height (need more scrunch) or percent scrunch is < 100 (may need less scrunch)
	var scrunchRequired bool = (currSize > maxSize) || (p.scrunchPercentage < 100)
	if !scrunchRequired {
		return
	}

	// println("scrunching", p.Class, "scrunchPercentage now ", p.scrunchPercentage)
	var percent int
	for percent = 100; percent > 50; percent -= 5 {
		testSize := fnCalcSize(percent)
		// println(percent, testSize, maxSize)
		if testSize <= maxSize {
			break
		}
	}
	if percent != p.scrunchPercentage {
		p.scrunchPercentage = percent
		tmp := p.popAllWithoutFlip()
		// println("scrunchPercentage now ", p.scrunchPercentage)
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
		// p.fannedheight is a method value, a function that binds a method (Pile.fannedHeight) to a specific receiver value (p)
		// this function can then be invoked without the receiver value; it needs only the non-receiver arguments
		// it's a kind of closure
		// a method expression would be Pile.fannedHeight or (*Pile).fannedHeight
		// which yields a function value with a regular first parameter taking the place of the receiver
		p.scrunch(s*CardHeight, p.fannedHeight) // method value
	case "Right":
		p.scrunch(s*CardWidth, p.fannedWidth) // method value
	}
}
