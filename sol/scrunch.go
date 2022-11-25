package sol

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"fmt"
)

// func (b *Baize) FindBuddyPiles() {
// 	for _, p1 := range b.piles {
// 		switch (p1).(type) {
// 		case *Tableau, *Reserve:
// 			p1.buddyPos = image.Point{0, 0}
// 			for _, p2 := range b.piles {
// 				switch (p2).(type) {
// 				case *Tableau, *Reserve:
// 					switch p1.fanType {
// 					case FAN_DOWN:
// 						if p1.slot.X == p2.slot.X && p2.slot.Y > p1.slot.Y {
// 							p1.buddyPos = p2.pos
// 						}
// 					case FAN_LEFT:
// 						if p1.slot.Y == p2.slot.Y && p2.slot.X < p1.slot.X {
// 							p1.buddyPos = p2.pos
// 						}
// 					case FAN_RIGHT:
// 						if p1.slot.Y == p2.slot.Y && p2.slot.X > p1.slot.X {
// 							p1.buddyPos = p2.pos
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (b *Baize) CalcScrunchDims(w, h int) {
// 	for _, pile := range b.piles {
// 		switch (pile).(type) {
// 		case *Tableau, *Reserve:
// 			switch pile.FanType() {
// 			case FAN_DOWN:
// 				if pile.buddyPos.Y != 0 {
// 					pile.scrunchDims.Y = pile.buddyPos.Y - pile.pos.Y
// 				} else {
// 					// baize->dragOffset is always -ve
// 					pile.scrunchDims.Y = h - pile.pos.Y + util.Abs(b.dragOffset.Y)
// 				}
// 			case FAN_LEFT:
// 				if pile.buddyPos.X != 0 {
// 					pile.scrunchDims.X = pile.buddyPos.X - pile.pos.X
// 				} else {
// 					pile.scrunchDims.X = pile.pos.X
// 				}
// 			case FAN_RIGHT:
// 				if pile.buddyPos.X != 0 {
// 					pile.scrunchDims.X = pile.buddyPos.X - pile.pos.X
// 				} else {
// 					// baize->dragOffset is always -ve
// 					pile.scrunchDims.X = w - pile.pos.X + util.Abs(b.dragOffset.X)
// 				}
// 			}
// 			pile.fanFactor = DefaultFanFactor[pile.fanType]
// 		}
// 	}
// }

// SizeWithFanFactor calculates the width or height this pile would be if it had a specified fan factor
func (self *Pile) SizeWithFanFactor(fanFactor float64) int {
	var max int
	switch self.fanType {
	case FAN_DOWN:
		for i := 0; i < len(self.cards)-1; i++ {
			c := self.cards[i]
			if c.Prone() {
				max += int(float64(CardHeight) / CARD_BACK_FAN_FACTOR)
			} else {
				max += int(float64(CardHeight) / fanFactor)
			}
		}
		max += CardHeight
	case FAN_LEFT, FAN_RIGHT:
		for i := 0; i < len(self.cards)-1; i++ {
			c := self.cards[i]
			if c.Prone() {
				max += int(float64(CardWidth) / CARD_BACK_FAN_FACTOR)
			} else {
				max += int(float64(CardWidth) / fanFactor)
			}
		}
		max += CardWidth
	}
	return max
}

// Scrunch prepares to refan cards after Push() or Pop(), adjusting the amount of overlap to try to keep them fitting on the screen
// only Scrunch piles with fanType LEFT/RIGHT/UP/DOWN, ignore the waste-style piles and those that do not fan
func (self *Pile) Scrunch() {

	self.fanFactor = DefaultFanFactor[self.fanType]

	if NoScrunch || len(self.cards) < 2 {
		self.Refan()
		return
	}

	var maxPileSize int
	switch self.fanType {
	case FAN_DOWN:
		// baize->dragOffset is always -ve
		// statusbar height is 24
		// maxPileSize = TheBaize.WindowHeight - scpos.Y + util.Abs(TheBaize.dragOffset.Y)
		maxPileSize = TheBaize.WindowHeight - self.ScreenPos().Y + (CardHeight / 2)
	case FAN_LEFT:
		maxPileSize = self.ScreenPos().X
	case FAN_RIGHT:
		// baize->dragOffset is always -ve
		// maxPileSize = TheBaize.WindowWidth - scpos.X + util.Abs(TheBaize.dragOffset.X)
		maxPileSize = TheBaize.WindowWidth - self.ScreenPos().X
	}
	if maxPileSize == 0 {
		// this pile doesn't need scrunching
		self.Refan()
		return
	}

	var nloops int
	var fanFactor float64
	for fanFactor = DefaultFanFactor[self.fanType]; fanFactor < 7.0; fanFactor += 0.1 {
		size := self.SizeWithFanFactor(fanFactor)
		switch self.fanType {
		case FAN_DOWN:
			if size < maxPileSize {
				goto exitloop
			}
		case FAN_LEFT, FAN_RIGHT:
			if size < maxPileSize {
				goto exitloop
			}
		default:
			goto exitloop
		}
		nloops++
	}
exitloop:
	self.fanFactor = fanFactor
	if DebugMode && nloops > 0 {
		fmt.Printf("%d loops to go from %f to %f", nloops, DefaultFanFactor[self.fanType], self.fanFactor)
		fmt.Printf(" WindowWidth, Height = %d,%d\n", TheBaize.WindowWidth, TheBaize.WindowHeight)
	}
	self.Refan()
}
