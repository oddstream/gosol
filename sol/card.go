package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"strconv"

	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

const (
	// FLIPSTEP is the amount we shrink/grow the flipping card width every tick
	FLIPSTEP = 0.1
	// SHAKEAMOUNT the number of pixels to shake a card by
	SHAKEAMOUNT = 2
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
)

func init() {
	// https://ebiten.org/examples/tiles.html
	// uses
	// screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// i.e. draws the tile direct from a subimage of imagesheet

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

	dc := gg.NewContext(71, 96)
	dc.SetRGBA(0.1, 0.1, 0.1, 0.9)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(0, 0, float64(71), float64(96), 4)
	dc.Fill()
	dc.Stroke()
	shadowImage = ebiten.NewImageFromImage(dc.Image())
}

// Card object
type Card struct {
	owner *Pile

	pack    int
	suit    string
	ordinal int
	prone   bool
	id      string
	color   color.RGBA

	faceX, faceY     int // position of this card's image in image sheet
	backX, backY     int // position of this card's image in image sheet
	screenX, screenY int // current position on screen (after Fan)

	srcX, srcY float64 // smoothstep origin
	dstX, dstY float64 // smoothstep destination
	lerpStep   float64 // current lerp value 0.0 .. 1.0
	lerping    bool    // true if this card is smoothstepping

	dragging               bool // true if this card is being dragged
	dragStartX, dragStartY int  // starting point for dragging

	flipStep  float64 // if 0, we are not flipping
	flipWidth float64 // scale of the card width while flipping

	shake shakeState
}

// NewCard is a factory for Card objects
func NewCard(pack int, suit string, ordinal int) *Card {
	c := &Card{pack: pack, suit: suit, ordinal: ordinal}
	if c.suit == "Heart" || c.suit == "Diamond" {
		c.color = BasicColors["Red"]
	} else {
		c.color = BasicColors["Black"]
	}
	c.prone = true
	c.id = c.String()

	switch c.suit {
	case "Club":
		c.faceY = 0
	case "Diamond":
		c.faceY = 96
	case "Heart":
		c.faceY = 96 + 96
	case "Spade":
		c.faceY = 96 + 96 + 96
	}
	c.faceX = (c.ordinal - 1) * 71

	pt := backFrames[TheUserData.CardBack]
	c.backX, c.backY = pt.X, pt.Y

	return c
}

func (c *Card) String() string {
	return fmt.Sprintf("%d%c%02d", c.pack, c.suit[0], c.ordinal)
}

// ParseID decomposes a string id into Card members pack, suit, ordinal
func parseID(id string) (pack int, suit string, ordinal int) {
	var err error
	pack, err = strconv.Atoi(id[0:1])
	if err != nil || pack > 9 {
		log.Fatal("error in Card id" + id)
	}
	switch id[1:1] {
	case "C":
		suit = "Club"
	case "D":
		suit = "Diamond"
	case "H":
		suit = "Heart"
	case "S":
		suit = "Spade"
	default:
		log.Fatal("error in Card id" + id)
	}
	ordinal, err = strconv.Atoi(id[2:3]) // TODO beware leading 0
	if err != nil || ordinal < 1 || ordinal > 13 {
		log.Fatal("error in Card id" + id)
	}
	return // uses named return variables
}

// Position returns the x,y screen coords of this card
func (c *Card) Position() (int, int) {
	return c.screenX, c.screenY
}

// Rect gives the x,y screen coords of the card's top left and bottom right corners
func (c *Card) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = c.screenX
	y0 = c.screenY
	x1 = x0 + 71
	y1 = y0 + 96
	return // using named return parameters
}

// SetPosition sets the position of the Card
func (c *Card) SetPosition(x, y int) {
	c.screenX, c.screenY = x, y
}

// TransitionTo starts the transition of this Card
func (c *Card) TransitionTo(x, y int) {
	if x != c.screenX || y != c.screenY {
		c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
		c.dstX, c.dstY = float64(x), float64(y)
		c.lerpStep = 0
		c.lerping = true
	}
}

// IsBusy returns true of this card is lerping, flipping or being dragged
// func (c *Card) IsBusy() bool {
// 	return c.lerping || c.dragging || c.flipStep != 0
// }

