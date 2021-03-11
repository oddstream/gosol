package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"bytes"
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

const (
	// lerpStepAmount is the amount a transitioning card moves each tick
	lerpStepAmount = 0.025
	// flipStepAmount is the amount we shrink/grow the flipping card width every tick
	flipStepAmount = 0.1
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

// golang gotcha: go:embed cannot apply to var inside func

//go:embed assets/cards71x96.png
var faceBytes []byte

//go:embed assets/windows_16bit_cards.png
var backBytes []byte

var (
	shadowImage    *ebiten.Image
	faceImageSheet *ebiten.Image
	backImageSheet *ebiten.Image
	backFrames     = map[string]image.Point{
		"Aquarium":    {X: 5, Y: 4},
		"CardHand":    {X: 85, Y: 4},
		"Castle":      {X: 165, Y: 4},
		"JazzCup":     {X: 245, Y: 4},
		"Fishes":      {X: 325, Y: 4},
		"FlowerBlack": {X: 405, Y: 4},
		"FlowerBlue":  {X: 485, Y: 4},
		"PalmBeach":   {X: 5, Y: 140},
		"Pattern1":    {X: 85, Y: 140},
		"Pattern2":    {X: 165, Y: 140},
		"Robot":       {X: 245, Y: 140},
		"Roses":       {X: 325, Y: 140},
		"Shell":       {X: 405, Y: 140},
	}
	scalableFaceImages map[CardID]*ebiten.Image
	scalableBackImage  *ebiten.Image
)

func init() {
	// https://ebiten.org/examples/tiles.html
	// uses
	// screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// i.e. draws the tile direct from a subimage of imagesheet

	// load fixed-size images for retro cards from spritesheets

	println("loading spritesheets")

	img, _, err := image.Decode(bytes.NewReader(faceBytes))
	if err != nil {
		log.Fatal(err)
	}
	faceImageSheet = ebiten.NewImageFromImage(img)
	faceBytes = nil

	img, _, err = image.Decode(bytes.NewReader(backBytes))
	if err != nil {
		log.Fatal(err)
	}
	backImageSheet = ebiten.NewImageFromImage(img)
	backBytes = nil
}

// BuildScalableCardImages builds the card images that can change in scale, after CardWidth,Height have been set
func BuildScalableCardImages() {

	// build shadow image
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetRGBA(0.1, 0.1, 0.1, 0.9)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Fill()
	dc.Stroke()
	shadowImage = ebiten.NewImageFromImage(dc.Image())

	// build images for scalable cards

	scalableFaceImages = make(map[CardID]*ebiten.Image)

	for ord := 1; ord < 14; ord++ {
		for suit := 1; suit < 5; suit++ {
			ID := NewCardID(0, suit, ord)
			scalableFaceImages[ID] = createFaceImage(ID)
		}
	}

	scalableBackImage = createBackImage()
}

// Card object
type Card struct {
	ID    CardID // contains flags (prone, marked) pack, suit, ordinal
	owner *Pile

	screenX, screenY int // current position on screen (after Fan)

	srcX, srcY float64 // smoothstep origin
	dstX, dstY float64 // smoothstep destination
	lerpStep   float64 // current lerp value 0.0 .. 1.0; if < 1.0, card is lerping

	dragging               bool // true if this card is being dragged
	dragStartX, dragStartY int  // starting point for dragging

	flipStep  float64 // if 0, we are not flipping
	flipWidth float64 // scale of the card width while flipping

	shake shakeState

	faceImg, backImg *ebiten.Image // images used to draw this card
}

// NewCard is a factory for Card objects
func NewCard(pack, suit, ordinal int) *Card {
	c := &Card{ID: NewCardID(pack, suit, ordinal)}
	c.SetProne(true)

	switch TheUserData.CardStyle {
	case "retro":
		var faceX, faceY, backX, backY int
		switch c.StringSuit() {
		case "Club":
			faceY = 0
		case "Diamond":
			faceY = CardHeight
		case "Heart":
			faceY = CardHeight + CardHeight
		case "Spade":
			faceY = CardHeight + CardHeight + CardHeight
		}
		faceX = (c.Ordinal() - 1) * CardWidth

		c.faceImg = faceImageSheet.SubImage(image.Rect(faceX, faceY, faceX+CardWidth, faceY+CardHeight)).(*ebiten.Image)
		if c.faceImg == nil {
			log.Fatal("no face image")
		}
		pt := backFrames[TheUserData.CardBack]
		backX, backY = pt.X, pt.Y
		c.backImg = backImageSheet.SubImage(image.Rect(backX, backY, backX+CardWidth, backY+CardHeight)).(*ebiten.Image)
		if c.backImg == nil {
			log.Fatal("no back image")
		}

	case "default", "scalable", "bridge", "poker":
		subid := NewCardID(0, c.Suit(), c.Ordinal())
		var ok bool
		c.faceImg, ok = scalableFaceImages[subid]
		if !ok {
			log.Fatal(subid, "not in scalableFaceImages")
		}
		c.backImg = scalableBackImage
	default:
		log.Fatal("unknown CardStyle", TheUserData.CardStyle)
	}

	// could do c.lerpStep = 1.0 here, but a freshly created card is soon SetPosition()

	return c
}

// String satisfies the Stringer interface (defined by fmt package)
func (c *Card) String() string {
	return c.ID.String()
}

// Position returns the x,y screen coords of this card
func (c *Card) Position() (int, int) {
	return c.screenX, c.screenY
}

// Rect gives the x,y screen coords of the card's top left and bottom right corners
func (c *Card) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = c.screenX
	y0 = c.screenY
	x1 = x0 + CardWidth
	y1 = y0 + CardHeight
	return // using named return parameters
}

