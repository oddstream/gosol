// Copyright ©️ 2021 oddstream.games

package sol

import (
	"bytes"
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"image"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/fogleman/gg"
)

//go:embed assets/raccoon280x180.png
var logoBytes []byte

// Splash represents a game state.
type Splash struct {
	circleImage *ebiten.Image
	logoImage   *ebiten.Image
	circlePos   image.Point
	logoPos     image.Point
	skew        float64
}

// NewSplash creates and initializes a Splash/GameState object
func NewSplash() *Splash {
	s := &Splash{}

	dc := gg.NewContext(400, 400)
	dc.SetRGB(0.25, 0.25, 0.25)
	dc.DrawCircle(200, 200, 120)
	dc.Fill()
	dc.Stroke()
	img := dc.Image()
	s.circleImage = ebiten.NewImageFromImage(img)

	img, _, err := image.Decode(bytes.NewReader(logoBytes))
	if err != nil {
		log.Fatal(err)
	}
	s.logoImage = ebiten.NewImageFromImage(img)

	return s
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	yCenter := outsideHeight / 2

	cx, cy := s.circleImage.Size()
	s.circlePos = image.Point{X: xCenter - (cx / 2), Y: yCenter - (cy / 2)}

	lx, ly := s.logoImage.Size()
	s.logoPos = image.Point{X: xCenter - (lx / 2), Y: yCenter - 4 - (ly / 2)}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (s *Splash) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		os.Exit(0)
	}

	if s.skew < 90 {
		s.skew++
	} else {
		GSM.Switch(NewBaize())
	}

	return nil
}

// Draw draws the current GameState to the given screen
func (s *Splash) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	skewRadians := s.skew * math.Pi / 180

	{
		op := &ebiten.DrawImageOptions{}
		sx, sy := s.circleImage.Size()
		sx, sy = sx/2, sy/2
		op.GeoM.Translate(float64(-sx), float64(-sy))
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Skew(skewRadians, skewRadians)
		op.GeoM.Translate(float64(sx), float64(sy))

		op.GeoM.Translate(float64(s.circlePos.X), float64(s.circlePos.Y))
		screen.DrawImage(s.circleImage, op)
	}
	{
		op := &ebiten.DrawImageOptions{}
		sx, sy := s.logoImage.Size()
		sx, sy = sx/2, sy/2
		op.GeoM.Translate(float64(-sx), float64(-sy))
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Skew(skewRadians, skewRadians)
		op.GeoM.Translate(float64(sx), float64(sy))

		op.GeoM.Translate(float64(s.logoPos.X), float64(s.logoPos.Y))

		screen.DrawImage(s.logoImage, op)
	}
}
