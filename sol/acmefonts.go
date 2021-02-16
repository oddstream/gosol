// Copyright ©️ 2021 oddstream.games

package maze

import (
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// AcmeFonts contains references to small, normal, large and huge Acme fonts
type AcmeFonts struct {
	small  font.Face
	normal font.Face
	large  font.Face
	huge   font.Face
}

// NewAcmeFonts loads some fonts and returns a pointer to an object referencing them
func NewAcmeFonts() *AcmeFonts {

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
	tt, err := truetype.Parse(Acme_ttf)
	if err != nil {
		log.Fatal(err)
	}

	af := &AcmeFonts{}

	af.small = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.normal = truetype.NewFace(tt, &truetype.Options{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.large = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	af.huge = truetype.NewFace(tt, &truetype.Options{
		Size:    256,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return af
}
