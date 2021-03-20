package sol

import (
	"bytes"
	"image"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// func init() {
// 	println("Club", string(rune(9827)), "Diamond", string(rune(9830)), "Heart", string(rune(9829)), "Spade", string(rune(9824)))
// }

var (
	scalableFaceImages map[CardID]*ebiten.Image
	scalableBackImage  *ebiten.Image
	shadowImage        *ebiten.Image
)

func createFaceImage(ID CardID) *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(BasicColors["White"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Fill()

	dc.SetColor(BasicColors["Silver"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	dc.SetColor(ID.Color())
	dc.SetFontFace(schriftbank.CardOrdinal)
	dc.DrawStringAnchored(util.OrdinalToShortString(ID.Ordinal()), float64(CardWidth)/3.333, float64(CardHeight)/6.666, 0.5, 0.5)
	dc.Stroke()

	dc.SetFontFace(schriftbank.CardSymbolRegular)
	// https://unicodelookup.com/#club/1
	// https://www.fileformat.info/info/unicode/char/2665/index.htm
	// https://www.fileformat.info/info/unicode/char/2665/fontsupport.htm
	// https://github.com/fogleman/gg/blob/v1.3.0/context.go#L679
	var r rune
	switch ID.Suit() {
	case CLUB:
		r = 9827 //0x2663
	case DIAMOND:
		r = 9830 // 0x2666
	case HEART:
		r = 9829 //0x2665
	case SPADE:
		r = 9824 // 0x2660
	}
	// to make the symbols align with the ordinal short string, draw it down a little, hence /6 instead of /6.666
	dc.DrawStringAnchored(string(r), float64(CardWidth)-float64(CardWidth)/(3.333), float64(CardHeight)/6, 0.5, 0.5)
	dc.Stroke()

	if ID.Ordinal() == 1 && ID.Suit() == SPADE {
		img, _, err := image.Decode(bytes.NewReader(logoBytes))
		if err != nil {
			log.Fatal(err)
		}
		logoWidth := img.Bounds().Dx()
		logoHeight := img.Bounds().Dy()
		dcLogo := gg.NewContext(logoWidth, logoHeight)
		var scale float64 = float64(CardWidth) / float64(logoWidth)
		dcLogo.ScaleAbout(scale, scale, float64(logoWidth)/2, float64(logoHeight)/2)
		// dcLogo.RotateAbout(gg.Radians(-45), float64(logoWidth)/2, float64(logoHeight)/2)
		dcLogo.DrawImageAnchored(img, logoWidth/2, logoHeight/2, 0.5, 0.5)
		dc.DrawImageAnchored(dcLogo.Image(), CardWidth/2, CardHeight/2, 0.5, 0.4)
	} else if ID.Ordinal() == 1 || ID.Ordinal() > 10 {
		dc.SetFontFace(schriftbank.CardSymbolLarge)
		dc.DrawStringAnchored(string(r), float64(CardWidth)/2, float64(CardHeight)/1.75, 0.5, 0.5)
	}

	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func createBackImage() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(ExtendedColors[TheUserData.BackColor])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Fill()
	dc.SetColor(BasicColors["Silver"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func createShadowImage() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetRGBA(0.1, 0.1, 0.1, 0.9)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), float64(CardWidth)/12)
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// BuildScalables builds the card images that can change in scale, after CardWidth,Height have been set
func BuildScalables() {

	schriftbank.MakeCardFonts(CardWidth) // CardWidth/Height have now been set

	if TheUserData.CardStyle != "retro" {
		scalableFaceImages = make(map[CardID]*ebiten.Image)
		for ord := 1; ord < 14; ord++ {
			for suit := 1; suit < 5; suit++ {
				ID := NewCardID(0, suit, ord)
				scalableFaceImages[ID] = createFaceImage(ID)
			}
		}
		scalableBackImage = createBackImage()
	}

	shadowImage = createShadowImage()
}

// getScalableImages reloads the face and back image for this card
func (c *Card) getScalableImages() {
	subid := NewCardID(0, c.Suit(), c.Ordinal())
	c.faceImg = scalableFaceImages[subid]
	c.backImg = scalableBackImage
	// either faceImg or backImg may be nil if we are booting up
}