// SetPosition sets the position of the Card
func (c *Card) SetPosition(x, y int) {
	c.screenX, c.screenY = x, y
	c.lerpStep = 1.0 // stop any current lerp
}

// TransitionTo starts the transition of this Card
func (c *Card) TransitionTo(x, y int) {
	// if c.lerpStep < 1.0 {
	// 	println(c.ID.String(), "already lerping")
	// }
	if x == c.screenX && y == c.screenY {
		c.SetPosition(x, y)
	} else {
		c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
		c.dstX, c.dstY = float64(x), float64(y)
		c.lerpStep = 0.0 // trigger a lerp
	}
}

// StartDrag informs card that it is being dragged
func (c *Card) StartDrag() {
	c.dragStartX, c.dragStartY = c.screenX, c.screenY
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
	if c.lerpStep < 1.0 || c.dragging {
		return
	}
	// hijack the lerping src, dst positions
	c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
	c.dstX, c.dstY = float64(c.screenX-shakeAmount), float64(c.screenY)
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

// Flip toggles the card
// func (c *Card) Flip() {
// 	if c.prone {
// 		c.FlipUp()
// 	} else {
// 		c.FlipDown()
// 	}
// }

// Animating returns true if this card is lerping, dragging or flipping
func (c *Card) Animating() bool {
	if c.lerpStep < 1.0 {
		return true
	}
	if c.dragging {
		return true
	}
	if c.flipStep != 0 {
		return true
	}
	return false
}

// MarkMovable sets the movable state
func (c *Card) MarkMovable(bool) {
	// TODO
}

// Layout implements ebiten.Game's Layout.
func (c *Card) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the card state (transitions)
func (c *Card) Update() error {
	if c.lerpStep < 1.0 {
		c.screenX = int(util.Smoothstep(c.srcX, c.dstX, c.lerpStep))
		c.screenY = int(util.Smoothstep(c.srcY, c.dstY, c.lerpStep))
		// make Card settle faster when already close to it's destination
		if util.OverlapAreaFloat64(c.srcX, c.srcY, c.srcX+float64(CardWidth), c.srcY+float64(CardHeight), c.dstX, c.dstY, c.dstX+float64(CardWidth), c.dstY+float64(CardHeight)) > 0 {
			c.lerpStep += lerpStepAmount * 2
		} else {
			c.lerpStep += lerpStepAmount
		}
		if c.lerpStep >= 1.0 {
			c.screenX, c.screenY = int(c.dstX), int(c.dstY)
		}
	}
	if c.flipStep != 0.0 {
		c.flipWidth += c.flipStep
		if c.flipWidth <= 0.15 {
			c.flipStep = flipStepAmount // now make card wider
		} else if c.flipWidth >= 1.0 {
			c.flipWidth = 1.0
			c.flipStep = 0.0
		}
	}
	if c.shake != notShaking {
		switch c.shake {
		case shakingLeft:
			if float64(c.screenX) > c.dstX {
				c.screenX--
			} else {
				c.dstX = c.srcX + shakeAmount
				c.shake = shakingRight
			}
		case shakingRight:
			if float64(c.screenX) < c.dstX {
				c.screenX++
			} else {
				c.dstX = c.srcX
				c.shake = shakingCenter
			}
		case shakingCenter:
			if float64(c.screenX) > c.dstX {
				c.screenX--
			} else {
				c.shake = notShaking
			}
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
			img = c.backImg
		}
	} else {
		if c.Prone() {
			img = c.backImg
		} else {
			img = c.faceImg
		}
	}

	if c.flipStep != 0 {
		// img = ebiten.NewImageFromImage(img)
		op.GeoM.Translate(float64(-CardWidth/2), 0)
		op.GeoM.Scale(c.flipWidth, 1.0)
		op.GeoM.Translate(float64(CardWidth/2), 0)
	}

	op.GeoM.Translate(float64(c.screenX), float64(c.screenY))

	if c.flipStep == 0 && (c.lerpStep < 1.0 || c.dragging == true) {
		op.GeoM.Translate(2, 2)
		screen.DrawImage(shadowImage, op)
		op.GeoM.Translate(-2, -2)
	}

	screen.DrawImage(img, op)
}
