package sol

import (
	"image"

	"oddstream.games/gomps5/util"
)

func (b *Baize) FindBuddyPiles() {
	for _, p1 := range b.piles {
		switch t1 := p1.(type) {
		case *Tableau:
			t1.buddyPos = image.Point{0, 0}
			for _, p2 := range b.piles {
				switch t2 := p2.(type) {
				case *Tableau:
					switch t1.fanType {
					case FAN_DOWN:
						if t1.slot.X == t2.slot.X && t2.slot.Y > t1.slot.Y {
							t1.buddyPos = t2.pos
						}
					case FAN_LEFT:
						if t1.slot.Y == t2.slot.Y && t2.slot.X < t1.slot.X {
							t1.buddyPos = t2.pos
						}
					case FAN_RIGHT:
						if t1.slot.Y == t2.slot.Y && t2.slot.X > t1.slot.X {
							t1.buddyPos = t2.pos
						}
					}
				}
			}
		}
	}
}

func (b *Baize) CalcScrunchDims(w, h int) {
	for _, p := range b.piles {
		switch tab := p.(type) {
		case *Tableau:
			switch tab.fanType {
			case FAN_DOWN:
				if tab.buddyPos.Y != 0 {
					tab.scrunchDims.Y = tab.buddyPos.Y - tab.pos.Y
				} else {
					// baize->dragOffset is always -ve
					tab.scrunchDims.Y = h - tab.pos.Y + util.Abs(b.dragOffset.Y)
				}
			case FAN_LEFT:
				if tab.buddyPos.X != 0 {
					tab.scrunchDims.X = tab.buddyPos.X - tab.pos.X
				} else {
					tab.scrunchDims.X = tab.pos.X
				}
			case FAN_RIGHT:
				if tab.buddyPos.X != 0 {
					tab.scrunchDims.X = tab.buddyPos.X - tab.pos.X
				} else {
					// baize->dragOffset is always -ve
					tab.scrunchDims.X = w - tab.pos.X + util.Abs(b.dragOffset.X)
				}
			}
			tab.fanFactor = tab.defaultFanFactor
		}
	}
}

// CalcFannedRect calculates the width and height this pile would be if it had a specified fan factor
func (base *Base) CalcFannedRect(fanFactor float64) image.Point {
	dims := image.Point{CardWidth, CardHeight}
	if len(base.cards) < 2 {
		return dims
	}
	switch base.fanType {
	case FAN_NONE:
		// well, that was easy
	case FAN_DOWN3:
		switch len(base.cards) {
		case 0, 1:
		case 2:
			dims.Y += (CardHeight / CARD_FACE_FAN_FACTOR_V)
		default:
			dims.Y += (CardHeight / CARD_FACE_FAN_FACTOR_V) * 2
		}
	case FAN_LEFT3, FAN_RIGHT3:
		switch len(base.cards) {
		case 0, 1:
		case 2:
			dims.X += (CardWidth / CARD_FACE_FAN_FACTOR_H)
		default:
			dims.X += (CardWidth / CARD_FACE_FAN_FACTOR_H) * 2
		}
	case FAN_DOWN:
		for i := 0; i < len(base.cards)-1; i++ {
			c := base.cards[i]
			if c.Prone() {
				dims.Y += CardHeight / CARD_BACK_FAN_FACTOR
			} else {
				dims.Y += int(float64(CardHeight) / fanFactor)
			}
		}
	case FAN_LEFT, FAN_RIGHT:
		for i := 0; i < len(base.cards)-1; i++ {
			c := base.cards[i]
			if c.Prone() {
				dims.X += CardWidth / CARD_BACK_FAN_FACTOR
			} else {
				dims.X += int(float64(CardWidth) / fanFactor)
			}
		}
	}
	return dims
}

func (base *Base) Scrunch() {

	if len(base.cards) > 2 && (base.scrunchDims.X > CardWidth || base.scrunchDims.Y > CardWidth) {
		var fanFactor float64
		for fanFactor = base.defaultFanFactor; fanFactor < 7.0; fanFactor += 0.5 {
			dims := base.CalcFannedRect(fanFactor)
			switch base.fanType {
			case FAN_DOWN:
				if dims.Y < base.scrunchDims.Y {
					goto exitloop
				}
			case FAN_LEFT, FAN_RIGHT:
				if dims.X < base.scrunchDims.X {
					goto exitloop
				}
			default:
				goto exitloop
			}
			println("going around again", fanFactor)
		}
	exitloop:
		base.fanFactor = fanFactor
	}

	base.Refan()
}

func (b *Baize) Scrunch() {
	for _, p := range b.piles {
		p.Scrunch()
	}
}
