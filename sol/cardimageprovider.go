package sol

import "github.com/hajimehoshi/ebiten/v2"

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
