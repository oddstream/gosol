// Copyright ©️ 2021 oddstream.games

package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/Acme-Regular.ttf
var acmeFontBytes []byte

//go:embed assets/DejaVuSans-Bold.ttf
var symbolFontBytes []byte

// CardFonts contains references to regular and large fonts used on cards and piles
type CardFonts struct {
	acmeRegular   font.Face
	acmeLarge     font.Face
	symbolRegular font.Face
	symbolLarge   font.Face
}

// NewCardFonts loads some fonts and returns a pointer to an object referencing them
func NewCardFonts() *CardFonts {

	println("NewCardFonts with Width", CardWidth)

	cf := &CardFonts{}

	// path, err := filepath.Abs("/home/gilbert/Tetra/assets/Acme-Regular.ttf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// golang gotcha use full path so "go test" can find asset
	// bytes, err := ioutil.ReadFile("/home/gilbert/Tetra/assets/Acme-Regular.ttf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// https://pkg.go.dev/golang.org/x/image@v0.0.0-20201208152932-35266b937fa6/font

	tt, err := truetype.Parse(symbolFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	cf.symbolRegular = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth) / 2.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	cf.symbolLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth) / 1.25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	tt, err = truetype.Parse(acmeFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	cf.acmeRegular = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth) / 2.25,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	cf.acmeLarge = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return cf
}
