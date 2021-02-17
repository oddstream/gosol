package sol

import (
	// go:embed only allowed in Go files that import "embed"
	_ "embed"

	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

// golang gotcha: go:embed cannot apply to var inside func

//go:embed cards71x96.png
var faceBytes []byte

//go:embed windows_16bit_cards.png
var backBytes []byte

var (
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
	/*
		var err error
		faceImageSheet, _, err = ebitenutil.NewImageFromFile("sol/cards71x96.png")
		if err != nil {
			log.Fatal("cannot load sol/cards71x96.png")
		}
		backImageSheet, _, err = ebitenutil.NewImageFromFile("sol/windows_16bit_cards.png")
		if err != nil {
			log.Fatal("cannot load sol/windows_16bit_cards.png")
		}
	*/

	img, _, err := image.Decode(bytes.NewReader(faceBytes))
	if err != nil {
		log.Fatal(err)
	}
	faceImageSheet = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(backBytes))
	if err != nil {
		log.Fatal(err)
	}
	backImageSheet = ebiten.NewImageFromImage(img)
}

// Card object
type Card struct {
	pack    int
	suit    string
	ordinal int
	prone   bool
	id      string
	color   color.RGBA
	owner   CardOwner

	screenX, screenY int     // current position on screen
	srcX, srcY       float64 // smoothstep origin
	dstX, dstY       float64 // smoothstep destination
	lerpStep         float64 // current lerp value 0.0 .. 1.0
	lerping          bool    // true if this card is smoothstepping
}

// NewCard is the factory for Card objects
func NewCard(pack int, suit string, ordinal int) *Card {
	c := &Card{pack: pack, suit: suit, ordinal: ordinal}
	if c.suit == "Heart" || c.suit == "Diamond" {
		c.color = BasicColors["Red"]
	} else {
		c.color = BasicColors["Black"]
	}
	c.prone = false
	c.id = c.String()
	return c
}

func (c *Card) String() string {
	return fmt.Sprintf("%d%c%02d", c.pack, c.suit[0], c.ordinal)
}

// ParseID decomposes a string id into Card members pack, suit, ordinal
func (c *Card) ParseID(id string) {

}

// Rect gives the x,y screen coords of the card's top left and bottom right corners
func (c *Card) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = int(c.screenX)
	y0 = int(c.screenY)
	x1 = x0 + 71
	y1 = y0 + 96
	return // using named return parameters
}

// PositionTo sets the position of the Card
func (c *Card) PositionTo(x, y int) {
	c.screenX, c.screenY = x, y
}

// TransitionTo starts the transition of this Card
func (c *Card) TransitionTo(x, y int) {
	c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
	c.dstX, c.dstY = float64(x), float64(y)
	c.lerpStep = 0
	c.lerping = true
}

// TransitionBackToPile starts the transition of this Card back to it's Pile
func (c *Card) TransitionBackToPile() {
	c.srcX, c.srcY = float64(c.screenX), float64(c.screenY)
	x, y := c.owner.Position()
	c.dstX, c.dstY = float64(x), float64(y)
	c.lerpStep = 0
	c.lerping = true
}

// Layout implements ebiten.Game's Layout.
func (c *Card) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
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
	return nil
}

// Draw renders the card into the screen
func (c *Card) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.screenX), float64(c.screenY))
	if c.prone {
		pt := backFrames["JazzCup"]
		x, y := pt.X, pt.Y
		screen.DrawImage(backImageSheet.SubImage(image.Rect(x, y, x+71, y+96)).(*ebiten.Image), op)
	} else {
		var x, y int
		switch c.suit {
		case "Club":
			y = 0
		case "Diamond":
			y = 96
		case "Heart":
			y = 96 + 96
		case "Spade":
			y = 96 + 96 + 96
		}
		x = (c.ordinal - 1) * 71
		screen.DrawImage(faceImageSheet.SubImage(image.Rect(x, y, x+71, y+96)).(*ebiten.Image), op)
	}
}
