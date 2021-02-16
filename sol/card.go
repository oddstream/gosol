package sol

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

	var err error
	faceImageSheet, _, err = ebitenutil.NewImageFromFile("sol/cards71x96.png")
	if err != nil {
		log.Fatal("cannot load sol/cards71x96.png")
	}
	backImageSheet, _, err = ebitenutil.NewImageFromFile("sol/windows_16bit_cards.png")
	if err != nil {
		log.Fatal("cannot load sol/windows_16bit_cards.png")
	}
}

// Card object
type Card struct {
	pack    int
	suit    string
	ordinal int
	prone   bool
	id      string
	color   color.RGBA

	screenX, screenY float64
}

// NewCard is the factory for Card objects
func NewCard(pack int, suit string, ordinal int) *Card {
	c := &Card{pack: pack, suit: suit, ordinal: ordinal}
	if c.suit == "Heart" || c.suit == "Diamond" {
		c.color = BasicColors["Red"]
	} else {
		c.color = BasicColors["Black"]
	}
	c.prone = true
	c.id = c.String()
	return c
}

func (c *Card) String() string {
	return fmt.Sprintf("%d%c%02d", c.pack, c.suit[0], c.ordinal)
}

// ParseID decomposes a string id into Card members pack, suit, ordinal
func (c *Card) ParseID(id string) {

}

// Layout implements ebiten.Game's Layout.
func (c *Card) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update the baize state (transitions, user input)
func (c *Card) Update() error {
	return nil
}

// Draw renders the baize into the screen
func (c *Card) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.screenX, c.screenY)
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
