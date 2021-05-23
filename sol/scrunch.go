package sol

func (p *Pile) fannedHeight(scrunchPercent int) int {
	backDelta := CardHeight / backFanFactor * scrunchPercent / 100
	faceDelta := CardHeight / faceFanFactor * scrunchPercent / 100
	var y int
	for i := 0; i < p.CardCount()-1; i++ {
		c := p.Cards[i]
		if c.Prone() {
			y = y + backDelta
		} else {
			y = y + faceDelta
		}
	}
	y = y + CardHeight
	return y
}

func (p *Pile) fannedWidth(scrunchPercent int) int {
	backDelta := CardWidth / backFanFactor * scrunchPercent / 100
	faceDelta := CardWidth / faceFanFactor * scrunchPercent / 100
	var x int
	for i := 0; i < p.CardCount()-1; i++ {
		c := p.Cards[i]
		if c.Prone() {
			x = x + backDelta
		} else {
			x = x + faceDelta
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
	for percent = 100; percent > 50; percent -= 1 {
		testSize := fnCalcSize(percent)
		// println(percent, testSize, maxSize)
		if testSize <= maxSize {
			break
		}
	}
	if percent != p.scrunchPercentage {
		p.scrunchPercentage = percent
		p.RepushAllCards()
	}
}

// ScrunchCards alters the fan so that cards overlap more to fit in view
func (p *Pile) ScrunchCards() {

	if p.scrunchSize == 0 || p.CardCount() < 4 || p.Class != "Tableau" {
		return
	}
	switch p.Fan {
	case "Down":
		// p.fannedheight is a method value, a function that binds a method (Pile.fannedHeight) to a specific receiver value (p)
		// this function can then be invoked without the receiver value; it needs only the non-receiver arguments
		// it's a kind of closure
		// a method expression would be Pile.fannedHeight or (*Pile).fannedHeight
		// which yields a function value with a regular first parameter taking the place of the receiver
		p.scrunch(p.scrunchSize*CardHeight, p.fannedHeight) // method value
	case "Right":
		p.scrunch(p.scrunchSize*CardWidth, p.fannedWidth) // method value
	}

}

// TODO this needs to know height and width of window
// TODO this needs to be called when window dimensions change
func (b *Baize) calcScrunchSizev2() {
	for _, p1 := range b.Piles {
		if p1.Class != "Tableau" {
			continue
		}
		switch p1.Fan {
		case "Down":
			var downwardY PilePositionType
			for _, p2 := range b.Piles {
				if p2.Class != "Tableau" {
					continue
				}
				if p2.X == p1.X && p2.Y > p1.Y {
					downwardY = p2.Y
				}
			}
			if downwardY != 0 {
				p1.scrunchSize = int(downwardY - p1.Y)
			}
		case "Right":
			var rightwardY PilePositionType
			for _, p2 := range b.Piles {
				if p2.Class != "Tableau" {
					continue
				}
				if p2.Y == p1.Y && p2.X > p1.X {
					rightwardY = p2.Y
				}
			}
			if rightwardY != 0 {
				p1.scrunchSize = int(rightwardY - p1.Y)
			}
		}
	}
}
