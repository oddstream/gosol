package sol

import (
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/schriftbank"
	"oddstream.games/gomps5/util"
)

const (
	CTLX = 0.1333
	CTLY = 0.1333
	CTRX = 0.8666
	CTRY = 0.1333
	CBLX = 0.1333
	CBLY = 0.8666
	CBRX = 0.8666
	CBRY = 0.8666
	XL   = 0.375 // X Left
	XC   = 0.5   // X Center
	XR   = 0.625 // X Right
	YT   = 0.166 // Y Top
	Y1   = 0.33  // Y between top and center
	Y2   = 0.4
	YC   = 0.5 // Y Center
	Y3   = 0.6
	Y4   = 0.66  // Y between center and bottom
	YB   = 0.833 // Y Bottom
)

type PipInfo struct {
	X, Y float64
	SZ   int
}

// array of 13 slices of gg.Point
var Pips [13][]PipInfo = [13][]PipInfo{
	/* 1 */ {},
	/* 2 */ {
		{X: XC, Y: YT},
		{X: XC, Y: YB},
	},
	/* 3 */ {
		{X: XC, Y: YT},
		{X: XC, Y: YC},
		{X: XC, Y: YB},
	},
	/* 4 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 5 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XC, Y: YC},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 6 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XL, Y: YC}, {X: XR, Y: YC},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 7 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XC, Y: Y1},
		{X: XL, Y: YC}, {X: XR, Y: YC},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 8 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XC, Y: Y1},
		{X: XL, Y: YC}, {X: XR, Y: YC},
		{X: XC, Y: Y4},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 9 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XL, Y: Y2}, {X: XR, Y: Y2},
		{X: XC, Y: YC, SZ: -1}, // smaller
		{X: XL, Y: Y3}, {X: XR, Y: Y3},
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	/* 10 */ {
		{X: XL, Y: YT}, {X: XR, Y: YT},
		{X: XC, Y: Y1 - 0.03, SZ: -1}, // smaller
		{X: XL, Y: Y2}, {X: XR, Y: Y2},
		{X: XL, Y: Y3}, {X: XR, Y: Y3},
		{X: XC, Y: Y4 + 0.05, SZ: -1}, // smaller
		{X: XL, Y: YB}, {X: XR, Y: YB},
	},
	{},
	{},
	{},
}

