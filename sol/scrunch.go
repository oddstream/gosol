package sol

import (
	"image"

	"oddstream.games/gomps5/util"
)

func (b *Baize) FindBuddyPiles() {
	for _, p1 := range b.piles {
		switch (p1.subtype).(type) {
		case *Tableau: // TODO *Reserve
			p1.buddyPos = image.Point{0, 0}
			for _, p2 := range b.piles {
				switch (p2.subtype).(type) {
				case *Tableau:
					switch p1.fanType {
					case FAN_DOWN:
						if p1.slot.X == p2.slot.X && p2.slot.Y > p1.slot.Y {
							p1.buddyPos = p2.pos
						}
					case FAN_LEFT:
						if p1.slot.Y == p2.slot.Y && p2.slot.X < p1.slot.X {
							p1.buddyPos = p2.pos
						}
					case FAN_RIGHT:
						if p1.slot.Y == p2.slot.Y && p2.slot.X > p1.slot.X {
							p1.buddyPos = p2.pos
						}
					}
				}
			}
		}
	}
}

func (b *Baize) CalcScrunchDims(w, h int) {
	for _, p := range b.piles {
		switch (p.subtype).(type) {
		case *Tableau: // TODO *Reserve
			switch p.fanType {
			case FAN_DOWN:
				if p.buddyPos.Y != 0 {
					p.scrunchDims.Y = p.buddyPos.Y - p.pos.Y
				} else {
					// baize->dragOffset is always -ve
					p.scrunchDims.Y = h - p.pos.Y + util.Abs(b.dragOffset.Y)
				}
			case FAN_LEFT:
				if p.buddyPos.X != 0 {
					p.scrunchDims.X = p.buddyPos.X - p.pos.X
				} else {
					p.scrunchDims.X = p.pos.X
				}
			case FAN_RIGHT:
				if p.buddyPos.X != 0 {
					p.scrunchDims.X = p.buddyPos.X - p.pos.X
				} else {
					// baize->dragOffset is always -ve
					p.scrunchDims.X = w - p.pos.X + util.Abs(b.dragOffset.X)
				}
			}
			p.fanFactor = p.defaultFanFactor
		}
	}
}

// CalcFannedRect calculates the width and height this pile would be if it had a specified fan factor
func (p *Pile) CalcFannedRect(fanFactor float64) image.Point {
	dims := image.Point{CardWidth, CardHeight}
	if len(p.cards) < 2 {
		return dims
	}
	switch p.fanType {
	case FAN_NONE:
		// well, that was easy
	case FAN_DOWN3:
		switch len(p.cards) {
		case 0, 1:
		case 2:
			dims.Y += (CardHeight / CARD_FACE_FAN_FACTOR_V)
		default:
			dims.Y += (CardHeight / CARD_FACE_FAN_FACTOR_V) * 2
		}
	case FAN_LEFT3, FAN_RIGHT3:
		switch len(p.cards) {
		case 0, 1:
		case 2:
			dims.X += (CardWidth / CARD_FACE_FAN_FACTOR_H)
		default:
			dims.X += (CardWidth / CARD_FACE_FAN_FACTOR_H) * 2
		}
	case FAN_DOWN:
		for i := 0; i < len(p.cards)-1; i++ {
			c := p.cards[i]
			if c.Prone() {
				dims.Y += CardHeight / CARD_BACK_FAN_FACTOR
			} else {
				dims.Y += int(float64(CardHeight) / fanFactor)
			}
		}
	case FAN_LEFT, FAN_RIGHT:
		for i := 0; i < len(p.cards)-1; i++ {
			c := p.cards[i]
			if c.Prone() {
				dims.X += CardWidth / CARD_BACK_FAN_FACTOR
			} else {
				dims.X += int(float64(CardWidth) / fanFactor)
			}
		}
	}
	return dims
}

func (p *Pile) Scrunch() {

	if NoScrunch {
		p.Refan()
		return
	}

	if len(p.cards) > 2 && (p.scrunchDims.X > CardWidth || p.scrunchDims.Y > CardWidth) {
		var fanFactor float64
		for fanFactor = p.defaultFanFactor; fanFactor < 7.0; fanFactor += 0.5 {
			dims := p.CalcFannedRect(fanFactor)
			switch p.fanType {
			case FAN_DOWN:
				if dims.Y < p.scrunchDims.Y {
					goto exitloop
				}
			case FAN_LEFT, FAN_RIGHT:
				if dims.X < p.scrunchDims.X {
					goto exitloop
				}
			default:
				goto exitloop
			}
			println("going around again", fanFactor)
		}
	exitloop:
		p.fanFactor = fanFactor
	}

	p.Refan()
}

func (b *Baize) Scrunch() {
	for _, p := range b.piles {
		p.Scrunch()
	}
}
