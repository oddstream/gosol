package sol

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// func init() {
// 	println("Club", string(rune(9827)), "Diamond", string(rune(9830)), "Heart", string(rune(9829)), "Spade", string(rune(9824)))
// }

type ModernCardImageProvider struct {
	CardImages
}

func cardCornerRadius() float64 {
	return float64(CardWidth) / 15
}

func createModernFaceImage(ID CardID) *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), cardCornerRadius())
	dc.Fill()

	dc.SetLineWidth(2)
	dc.SetRGBA(0, 0, 0, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), cardCornerRadius())
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

	// the following increases duration from 50 to 60ms
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
		// dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

func createModernBackImage(width, height int, backColor color.Color) *ebiten.Image {
	dc := gg.NewContext(width, height)
	dc.SetColor(backColor)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cardCornerRadius())
	dc.Fill()
	dc.SetLineWidth(2)
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.DrawRoundedRectangle(1, 1, float64(width-2), float64(height-2), cardCornerRadius())
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func createModernShadowImage(width, height int) *ebiten.Image {
	dc := gg.NewContext(width, height)
	dc.SetRGBA(0, 0, 0, 0.5)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cardCornerRadius())
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// func createModernMovableImage(width, height int) *ebiten.Image {
// 	dc := gg.NewContext(width, height)
// 	dc.SetColor(ExtendedColors["Gold"])
// 	dc.SetLineWidth(2)
// 	dc.DrawRoundedRectangle(1, 1, float64(width-2), float64(height-2), cardCornerRadius())
// 	dc.Stroke()
// 	return ebiten.NewImageFromImage(dc.Image())
// }

func NewModernCardImageProvider() *ModernCardImageProvider {
	ip := &ModernCardImageProvider{}
	ip.faceImgs = make(map[CardID]*ebiten.Image)
	for ord := 1; ord < 14; ord++ {
		for _, suit := range []int{CLUB, DIAMOND, HEART, SPADE} {
			ID := NewCardID(0, suit, ord)
			ip.faceImgs[ID] = createModernFaceImage(ID)
		}
	}
	ip.backImgs = make(map[string]*ebiten.Image)
	for k, v := range ExtendedColors {
		ip.backImgs[k] = createModernBackImage(CardWidth, CardHeight, v)
	}
	ip.shadowImg = createModernShadowImage(CardWidth, CardHeight)
	// ip.movableImg = createModernMovableImage(CardWidth, CardHeight)
	return ip
}

func (ip *ModernCardImageProvider) FaceImage(ID CardID) *ebiten.Image {
	ID = ID & CardID(suitMask|ordinalMask)
	img, ok := ip.faceImgs[ID]
	if !ok {
		log.Panic("missing scalable face image")
	}
	return img
}

func (ip *ModernCardImageProvider) BackImage(colorName string) *ebiten.Image {
	return ip.backImgs[colorName]
}

func (ip *ModernCardImageProvider) BackImages() map[string]*ebiten.Image {
	return ip.backImgs
}

func (ip *ModernCardImageProvider) ShadowImage() *ebiten.Image {
	return ip.shadowImg
}

// func (ip *ModernCardImageProvider) MovableImage() *ebiten.Image {
// 	return ip.movableImg
// }
