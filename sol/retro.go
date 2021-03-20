package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/cards71x96.png
var faceBytes []byte

//go:embed assets/windows_16bit_cards.png
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

	// load fixed-size images for retro cards from spritesheets

	println("loading retro card spritesheets")

	img, _, err := image.Decode(bytes.NewReader(faceBytes))
	if err != nil {
		log.Panic(err)
	}
	faceImageSheet = ebiten.NewImageFromImage(img)
	faceBytes = nil

	img, _, err = image.Decode(bytes.NewReader(backBytes))
	if err != nil {
		log.Panic(err)
	}
	backImageSheet = ebiten.NewImageFromImage(img)
	backBytes = nil
}

// getRetroImages reloads the face and back image for this card
func (c *Card) getRetroImages() {
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
		log.Panic("no face image")
	}
	pt := backFrames[TheUserData.CardBack]
	backX, backY = pt.X, pt.Y
	c.backImg = backImageSheet.SubImage(image.Rect(backX, backY, backX+CardWidth, backY+CardHeight)).(*ebiten.Image)
	if c.backImg == nil {
		log.Panic("no back image")
	}
}
