package sol

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

const (
	// lerpStepAmount is the amount a transitioning card moves each tick
	lerpStepAmount = 0.025
	// flipStepAmount is the amount we shrink/grow the flipping card width every tick
	flipStepAmount = 0.075
	// shakeAmount the number of pixels to shake a card by
	shakeAmount = 2
)

type shakeState int

const (
	notShaking shakeState = iota
	shakingLeft
	shakingRight
	shakingCenter
)

/*
	Cards have several states: idle, being dragged, transitioning, shaking, spinning, flipping
	You'd think that cards should have a 'state' enum, but the states can overlap (eg a card
	can transition and flip at the same time)
*/

// Card object
type Card struct {
	ID    CardID // contains flags (prone, marked) pack, suit, ordinal
	owner *Pile

	baizeX, baizeY int // current position on baize (after Fan)

	srcX, srcY float64 // smoothstep origin
	dstX, dstY float64 // smoothstep destination
	lerpStep   float64 // current lerp value 0.0 .. 1.0; if < 1.0, card is lerping

	dragging               bool // true if this card is being dragged
	dragStartX, dragStartY int  // starting point for dragging

	flipStep  float64 // if 0, we are not flipping
	flipWidth float64 // scale of the card width while flipping

	shake shakeState

	directionX, directionY int // direction vector when card is spinning
	directionZ, scaleZ     float64
	angle, spin            float64 // current angle and spin when card is spinning

	movable int // 0=not movable, 1=movable but boring, 2=movable and possibly helpful, 3=movable to foundation

	faceImg *ebiten.Image // image used to draw this card
}

// NewCard is a factory for Card objects
func NewCard(pack, suit, ordinal int) *Card {
	c := &Card{ID: NewCardID(pack, suit, ordinal)}
	c.SetProne(true)
	c.RefreshFaceImage()
	// could do c.lerpStep = 1.0 here, but a freshly created card is soon SetPosition()'ed
	return c
}

func (c *Card) RefreshFaceImage() {
	c.faceImg = TheCIP.FaceImage(c.ID)
}

// String satisfies the Stringer interface (defined by fmt package)
func (c *Card) String() string {
	return c.ID.String()
}

// BaizePosition returns the x,y baize coords of this card
func (c *Card) BaizePosition() (int, int) {
	return c.baizeX, c.baizeY
}

// BaizeRect gives the x,y baize coords of the card's top left and bottom right corners
func (c *Card) BaizeRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0 = c.BaizePosition()
	x1 = x0 + CardWidth
	y1 = y0 + CardHeight
	return // using named return parameters
}

// ScreenRect gives the x,y screen coords of the card's top left and bottom right corners
func (c *Card) ScreenRect() (x0 int, y0 int, x1 int, y1 int) {
	x0, y0 = c.BaizePosition()
	x0 += TheBaize.DragOffsetX
	y0 += TheBaize.DragOffsetY
	x1 = x0 + CardWidth
	y1 = y0 + CardHeight
	return // using named return parameters
}

// SetPosition sets the position of the Card
func (c *Card) SetPosition(x, y int) {
	c.baizeX, c.baizeY = x, y
	c.lerpStep = 1.0 // stop any current lerp
}

// TransitionTo starts the transition of this Card
func (c *Card) TransitionTo(x, y int) {
	// if c.lerpStep < 1.0 {
	// 	println(c.ID.String(), "already lerping")
	// }
	if x == c.baizeX && y == c.baizeY {
		c.SetPosition(x, y)
	} else {
		c.srcX, c.srcY = float64(c.baizeX), float64(c.baizeY)
		c.dstX, c.dstY = float64(x), float64(y)
		c.lerpStep = 0.0 // trigger a lerp
	}
}

// StartDrag informs card that it is being dragged
func (c *Card) StartDrag() {
	c.dragStartX, c.dragStartY = c.baizeX, c.baizeY
	c.dragging = true
	// println("start drag", c.ID.String(), "start", c.dragStartX, c.dragStartY)
}

// DragBy repositions the card by the distance it has been dragged
func (c *Card) DragBy(dx, dy int) {
	c.SetPosition(c.dragStartX+dx, c.dragStartY+dy)
}

// DragStartPosition returns the x,y screen coords of this card before dragging started
func (c *Card) DragStartPosition() (int, int) {
	return c.dragStartX, c.dragStartY
}

// StopDrag informs card that it is no longer being dragged
func (c *Card) StopDrag() {
	// println("stop drag", c.ID.String())
	c.dragging = false
}

// CancelDrag informs card that it is no longer being dragged
func (c *Card) CancelDrag() {
	// println("cancel drag", c.ID.String(), "start", c.dragStartX, c.dragStartY, "screen", c.screenX, c.screenY)
	c.TransitionTo(c.dragStartX, c.dragStartY)
	// TODO should go back to Pile.PushedFannedPosition in case of a mis-drag
	c.dragging = false
}

// Shake starts the transition of this Card left, right, center
func (c *Card) Shake() {
	if c.shake != notShaking {
		return
	}
	if c.Transitioning() {
		return
	}
	// hijack the lerping src, dst positions
	c.srcX, c.srcY = float64(c.baizeX), float64(c.baizeY)
	c.dstX, c.dstY = float64(c.baizeX-shakeAmount), float64(c.baizeY)
	c.shake = shakingLeft
}

