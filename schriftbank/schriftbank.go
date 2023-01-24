// Package schriftbank provides a collection for fonts for package sol
package schriftbank

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/Acme-Regular.ttf
var acmeFontBytes []byte

//go:embed assets/DejaVuSansCondensed-Bold.ttf
var symbolFontBytes []byte

//go:embed assets/Roboto-Regular.ttf
var robotoRegularFontBytes []byte

//go:embed assets/Roboto-Medium.ttf
var robotoMediumFontBytes []byte

var (
	// RobotoRegular14 used by UI toast
	RobotoRegular14 font.Face
	// RobotoMedium24 used by UI
	RobotoMedium24 font.Face
	// CardSymbolSmall is used to draw the suit symbol under the card ordinal
	CardSymbolSmall font.Face
	// CardSymbolRegular is used to draw the suit symbols as pips in the middle of the card
	CardSymbolRegular font.Face
	// CardSymbolSimple is used in Simple card faces
	// CardSymbolSimple font.Face
	// CardSymbolLarge is used to draw the large suit symbol
	CardSymbolLarge font.Face
	// CardSymbolHuge is used to draw the recycle symbol
	CardSymbolHuge font.Face
	// CardOrdinalSmall is used to draw the card ordinal (A to K)
	CardOrdinalSmall font.Face
	// CardOrdinal is used to draw the card ordinal (A to K)
	CardOrdinal font.Face
	// CardOrdinalSimple is used in Simple card faces
	// CardOrdinalSimple font.Face
	// CardOrdinalLarge is used to draw the card ordinal (J, Q, K)
	CardOrdinalLarge font.Face
	CardOrdinalHuge  font.Face
)

func init() {

	// defer util.Duration(time.Now(), "init schriftbank")

	tt, err := truetype.Parse(robotoRegularFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	RobotoRegular14 = truetype.NewFace(tt, &truetype.Options{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	tt, err = truetype.Parse(robotoMediumFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	RobotoMedium24 = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// MakeCardFonts creates the fonts used for Card, once size of card is known (or has changed)
func MakeCardFonts(cardWidth int) {
	// defer util.Duration(time.Now(), "MakeCardFonts")

	tt, err := truetype.Parse(acmeFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	CardOrdinalSmall = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.233,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	CardOrdinal = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.275,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	// CardOrdinalSimple = truetype.NewFace(tt, &truetype.Options{
	// 	Size:    float64(cardWidth) * 0.4,
	// 	DPI:     72,
	// 	Hinting: font.HintingFull,
	// })
	CardOrdinalLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.75,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	tt, err = truetype.Parse(symbolFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	CardSymbolSmall = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.2,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	CardSymbolRegular = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.275,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	// CardSymbolSimple = truetype.NewFace(tt, &truetype.Options{
	// 	Size:    float64(cardWidth) * 0.4,
	// 	DPI:     72,
	// 	Hinting: font.HintingFull,
	// })
	CardSymbolLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) * 0.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	CardSymbolHuge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth),
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
