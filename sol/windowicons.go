package sol

import (
	"image"

	"github.com/fogleman/gg"
	"oddstream.games/gomps5/schriftbank"
)

func createImg(size float64) image.Image {
	var iSize int = int(size)
	var halfSize float64 = size / 2.0
	dc := gg.NewContext(iSize, iSize)
	dc.SetColor(ExtendedColors["BaizeGreen"])
	dc.DrawCircle(halfSize, halfSize, halfSize)
	dc.Fill()
	dc.SetColor(BasicColors["White"])
	dc.SetFontFace(schriftbank.Symbol24)
	dc.DrawStringAnchored(string(rune(9829)), halfSize, halfSize, 0.5, 0.4)
	dc.Stroke()
	return dc.Image()
}

func WindowIcons() []image.Image {
	var images []image.Image
	var sizes []float64 = []float64{16, 32, 48}
	for i := 0; i < len(sizes); i++ {
		images = append(images, createImg(sizes[i]))
	}
	return images
}
