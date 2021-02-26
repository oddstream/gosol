package sol

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/util"
)

func init() {
	println("Heart", string(rune(9829)), "Diamond", string(rune(9830)), "Spade", string(rune(9824)), "Club", string(rune(9827)))
}

func createFaceImage(suit string, ord int, textColor *color.RGBA) *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(BasicColors["White"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Fill()

	dc.SetColor(BasicColors["Silver"])
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), 6)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	dc.SetColor(textColor)
	dc.SetFontFace(TheCardFonts.regular)

	if ord == 10 {
		dc.DrawString("X", float64(CardWidth)/8, float64(CardHeight)/3.5)
		// dc.DrawString(util.OrdinalToChar(ord), float64(CardWidth)/12, float64(CardHeight)/3.5)
	} else {
		dc.DrawString(util.OrdinalToChar(ord), float64(CardWidth)/8, float64(CardHeight)/3.5)
	}
	// https://unicodelookup.com/#club/1
	// https://www.fileformat.info/info/unicode/char/2665/index.htm
	// https://www.fileformat.info/info/unicode/char/2665/fontsupport.htm
	// https://github.com/fogleman/gg/blob/v1.3.0/context.go#L679
	var r rune
	switch suit {
	case "H":
		r = 9829 //0x2665
	case "D":
		r = 9830 // 0x2666
	case "S":
		r = 9824 // 0x2660
	case "C":
		r = 9827 //0x2663
	}
	dc.DrawString(string(r), float64(CardWidth)/1.75, float64(CardHeight)/3.5)

	dc.SetFontFace(TheCardFonts.large)
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