// TransitionBackToPile starts the transition of this Card back to it's Pile TODO broken when fanned
// func (c *Card) TransitionBackToPile() {
// 	c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
// 	x, y := c.owner.Position()
// 	c.dstX, c.dstY = float64(x), float64(y)
// 	c.lerpStep = 0
// 	c.lerping = true
// }

// StartDrag informs card that it is being dragged
func (c *Card) StartDrag() {
	c.dragStartX, c.dragStartY = c.screenX, c.screenY
	c.dragging = true
	// println("start drag", c.id, "start", c.dragStartX, c.dragStartY)
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
	// println("stop drag", c.id)
	c.dragging = false
}

// CancelDrag informs card that it is no longer being dragged
func (c *Card) CancelDrag() {
	// println("cancel drag", c.id, "start", c.dragStartX, c.dragStartY, "screen", c.screenX, c.screenY)
	// c.TransitionTo(c.dragStartX, c.dragStartY)
	CTQ.Add(c, c.dragStartX, c.dragStartY)
	c.dragging = false
}

// Shake starts the transition of this Card left, right, center
func (c *Card) Shake() {
	if c.shake != notShaking {
		return
	}
	if c.lerping || c.dragging {
		return
	}
	// hijack the lerping src, dst positions
	c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
	c.dstX, c.dstY = float64(c.screenX-SHAKEAMOUNT), float64(c.screenY)
	c.shake = shakingLeft
}

// FlipUp flips the card face up
func (c *Card) FlipUp() {
	if c.prone && c.flipStep == 0.0 {
		c.prone = false        // card is immediately face up, else fan isn't correct
		c.flipStep = -FLIPSTEP // start by making card narrower
		c.flipWidth = 1.0
	}
}

// FlipDown flips the card face down
func (c *Card) FlipDown() {
	if !c.prone && c.flipStep == 0.0 {
		c.prone = true         // card is immediately face down, else fan isn't correct
		c.flipStep = -FLIPSTEP // start by making card narrower
		c.flipWidth = 1.0
	}
}

// Flip toggles the card
func (c *Card) Flip() {
	if c.prone {
		c.FlipUp()
	} else {
		c.FlipDown()
	}
}

// Animating returns true if this card is lerping, dragging or flipping
func (c *Card) Animating() bool {
	if c.lerping || c.dragging {
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
	if c.lerping {
		if c.lerpStep >= 1 {
			c.screenX, c.screenY = int(c.dstX), int(c.dstY)
			c.lerping = false
		} else {
			c.screenX = int(util.Smoothstep(c.srcX, c.dstX, c.lerpStep))
			c.screenY = int(util.Smoothstep(c.srcY, c.dstY, c.lerpStep))
			c.lerpStep += 0.05
		}
	}
	if c.flipStep != 0.0 {
		c.flipWidth += c.flipStep
		if c.flipWidth <= 0.15 {
			c.flipStep = FLIPSTEP // now make card wider
			// c.prone = !c.prone
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
				c.dstX = c.srcX + SHAKEAMOUNT
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
		if c.prone {
			// card is getting narrower, and it's going to show face down, but show face up
			img = faceImageSheet.SubImage(image.Rect(c.faceX, c.faceY, c.faceX+71, c.faceY+96)).(*ebiten.Image)
		} else {
			// card is getting narrower, and it's going to show face up, but show face down
			img = backImageSheet.SubImage(image.Rect(c.backX, c.backY, c.backX+71, c.backY+96)).(*ebiten.Image)
		}
	} else {
		if c.prone {
			img = backImageSheet.SubImage(image.Rect(c.backX, c.backY, c.backX+71, c.backY+96)).(*ebiten.Image)
		} else {
			img = faceImageSheet.SubImage(image.Rect(c.faceX, c.faceY, c.faceX+71, c.faceY+96)).(*ebiten.Image)
		}
	}

	if c.flipStep != 0 {
		// img = ebiten.NewImageFromImage(img)
		op.GeoM.Translate(float64(-71/2), 0)
		op.GeoM.Scale(c.flipWidth, 1.0)
		op.GeoM.Translate(float64(71/2), 0)
	}

	op.GeoM.Translate(float64(c.screenX), float64(c.screenY))

	if c.flipStep == 0 && (c.lerping == true || c.dragging == true) {
		op.GeoM.Translate(2, 2)
		screen.DrawImage(shadowImage, op)
		op.GeoM.Translate(-2, -2)
	}

	screen.DrawImage(img, op)
}
