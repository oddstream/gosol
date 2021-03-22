package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"time"

	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
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
	retroFaceImages map[CardID]*ebiten.Image
	retroBackImages map[string]*ebiten.Image
)

func init() {
	// https://ebiten.org/examples/tiles.html
	// uses
	// screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	// i.e. draws the tile direct from a subimage of imagesheet

	// load fixed-size images for retro cards from spritesheets

	println("loading retro card spritesheets")
	defer util.Duration(time.Now(), "retro init")
	img, _, err := image.Decode(bytes.NewReader(faceBytes))
	if err != nil {
		log.Panic(err)
	}
	faceImageSheet = ebiten.NewImageFromImage(img)
	faceBytes = nil

	retroFaceImages = make(map[CardID]*ebiten.Image)
	for ord := 1; ord < 14; ord++ {
		for suit := 1; suit < 5; suit++ {
			ID := NewCardID(0, suit, ord)
			retroFaceImages[ID] = createRetroFaceImage(ID)
		}
	}

	img, _, err = image.Decode(bytes.NewReader(backBytes))
	if err != nil {
		log.Panic(err)
	}
	backImageSheet = ebiten.NewImageFromImage(img)
	backBytes = nil

	retroBackImages = make(map[string]*ebiten.Image)
	for name, pt := range backFrames {
		backX, backY := pt.X, pt.Y
		backImg := backImageSheet.SubImage(image.Rect(backX, backY, backX+71, backY+96)).(*ebiten.Image)
		if backImg == nil {
			log.Panic("no back image")
		}
		retroBackImages[name] = backImg
	}
}

func createRetroFaceImage(ID CardID) *ebiten.Image {
	var faceX, faceY int
	switch ID.Suit() {
	case CLUB:
		faceY = 0
	case DIAMOND:
		faceY = 96
	case HEART:
		faceY = 96 + 96
	case SPADE:
		faceY = 96 + 96 + 96
	}
	faceX = (ID.Ordinal() - 1) * 71

	faceImg := faceImageSheet.SubImage(image.Rect(faceX, faceY, faceX+71, faceY+96)).(*ebiten.Image)
	if faceImg == nil {
		log.Panic("no face image")
	}
	return faceImg
}

// getRetroImages reloads the face and back image for this card
func (c *Card) getRetroImages() {
	// var faceX, faceY, backX, backY int
	// switch c.StringSuit() {
	// case "Club":
	// 	faceY = 0
	// case "Diamond":
	// 	faceY = 96
	// case "Heart":
	// 	faceY = 96 + 96
	// case "Spade":
	// 	faceY = 96 + 96 + 96
	// }
	// faceX = (c.Ordinal() - 1) * 71

	// c.faceImg = faceImageSheet.SubImage(image.Rect(faceX, faceY, faceX+71, faceY+96)).(*ebiten.Image)
	// if c.faceImg == nil {
	// 	log.Panic("no face image")
	// }
	// // dc := gg.NewContextForImage(ebiten.NewImageFromImage(c.faceImg))
	// // dc.SetColor(BasicColors["Black"])
	// // // draw the RoundedRect entirely INSIDE the context
	// // dc.DrawRoundedRectangle(1, 1, float64(dc.Width()-2), float64(dc.Height()-2), 2)
	// // dc.Stroke() // otherwise outline gets drawn in textColor (!?)
	// // c.faceImg = ebiten.NewImageFromImage(dc.Image())

	// pt := backFrames[TheUserData.CardBack]
	// backX, backY = pt.X, pt.Y
	// c.backImg = backImageSheet.SubImage(image.Rect(backX, backY, backX+71, backY+96)).(*ebiten.Image)
	// if c.backImg == nil {
	// 	log.Panic("no back image")
	// }
	var ID CardID = c.ID & CardID(suitMask|ordinalMask)
	c.faceImg = retroFaceImages[ID]
	c.backImg = retroBackImages[TheUserData.CardBack]
}

func (b *Baize) ShowCardBackPicker() {
	b.ui.ShowCardBackPicker()
}
