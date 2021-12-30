package sol

//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"fmt"
	"image"

	"oddstream.games/gomps5/util"
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
func (b *Baize) CalcScrunchDims(w, h int) {
	for _, pile := range b.piles {
		switch (pile).(type) {
		case *Tableau, *Reserve:
			switch pile.FanType() {
			case FAN_DOWN:
				// baize->dragOffset is always -ve
				pile.SetScrunchDims(image.Point{X: 0, Y: h - pile.BaizePos().Y + util.Abs(b.dragOffset.Y)})
			case FAN_LEFT:
				pile.SetScrunchDims(image.Point{X: pile.BaizePos().X, Y: 0})
			case FAN_RIGHT:
				// baize->dragOffset is always -ve
				pile.SetScrunchDims(image.Point{X: w - pile.BaizePos().X + util.Abs(b.dragOffset.X), Y: 0})
			}
			pile.SetFanFactor(DefaultFanFactor[pile.FanType()])
		}
	}
}

// CalcFannedRect calculates the width and height this pile would be if it had a specified fan factor
func (self *Core) CalcFannedRect(fanFactor float64) image.Point {
	dims := image.Point{CardWidth, CardHeight}
	if len(self.cards) < 2 {
		return dims
	}
	switch self.fanType {
	case FAN_NONE:
		// well, that was easy
	case FAN_DOWN3:
		switch len(self.cards) {
		case 0, 1:
		case 2:
			dims.Y += int(float64(CardHeight) / CARD_FACE_FAN_FACTOR_V)
		default:
			dims.Y += int(float64(CardHeight)/CARD_FACE_FAN_FACTOR_V) * 2
		}
	case FAN_LEFT3, FAN_RIGHT3:
		switch len(self.cards) {
		case 0, 1:
		case 2:
			dims.X += int(float64(CardWidth) / CARD_FACE_FAN_FACTOR_H)
		default:
			dims.X += int(float64(CardWidth)/CARD_FACE_FAN_FACTOR_H) * 2
		}
	case FAN_DOWN:
		for i := 0; i < len(self.cards)-1; i++ {
			c := self.cards[i]
			if c.Prone() {
				dims.Y += int(float64(CardHeight) / CARD_BACK_FAN_FACTOR)
			} else {
				dims.Y += int(float64(CardHeight) / fanFactor)
			}
		}
	case FAN_LEFT, FAN_RIGHT:
		for i := 0; i < len(self.cards)-1; i++ {
			c := self.cards[i]
			if c.Prone() {
				dims.X += int(float64(CardWidth) / CARD_BACK_FAN_FACTOR)
			} else {
				dims.X += int(float64(CardWidth) / fanFactor)
			}
		}
	}
	return dims
}

func (self *Core) Scrunch() {

	if NoScrunch {
		goto RefanAndReturn
	}

	if len(self.cards) > 2 && (self.scrunchDims.X > CardWidth || self.scrunchDims.Y > CardWidth) {
		var nloops int
		var fanFactor float64
		for fanFactor = DefaultFanFactor[self.fanType]; fanFactor < 7.0; fanFactor += 0.5 {
			dims := self.CalcFannedRect(fanFactor)
			switch self.fanType {
			case FAN_DOWN:
				if dims.Y < self.scrunchDims.Y {
					goto exitloop
				}
			case FAN_LEFT, FAN_RIGHT:
				if dims.X < self.scrunchDims.X {
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
			fmt.Printf("%d loops to go from %f to %f\n", nloops, DefaultFanFactor[self.fanType], self.fanFactor)
		}
	}

RefanAndReturn:
	self.Refan()
}

func (b *Baize) Scrunch() {
	for _, p := range b.piles {
		p.Scrunch()
	}
}
