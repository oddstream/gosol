package sol

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

type CardImages struct {
	faceImgs  map[CardID]*ebiten.Image
	backImgs  map[string]*ebiten.Image
	shadowImg *ebiten.Image
	// movableImg *ebiten.Image
}

type CardImageProvider interface {
	FaceImage(CardID) *ebiten.Image
	BackImage(string) *ebiten.Image
	BackImages() map[string]*ebiten.Image
	ShadowImage() *ebiten.Image
	// MovableImage() *ebiten.Image
}

func CreateCardImages() {
	schriftbank.MakeCardFonts(CardWidth) // CardWidth/Height have now been set

	if ThePreferences.RetroCards {
		TheCIP = NewRetroCardImageProvider()
		CardBackImage = TheCIP.BackImage(ThePreferences.CardBackPattern)
	} else {
		TheCIP = NewModernCardImageProvider()
		CardBackImage = TheCIP.BackImage(ThePreferences.CardBackColor)
	}
	CardShadowImage = TheCIP.ShadowImage()
	// CardMovableImage = TheCIP.MovableImage()
}
