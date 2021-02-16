package sol

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/*
   {x = 5, y = 4, width = 71, height = 96},    -- Aquarium
   {x = 85, y = 4, width = 71, height = 96},    -- CardHand
   {x = 165, y = 4, width = 71, height = 96},    -- Castle
   {x = 245, y = 4, width = 71, height = 96},    -- Empty
   {x = 325, y = 4, width = 71, height = 96},    -- Fishes
   {x = 405, y = 4, width = 71, height = 96},    -- FlowerBlack
   {x = 485, y = 4, width = 71, height = 96},    -- FlowerBlue
   {x = 5, y = 140, width = 71, height = 96},    -- PalmBeach
   {x = 85, y = 140, width = 71, height = 96},    -- Pattern1
   {x = 165, y = 140, width = 71, height = 96},    -- Pattern2
   {x = 245, y = 140, width = 71, height = 96},    -- Robot
   {x = 325, y = 140, width = 71, height = 96},    -- Roses
   {x = 405, y = 140, width = 71, height = 96},    -- Shell
*/

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

// Draw renders the baize into the screen
func (c *Card) Draw(screen *ebiten.Image) {
	if c.prone {
		pt := backFrames["JazzCup"]
		x, y := pt.X, pt.Y
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.screenX, c.screenY)
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
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.screenX, c.screenY)
		screen.DrawImage(faceImageSheet.SubImage(image.Rect(x, y, x+71, y+96)).(*ebiten.Image), op)
	}
}
