// Package schriftbank provides a collection for fonts for package sol
package schriftbank

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"oddstream.games/gosol/util"
)

//go:embed assets/Acme-Regular.ttf
var acmeFontBytes []byte

//go:embed assets/DejaVuSans-Bold.ttf
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
	// Symbol24 used by UI
	Symbol24 font.Face
	// CardSymbolRegular is used to draw the suit symbol
	CardSymbolRegular font.Face
	// CardSymbolLarge is used to draw the large suit symbol
	CardSymbolLarge font.Face
	// CardOrdinal is used to draw the card ordinal (A to K)
	CardOrdinal font.Face
)

func init() {

	println("loading fonts")
	defer util.Duration(time.Now(), "init schriftbank")

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

	tt, err = truetype.Parse(symbolFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	Symbol24 = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// MakeCardFonts creates the fonts used for Card, once size of card is known (or has changed)
func MakeCardFonts(cardWidth int) {
	defer util.Duration(time.Now(), "MakeCardFonts")
	tt, err := truetype.Parse(acmeFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	CardOrdinal = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) / 2.25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	tt, err = truetype.Parse(symbolFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	CardSymbolRegular = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) / 2.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	CardSymbolLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(cardWidth) / 1.25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

}
