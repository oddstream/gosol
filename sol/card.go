package sol

import (
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

const (
	cardmagic = 0x29041962

	lerpStartNormal = 0.1
	lerpStartClose  = 0.25

	debugSpeed  = 0.005
	slowSpeed   = 0.01
	normalSpeed = 0.02
	fastSpeed   = 0.04

	flipStepAmount = 0.075 // flipStepAmount is the amount we shrink/grow the flipping card width every tick
)

var CardStartPoint image.Point = image.Point{400, -100}

/*
	Cards have several states: idle, being dragged, transitioning, shaking, spinning, flipping
	You'd think that cards should have a 'state' enum, but the states can overlap (eg a card
	can transition and flip at the same time)
*/

// Card object
type Card struct {
	// static things
	magic uint32
	ID    CardID // contains pack, ordinal, suit, ordinal (and bonus prone and joker flag bits)

	// dynamic things
	owner Pile

	pos            image.Point
	src            image.Point // lerp origin
	dst            image.Point // lerp destination
	lerpStep       float64     // current lerp value 0.0 .. 1.0; if < 1.0, card is lerping
	lerpStepAmount float64     // the amount a transitioning card moves each tick

	dragStart image.Point // starting point for dragging

	flipStep  float64 // if 0, we are not flipping
	flipWidth float64 // scale of the card width while flipping

	directionX, directionY int // direction vector when card is spinning
	directionZ, scaleZ     float64
	angle, spin            float64 // current angle and spin when card is spinning

	movable bool
}

// NewCard is a factory for Card objects
func NewCard(pack, suit, ordinal int) Card {
	c := Card{magic: cardmagic, ID: NewCardID(pack, suit, ordinal), pos: CardStartPoint}
	// a joker ID will be created by having NOSUIT (0) and ordinal == 0
	c.SetProne(true)
	// could do c.lerpStep = 1.0 here, but a freshly created card is soon SetPosition()'ed
	return c
}

func (c *Card) Valid() bool {
	return c != nil && c.magic == cardmagic
}

func (c *Card) SetOwner(p Pile) {
	// p may be nil if we have just popped the card
	c.owner = p
}

func (c *Card) Owner() Pile {
	return c.owner
}

// String satisfies the Stringer interface (defined by fmt package)
func (c *Card) String() string {
	return c.ID.String()
}

// Pos returns the x,y baize coords of this card
func (c *Card) BaizePos() image.Point {
	return c.pos
}

// SetPosition sets the position of the Card
func (c *Card) SetBaizePos(pos image.Point) {
	c.pos = pos
	c.lerpStep = 1.0 // stop any current lerp
}

// Rect gives the x,y baize coords of the card's top left and bottom right corners
func (c *Card) BaizeRect() image.Rectangle {
	var r image.Rectangle
	r.Min = c.pos
	r.Max = r.Min.Add(image.Point{CardWidth, CardHeight})
	return r
}

// ScreenRect gives the x,y screen coords of the card's top left and bottom right corners
func (c *Card) ScreenRect() image.Rectangle {
	var r image.Rectangle = c.BaizeRect()
	r.Min = r.Min.Add(TheBaize.dragOffset)
	r.Max = r.Max.Add(TheBaize.dragOffset)
	return r
}

// TransitionTo starts the transition of this Card
func (c *Card) TransitionTo(pos image.Point) {
	// if c.lerpStep < 1.0 {
	// 	println(c.ID.String(), "already lerping")
	// }
	if NoCardLerp || pos.Eq(c.pos) {
		c.SetBaizePos(pos)
		return
	}

	if pos.Eq(c.dst) && c.lerpStep < 1.0 {
		// println("repeat", c.String(), "from", c.pos.String(), "to", c.dst.String())
		return
	}

	if c.Spinning() {
		return
	}

	c.src = c.pos
	c.dst = pos
	// the further the card has to travel, the smaller the lerp step amount
	// eg when dropping a card onto a pile
	// or moving a card from stock to waste
	if c.src.X < 0 || c.src.Y < 0 {
		c.lerpStep = lerpStartNormal
		c.lerpStepAmount = slowSpeed
	} else {
		dist := util.Distance(c.src, c.dst)
		if int(dist) < CardWidth {
			c.lerpStep = lerpStartClose
			c.lerpStepAmount = normalSpeed
			// println("fast", dist, c.String())
		} else {
			c.lerpStep = lerpStartNormal
			c.lerpStepAmount = normalSpeed
		}
	}
}

// StartDrag informs card that it is being dragged
func (c *Card) StartDrag() {
	if c.Transitioning() {
		println("Dragging a transitioning card", c.String())
		// set the drag origin to the be transition destination,
		// so that cancelling this drag will return the card
		// to where it thought it was going
		c.dragStart = c.dst
	} else {
		c.dragStart = c.pos
	}
	// println("start drag", c.ID.String(), "start", c.dragStartX, c.dragStartY)
}

// DragBy repositions the card by the distance it has been dragged
func (c *Card) DragBy(dx, dy int) {
	// println("Card.DragBy(", c.dragStartX+dx-c.baizeX, c.dragStartY+dy-c.baizeY, ")")
	c.SetBaizePos(c.dragStart.Add(image.Point{dx, dy}))
}

// DragStartPosition returns the x,y screen coords of this card before dragging started
// func (c *Card) DragStartPosition() (int, int) {
// return c.dragStartX, c.dragStartY
// }

// StopDrag informs card that it is no longer being dragged
func (c *Card) StopDrag() {
	// println("stop drag", c.ID.String())
}

// CancelDrag informs card that it is no longer being dragged
func (c *Card) CancelDrag() {
	// println("cancel drag", c.ID.String(), "start", c.dragStartX, c.dragStartY, "screen", c.screenX, c.screenY)
	c.TransitionTo(c.dragStart)
}

// WasDragged returns true of this card has been dragged
func (c *Card) WasDragged() bool {
	return !c.pos.Eq(c.dragStart)
}

// FlipUp flips the card face up
func (c *Card) FlipUp() {
	if c.Prone() {
		c.SetProne(false) // card is immediately face up, else fan isn't correct
		if !NoCardFlip {
			c.flipStep = -flipStepAmount // start by making card narrower
			c.flipWidth = 1.0
		}
	}
}

// FlipDown flips the card face down
func (c *Card) FlipDown() {
	if !c.Prone() {
		c.SetProne(true) // card is immediately face down, else fan isn't correct
		if !NoCardFlip {
			c.flipStep = -flipStepAmount // start by making card narrower
			c.flipWidth = 1.0
		}
	}
}

// Flip turns the card over
func (c *Card) Flip() {
	if c.Prone() {
		c.FlipUp()
	} else {
		c.FlipDown()
	}
}

// SetFlip turns the card over
func (c *Card) SetFlip(prone bool) {
	if prone {
		c.FlipDown()
	} else {
		c.FlipUp()
	}
}

// StartSpinning tells the card to start spinning
func (c *Card) StartSpinning() {
	c.directionX = rand.Intn(9) - 4
	c.directionY = rand.Intn(9) - 4
	c.directionZ = (rand.Float64() - 0.5) / 100
	c.scaleZ = 1.0
	c.spin = rand.Float64() - 0.5
}

// StopSpinning tells the card to stop spinning and return to it's upright state
func (c *Card) StopSpinning() {
	c.directionX, c.directionY = 0, 0
	c.angle, c.spin = 0, 0
	c.scaleZ = 1.0
}

// Spinning returns true if this card is spinning
func (c *Card) Spinning() bool {
	return c.directionX != 0 || c.directionY != 0 || c.angle != 0 || c.spin != 0
}

// Transitioning returns true if this card is lerping
func (c *Card) Transitioning() bool {
	// return c.lerpStep < 1.0
	return c.lerpStep < 1.0 && c.pos != c.dst
}

// Dragging returns true if this card is being dragged
func (c *Card) Dragging() bool {
	// if TheBaize.tail == nil {
	// 	return false
	// }
	for _, card := range TheBaize.tail {
		if card == c {
			return true
		}
	}
	return false
	// return c.dragging
}

// Flipping returns true if this card is flipping
func (c *Card) Flipping() bool {
	if NoCardFlip {
		return false
	}
	return c.flipStep != 0.0
}

// Layout implements ebiten.Game's Layout.
// func (c *Card) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return outsideWidth, outsideHeight
// }

// Update the card state (transitions)
func (c *Card) Update() error {
	if c.pos != c.dst && c.lerpStep < 1.0 {
		c.pos.X = int(util.Smoothstep(float64(c.src.X), float64(c.dst.X), c.lerpStep))
		c.pos.Y = int(util.Smoothstep(float64(c.src.Y), float64(c.dst.Y), c.lerpStep))
		if c.lerpStep += c.lerpStepAmount; c.lerpStep >= 1.0 {
			c.pos = c.dst
			c.src = image.Point{0, 0}
		}
	}
	if c.Flipping() {
		c.flipWidth += c.flipStep
		if c.flipWidth <= 0 {
			c.flipStep = flipStepAmount // now make card wider
		} else if c.flipWidth >= 1.0 {
			c.flipWidth = 1.0
			c.flipStep = 0.0
		}
	}
	if c.Spinning() {
		c.pos.X += c.directionX
		c.pos.Y += c.directionY
		c.scaleZ += c.directionZ
		if c.scaleZ < 0.5 || c.scaleZ > 1.5 {
			c.directionZ = -c.directionZ
		}
		c.angle += c.spin
		if c.angle > 360 {
			c.angle -= 360
		} else if c.angle < 0 {
			c.angle += 360
		}
	}
	return nil
}

// Draw renders the card into the screen
func (c *Card) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	var img *ebiten.Image
	// card prone has already been set to destination state
	if c.flipStep < 0 {
		if c.Prone() {
			// card is getting narrower, and it's going to show face down, but show face up
			img = TheCardFaceImageLibrary[(c.Suit()*13)+(c.Ordinal()-1)]
		} else {
			// card is getting narrower, and it's going to show face up, but show face down
			img = CardBackImage
		}
	} else {
		if c.Prone() {
			img = CardBackImage
		} else {
			img = TheCardFaceImageLibrary[(c.Suit()*13)+(c.Ordinal()-1)]
		}
	}

	if c.Flipping() {
		// img = ebiten.NewImageFromImage(img)
		op.GeoM.Translate(float64(-CardWidth/2), 0)
		op.GeoM.Scale(c.flipWidth, 1.0)
		op.GeoM.Translate(float64(CardWidth/2), 0)
	}

	if c.Spinning() {
		// do this before the baize position translate
		op.GeoM.Translate(float64(-CardWidth/2), float64(-CardHeight/2))
		op.GeoM.Rotate(c.angle * 3.1415926535 / 180.0)
		op.GeoM.Scale(c.scaleZ, c.scaleZ)
		op.GeoM.Translate(float64(CardWidth/2), float64(CardHeight/2))

		// naughty to do this here, but Draw knows the screen dimensions and Update doesn't
		w, h := screen.Size()
		w -= TheBaize.dragOffset.X
		h -= TheBaize.dragOffset.Y
		switch {
		case c.pos.X+CardWidth > w:
			c.directionX = -rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.pos.X < 0:
			c.directionX = rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.pos.Y+CardHeight > h:
			c.directionY = -rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.pos.Y < 0:
			c.directionY = rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		}
	}

	op.GeoM.Translate(float64(c.pos.X+TheBaize.dragOffset.X), float64(c.pos.Y+TheBaize.dragOffset.Y))

	if CardShadowImage != nil {
		if !c.Flipping() {
			switch {
			case c.Transitioning():
				offset := float64(CardWidth) / 20.0
				op.GeoM.Translate(offset, offset)
				screen.DrawImage(CardShadowImage, op)
				offset = -offset
				op.GeoM.Translate(offset, offset)
			case c.Dragging():
				offset := float64(CardWidth) / 20.0
				op.GeoM.Translate(offset, offset)
				screen.DrawImage(CardShadowImage, op)
				// move the offset PARTIALLY back, making the card appear "pressed" when pushed with the mouse (like a button)
				offset = -offset * 0.5
				op.GeoM.Translate(offset, offset)
				// this looks intuitively better than "lifting" the card with
				// op.GeoM.Translate(-offset*2, -offset*2)
				// even though "lifting" it (moving it up/left towards the light source) would be more "correct"
			}
		}
	}

	if img == nil {
		log.Panic("Card.Draw no image for ", c.String(), " prone: ", c.Prone())
	}

	if c.Owner().Target() && c == c.Owner().Peek() {
		op.ColorM.Scale(0.9, 0.9, 0.9, 1)
	}

	if DebugMode && ThePreferences.MarkMovableCards && c.movable {
		op.ColorM.Scale(0.9, 0.9, 0.9, 1)
	}

	screen.DrawImage(img, op)
}
