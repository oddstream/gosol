package sol

import (
	"image"

	"github.com/fogleman/gg"
)

func createImg(size int) image.Image {
	var s float64 = float64(size)
	var halfSize float64 = s / 2.0
	dc := gg.NewContext(size, size)
	dc.SetColor(ExtendedColors["BaizeGreen"])
	dc.DrawCircle(halfSize, halfSize, halfSize)
	dc.Fill()

	dc.SetColor(ExtendedColors["Ivory"])
	// draw a scaled diamond (simplest suit shape)
	dc.MoveTo(s*0.5, s*0.2)  // top
	dc.LineTo(s*0.75, s*0.5) // right
	dc.LineTo(s*0.5, s*0.8)  // bottom
	dc.LineTo(s*0.25, s*0.5) // left
	dc.LineTo(s*0.5, s*0.25) // top
	dc.Fill()

	dc.Stroke()
	return dc.Image()
}

func WindowIcons() []image.Image {
	var images []image.Image
	var sizes []int = []int{16, 32, 48, 96, 128}
	for i := 0; i < len(sizes); i++ {
		images = append(images, createImg(sizes[i]))
	}
	return images
}