// FlipUp flips the card face up
func (c *Card) FlipUp() {
	if c.Prone() {
		c.SetProne(false)            // card is immediately face up, else fan isn't correct
		c.flipStep = -flipStepAmount // start by making card narrower
		c.flipWidth = 1.0
	}
}

// FlipDown flips the card face down
func (c *Card) FlipDown() {
	if !c.Prone() {
		c.SetProne(true)             // card is immediately face down, else fan isn't correct
		c.flipStep = -flipStepAmount // start by making card narrower
		c.flipWidth = 1.0
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

// StartSpinning tells the card to start spinning
func (c *Card) StartSpinning() {
	// var coinToss func() int = func() int {
	// 	if rand.Float64() < 0.5 {
	// 		return -1
	// 	}
	// 	return 1
	// }
	c.directionX = rand.Intn(9) - 4
	c.directionY = rand.Intn(9) - 4
	c.directionZ = (rand.Float64() - 0.5) / 100
	c.scaleZ = 1.0
	c.spin = rand.Float64() - 0.5
	c.movable = 0
}

// StopSpinning tells the card to stop spinning and return to it's upright state
func (c *Card) StopSpinning() {
	c.directionX, c.directionY, c.angle, c.spin = 0, 0, 0, 0
	c.scaleZ = 1.0
}

// Spinning returns true if this card is spinning
func (c *Card) Spinning() bool {
	return c.directionX != 0 || c.directionY != 0 || c.angle != 0 || c.spin != 0
}

// Transitioning returns true if this card is lerping
func (c *Card) Transitioning() bool {
	return c.lerpStep < 1.0
}

// Dragging returns true if this card is being dragged
func (c *Card) Dragging() bool {
	return c.dragging
}

// Flipping returns true if this card is flipping
func (c *Card) Flipping() bool {
	return c.flipStep != 0
}

// Layout implements ebiten.Game's Layout.
// func (c *Card) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return outsideWidth, outsideHeight
// }

// Update the card state (transitions)
func (c *Card) Update() error {
	if c.lerpStep < 1.0 {
		c.baizeX = int(util.Smoothstep(c.srcX, c.dstX, c.lerpStep))
		c.baizeY = int(util.Smoothstep(c.srcY, c.dstY, c.lerpStep))
		// make Card settle faster when already close to it's destination
		if util.OverlapAreaFloat64(c.srcX, c.srcY, c.srcX+float64(CardWidth), c.srcY+float64(CardHeight), c.dstX, c.dstY, c.dstX+float64(CardWidth), c.dstY+float64(CardHeight)) > 0 {
			c.lerpStep += lerpStepAmount * 2
		} else {
			c.lerpStep += lerpStepAmount
		}
		if c.lerpStep >= 1.0 {
			c.baizeX, c.baizeY = int(c.dstX), int(c.dstY)
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
	if c.shake != notShaking {
		switch c.shake {
		case shakingLeft:
			if float64(c.baizeX) > c.dstX {
				c.baizeX--
			} else {
				c.dstX = c.srcX + shakeAmount
				c.shake = shakingRight
			}
		case shakingRight:
			if float64(c.baizeX) < c.dstX {
				c.baizeX++
			} else {
				c.dstX = c.srcX
				c.shake = shakingCenter
			}
		case shakingCenter:
			if float64(c.baizeX) > c.dstX {
				c.baizeX--
			} else {
				c.shake = notShaking
			}
		}
	}
	if c.Spinning() {
		c.baizeX += c.directionX
		c.baizeY += c.directionY
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
			img = c.faceImg
		} else {
			// card is getting narrower, and it's going to show face up, but show face down
			img = CardBackImage
		}
	} else {
		if c.Prone() {
			img = CardBackImage
		} else {
			img = c.faceImg
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
		w -= TheBaize.DragOffsetX
		h -= TheBaize.DragOffsetY
		switch {
		case c.baizeX+CardWidth > w:
			c.directionX = -rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.baizeX < 0:
			c.directionX = rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.baizeY+CardHeight > h:
			c.directionY = -rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		case c.baizeY < 0:
			c.directionY = rand.Intn(5)
			c.spin = rand.Float64() - 0.5
		}
	}

	op.GeoM.Translate(float64(c.baizeX+TheBaize.DragOffsetX), float64(c.baizeY+TheBaize.DragOffsetY))

	if !c.Flipping() {
		switch {
		case c.Transitioning():
			offset := float64(CardWidth) / 50.0
			op.GeoM.Translate(offset, offset)
			screen.DrawImage(CardShadowImage, op)
			op.GeoM.Translate(-offset, -offset)
		case c.Dragging():
			offset := float64(CardWidth) / 25.0
			op.GeoM.Translate(offset, offset)
			screen.DrawImage(CardShadowImage, op)
			op.GeoM.Translate(-offset/2, -offset/2)
		}
	}

	if TheUserData.HighlightMovable && !c.Prone() && !c.Spinning() {
		var colorValue float64 = util.MapValue(float64(c.movable), 0, 4, 0.8, 1)
		op.ColorM.Scale(colorValue, colorValue, colorValue, 1)
	}

	screen.DrawImage(img, op)

	// if DebugMode && c.movable > 0 {
	// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", c.movable), c.baizeX, c.baizeY+TheBaize.DragOffsetY)
	// }

	// if c.Movable() {
	// 	screen.DrawImage(CardMovableImage, op)
	// }
}
