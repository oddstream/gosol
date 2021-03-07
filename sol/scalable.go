package sol

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

// func init() {
// 	println("Club", string(rune(9827)), "Diamond", string(rune(9830)), "Heart", string(rune(9829)), "Spade", string(rune(9824)))
// }

func createFaceImage(ID uint32) *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(BasicColors["White"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Fill()

	dc.SetColor(BasicColors["Silver"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	dc.SetColor(colorFromCardID(ID))
	dc.SetFontFace(TheCardFonts.acmeRegular)
	dc.DrawStringAnchored(util.OrdinalToShortString(ordinalFromCardID(ID)), float64(CardWidth)/3.333, float64(CardHeight)/6.666, 0.5, 0.5)
	dc.Stroke()

	dc.SetFontFace(TheCardFonts.symbolRegular)
	// https://unicodelookup.com/#club/1
	// https://www.fileformat.info/info/unicode/char/2665/index.htm
	// https://www.fileformat.info/info/unicode/char/2665/fontsupport.htm
	// https://github.com/fogleman/gg/blob/v1.3.0/context.go#L679
	var r rune
	switch suitFromCardID(ID) {
	case 1: // Club
		r = 9827 //0x2663
	case 2: // Diamond
		r = 9830 // 0x2666
	case 3: // Heart
		r = 9829 //0x2665
	case 4: // Spade
		r = 9824 // 0x2660
	}
	// to make the symbols align with the ordinal short string, draw it down a little, hence /6 instead of /6.666
	dc.DrawStringAnchored(string(r), float64(CardWidth)-float64(CardWidth)/(3.333), float64(CardHeight)/6, 0.5, 0.5)
	dc.Stroke()

	dc.SetFontFace(TheCardFonts.symbolLarge)
	dc.DrawStringAnchored(string(r), float64(CardWidth)/2, float64(CardHeight)/1.75, 0.5, 0.5)

	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func createBackImage() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(ExtendedColors[TheUserData.BackColor])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Fill()
	dc.SetColor(BasicColors["Silver"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
