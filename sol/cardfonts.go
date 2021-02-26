// Copyright ©️ 2021 oddstream.games

package sol

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/DejaVuSans-Bold.ttf
var cardFontBytes []byte

// CardFonts contains references to regular and large fonts used on cards and piles
type CardFonts struct {
	regular font.Face
	large   font.Face
}

// NewCardFonts loads some fonts and returns a pointer to an object referencing them
func NewCardFonts() *CardFonts {

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
	tt, err := truetype.Parse(cardFontBytes)
	if err != nil {
		log.Fatal(err)
	}

	cf := &CardFonts{}

	cf.regular = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth) / 2,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	cf.large = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(CardWidth),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return cf
}
