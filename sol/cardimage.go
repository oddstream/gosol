package sol

import (
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/cardid"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

const (
	// Ordinals
	// 1.0 / 7.5 = 0.1333333
	// 1.0 - 0.1333333 = 0.866666
	COTLX = 0.1333 // Card Ordinal Top Left X
	COTLY = 0.1333 // Card Ordinal Top Left Y
	COTRX = 0.8666 // Card Ordinal Top Right X
	COTRY = 0.1333 // Card Ordinal Top Right Y
	COBLX = 0.1333 // Card Ordinal Bottom Left X
	COBLY = 0.8666 // Card Ordinal Bottom Left Y
	COBRX = 0.8666 // Card Ordinal Bottom Right X
	COBRY = 0.8666 // Card Ordinal Bottom Right Y
	// Pips
	XL = 0.375 // X Left
	XC = 0.5   // X Center
	XR = 0.625 // X Right
	YT = 0.166 // Y Top
	Y1 = 0.333 // Y between top and center
	Y2 = 0.4
	YC = 0.5 // Y Center
	Y3 = 0.6
	Y4 = 0.666 // Y between center and bottom
	YB = 0.833 // Y Bottom
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

func cardColor(cid cardid.CardID) color.RGBA {
	suit := cid.Suit()
	if TheSettings.ColorfulCards {
		switch TheBaize.script.CardColors() {
		case 4:
			switch suit {
			case cardid.NOSUIT:
				return BasicColors["Silver"]
			case cardid.CLUB:
				return ExtendedColors[TheSettings.ClubColor]
			case cardid.DIAMOND:
				return ExtendedColors[TheSettings.DiamondColor]
			case cardid.HEART:
				return ExtendedColors[TheSettings.HeartColor]
			case cardid.SPADE:
				return ExtendedColors[TheSettings.SpadeColor]
			}
		case 2:
			switch suit {
			case cardid.NOSUIT:
				return BasicColors["Silver"]
			case cardid.CLUB, cardid.SPADE:
				return ExtendedColors[TheSettings.BlackColor]
			case cardid.DIAMOND, cardid.HEART:
				return ExtendedColors[TheSettings.RedColor]
			}
		case 1:
			return ExtendedColors[TheSettings.SpadeColor]
		}
	} else {
		switch suit {
		case cardid.NOSUIT:
			return BasicColors["Silver"]
		case cardid.CLUB, cardid.SPADE:
			return ExtendedColors[TheSettings.BlackColor]
		case cardid.DIAMOND, cardid.HEART:
			return ExtendedColors[TheSettings.RedColor]
		}
	}
	return BasicColors["Purple"]
}

// createFaceImage tries to draw an image for this card that looks like kenney.nl playingCards.png
func createFaceImage(ID cardid.CardID) *ebiten.Image {
	w := float64(CardWidth)
	h := float64(CardHeight)

	dc := gg.NewContext(CardWidth, CardHeight)

	// draw the basic card face
	dc.SetColor(ExtendedColors[TheSettings.CardFaceColor])
	dc.DrawRoundedRectangle(0, 0, w, h, CardCornerRadius)
	dc.Fill()

	// surround with a thin border
	dc.SetLineWidth(1)
	// card face is probably light, so darken the border a bit
	dc.SetRGBA(0, 0, 0, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, w-2, h-2, CardCornerRadius)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	var cardOrdinal = ID.Ordinal()
	var suitRune rune = ID.SuitRune()
	var cardColor color.RGBA = cardColor(ID)
	// if ID.Joker() {
	// 	// if a joker is pretending to be a certain card, then show it's pretend ordinal and suit, but faded
	// 	cardColor.A = 64
	// }

	// draw the card ordinals in top left and bottom right corners
	dc.SetColor(cardColor)
	if cardOrdinal == 10 {
		dc.SetFontFace(schriftbank.CardOrdinalSmall)
	} else {
		dc.SetFontFace(schriftbank.CardOrdinal)
	}
	dc.DrawStringAnchored(util.OrdinalToShortString(cardOrdinal), w*COTLX, h*COTLY, 0.5, 0.4)
	dc.RotateAbout(gg.Radians(180), w*COBRX, h*COBRY)
	dc.DrawStringAnchored(util.OrdinalToShortString(cardOrdinal), w*COBRX, h*COBRY, 0.5, 0.4)
	dc.RotateAbout(gg.Radians(180), w*COBRX, h*COBRY)
	dc.Stroke()

	// https://unicodelookup.com/#club/1
	// https://www.fileformat.info/info/unicode/char/2665/index.htm

	// https://www.fileformat.info/info/unicode/char/2663/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2666/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2665/fontsupport.htm
	// https://www.fileformat.info/info/unicode/char/2660/fontsupport.htm

	// https://github.com/fogleman/gg/blob/v1.3.0/context.go#L679
	if suitRune != 0 {
		if cardOrdinal == 1 || cardOrdinal > 10 {
			// Ace, Jack, Queen, King get suit runes at top right and bottom left
			// so the suit can be seen when fanned
			// they also get purdy rectangles in the middle
			dc.SetFontFace(schriftbank.CardSymbolRegular)
			dc.SetColor(cardColor)
			dc.DrawStringAnchored(string(suitRune), w*COTRX, h*COTRY, 0.5, 0.4)
			dc.RotateAbout(gg.Radians(180), w*COBLX, h*COBRY)
			dc.DrawStringAnchored(string(suitRune), w*COBLX, h*COBLY, 0.5, 0.4)
			dc.RotateAbout(gg.Radians(180), w*COBLX, h*COBRY)
			// dc.DrawStringAnchored(string(r), w*0.16, h*0.27, 0.5, 0.5)
			// dc.DrawStringAnchored(string(r), w*0.84, h*0.73, 0.5, 0.5)
			dc.Stroke()

			dc.SetRGBA(0, 0, 0, 0.05)
			dc.DrawRectangle(w*0.25, h*0.25, w*0.5, h*0.5)
			dc.Fill()

			dc.SetColor(cardColor)
			dc.SetFontFace(schriftbank.CardSymbolLarge)
			dc.DrawStringAnchored(string(suitRune), w*0.5, h*0.44, 0.5, 0.5)

		} else if cardOrdinal > 0 {
			// a blank joker will have an ordinal of zero
			dc.SetColor(cardColor)
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
				dc.DrawStringAnchored(string(suitRune), w*pip.X, h*pip.Y, 0.5, 0.5)
			}
		}
		dc.Stroke()
	}
	return ebiten.NewImageFromImage(dc.Image())
}

/***
func createSimpleFaceImage(ID CardID) *ebiten.Image {
	w := float64(CardWidth)
	h := float64(CardHeight)

	dc := gg.NewContext(CardWidth, CardHeight)

	// draw the basic card face
	dc.SetColor(ExtendedColors[TheSettings.CardFaceColor])
	dc.DrawRoundedRectangle(0, 0, w, h, CardCornerRadius)
	dc.Fill()

	// surround with a thin border
	dc.SetLineWidth(1)
	// card face is probably light, so darken the border a bit
	dc.SetRGBA(0, 0, 0, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, w-2, h-2, CardCornerRadius)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	var cardOrdinal = ID.Ordinal()
	var suitRune rune = ID.SuitRune()
	var cardColor color.RGBA = ID.Color()

	dc.SetColor(cardColor)

	// draw the card ordinals in top left corner
	dc.SetFontFace(schriftbank.CardOrdinalSimple)
	dc.DrawStringAnchored(util.OrdinalToShortString(cardOrdinal), w*CTLX, h*CTLY, 0.35, 0.5)
	dc.Stroke()

	// draw the suit rune in top right corner
	if suitRune != 0 {
		dc.SetFontFace(schriftbank.CardSymbolSimple)
		dc.DrawStringAnchored(string(suitRune), w*CTRX, h*CTRY, 0.65, 0.5)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}
***/

func CreateCardBackImage(color string) *ebiten.Image {
	w := float64(CardWidth)
	h := float64(CardHeight)

	dc := gg.NewContext(CardWidth, CardHeight)

	dc.SetColor(ExtendedColors[color])
	dc.DrawRoundedRectangle(0, 0, w, h, CardCornerRadius)
	dc.Fill()

	dc.SetLineWidth(1)
	// card back probably dark, so lighten the border a bit
	dc.SetRGBA(1, 1, 1, 0.1)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, w-2, h-2, CardCornerRadius)
	dc.Stroke() // otherwise outline gets drawn in textColor (!?)

	dc.SetFontFace(schriftbank.CardSymbolRegular)
	dc.SetRGBA(0, 0, 0, 0.2)
	dc.DrawStringAnchored(string(cardid.SPADE_RUNE), w*0.4, h*0.4, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.DrawStringAnchored(string(cardid.HEART_RUNE), w*0.6, h*0.4, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.1)
	dc.DrawStringAnchored(string(cardid.DIAMOND_RUNE), w*0.4, h*0.6, 0.5, 0.5)
	dc.SetRGBA(0, 0, 0, 0.2)
	dc.DrawStringAnchored(string(cardid.CLUB_RUNE), w*0.6, h*0.6, 0.5, 0.5)
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

func CreateCardImages() {
	// defer util.Duration(time.Now(), "CreateCardImages")
	if CardWidth == 0 || CardHeight == 0 {
		log.Println("CreateCardImages called with zero card dimensions") // seen to happen in WASM
		return
	}
	schriftbank.MakeCardFonts(CardWidth)
	for _, suit := range []int{cardid.NOSUIT, cardid.CLUB, cardid.DIAMOND, cardid.HEART, cardid.SPADE} {
		for ord := 1; ord < 14; ord++ {
			ID := cardid.NewCardID(0, suit, ord)
			TheCardFaceImageLibrary[(suit*13)+(ord-1)] = createFaceImage(ID)
		}
	}
	CardBackImage = CreateCardBackImage(TheSettings.CardBackColor)
	MovableCardBackImage = CreateCardBackImage(TheSettings.MovableCardBackColor)
	CardShadowImage = CreateCardShadowImage()
}