// createFaceImage tries to draw an image for this card that looks like kenney.nl playingCards.png
func createFaceImage(ID CardID) *ebiten.Image {
	w := float64(CardWidth)
	h := float64(CardHeight)

	dc := gg.NewContext(CardWidth, CardHeight)

	// draw the basic card face
	dc.SetColor(ExtendedColors[ThePreferences.CardFaceColor])
	dc.DrawRoundedRectangle(0, 0, w, h, CardCornerRadius)
	dc.Fill()

	// surround with a border
	dc.SetLineWidth(2)
	// card face is probably light, so darken the border a bit
	dc.SetRGBA(0, 0, 0, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, w-2, h-2, CardCornerRadius)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	// draw the card ordinals in top left and bottom right corners
	dc.SetColor(ID.Color())
	if ID.Ordinal() == 10 {
		dc.SetFontFace(schriftbank.CardOrdinalSmall)
	} else {
		dc.SetFontFace(schriftbank.CardOrdinal)
	}
	dc.DrawStringAnchored(util.OrdinalToShortString(ID.Ordinal()), w*CTLX, h*CTLY, 0.5, 0.4)
	dc.RotateAbout(gg.Radians(180), w*CBRX, h*CBRY)
	dc.DrawStringAnchored(util.OrdinalToShortString(ID.Ordinal()), w*CBRX, h*CBRY, 0.5, 0.4)
	dc.RotateAbout(gg.Radians(180), w*CBRX, h*CBRY)
	dc.Stroke()

	// https://unicodelookup.com/#club/1
	// https://www.fileformat.info/info/unicode/char/2665/index.htm

	// https://www.fileformat.info/info/unicode/char/2663/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2666/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2665/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2660/fontsupport.htm

	// https://github.com/fogleman/gg/blob/v1.3.0/context.go#L679
	var r rune = ID.SuitRune()
	if r != 0 {
		if ID.Ordinal() == 1 || ID.Ordinal() > 10 {
			// Ace, Jack, Queen, King get suit runes at top right and bottom left
			// so the suit can be seen when fanned
			// they also get purdy rectangles in the middle
			dc.SetFontFace(schriftbank.CardSymbolRegular)
			dc.SetColor(ID.Color())
			dc.DrawStringAnchored(string(r), w*CTRX, h*CTRY, 0.5, 0.4)
			dc.RotateAbout(gg.Radians(180), w*CBLX, h*CBRY)
			dc.DrawStringAnchored(string(r), w*CBLX, h*CBLY, 0.5, 0.4)
			dc.RotateAbout(gg.Radians(180), w*CBLX, h*CBRY)
			// dc.DrawStringAnchored(string(r), w*0.16, h*0.27, 0.5, 0.5)
			// dc.DrawStringAnchored(string(r), w*0.84, h*0.73, 0.5, 0.5)
			dc.Stroke()

			dc.SetRGBA(0, 0, 0, 0.05)
			dc.DrawRectangle(w*0.25, h*0.25, w*0.5, h*0.5)
			dc.Fill()

			dc.SetColor(ID.Color())
			dc.SetFontFace(schriftbank.CardSymbolLarge)
			dc.DrawStringAnchored(string(r), w*0.5, h*0.44, 0.5, 0.5)

		} else {

			dc.SetColor(ID.Color())
			// TODO would really like to draw some crown-ish symbols here
			// the chess glyphs only have king and queen, and would look a bit off
			// so using J Q K will have to do for now
			var pips = Pips[ID.Ordinal()-1]
			for _, pip := range pips {
				switch pip.SZ {
				case -1:
					dc.SetFontFace(schriftbank.CardSymbolSmall)
				case 0:
					dc.SetFontFace(schriftbank.CardSymbolRegular)
				case 1:
					dc.SetFontFace(schriftbank.CardSymbolLarge)
				}
				dc.DrawStringAnchored(string(r), w*pip.X, h*pip.Y, 0.5, 0.5)
			}
		}
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

func CreateCardBackImage() *ebiten.Image {
	w := float64(CardWidth)
	h := float64(CardHeight)

	dc := gg.NewContext(CardWidth, CardHeight)

	dc.SetColor(ExtendedColors[ThePreferences.CardBackColor])
	dc.DrawRoundedRectangle(0, 0, w, h, CardCornerRadius)
	dc.Fill()

	dc.SetLineWidth(2)
	// card back probably dark, so lighten the border a bit
	dc.SetRGBA(1, 1, 1, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, w-2, h-2, CardCornerRadius)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	dc.SetFontFace(schriftbank.CardSymbolRegular)
	dc.SetRGBA(0, 0, 0, 0.2)
	dc.DrawStringAnchored(string(SPADE_RUNE), w*0.4, h*0.4, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.DrawStringAnchored(string(HEART_RUNE), w*0.6, h*0.4, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.DrawStringAnchored(string(DIAMOND_RUNE), w*0.4, h*0.6, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.2)
	dc.DrawStringAnchored(string(CLUB_RUNE), w*0.6, h*0.6, 0.5, 0.5)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

func CreateCardShadowImage() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetRGBA(0, 0, 0, 0.5)
	// dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(0, 0, float64(CardWidth), float64(CardHeight), CardCornerRadius)
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

func CreateCardFaceImageLibrary() {
	defer util.Duration(time.Now(), "CreateCardFaceImageLibrary")

	for _, suit := range []int{NOSUIT, CLUB, DIAMOND, HEART, SPADE} {
		for ord := 1; ord < 14; ord++ {
			ID := NewCardID(0, suit, ord)
			TheCardFaceImageLibrary[(suit*13)+(ord-1)] = createFaceImage(ID)
		}
	}
}

func CreateCardImages() {
	if CardWidth == 0 || CardHeight == 0 {
		println("CreateCardImages called with zero card dimensions") // seen to happen in WASM
		return
	}
	// TODO MAYBE turn off drawing globally while this runs
	schriftbank.MakeCardFonts(CardWidth)
	CreateCardFaceImageLibrary()
	CardBackImage = CreateCardBackImage()
	CardShadowImage = CreateCardShadowImage()
}
