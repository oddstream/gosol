package sol

import "github.com/hajimehoshi/ebiten/v2"

type CardImageProvider interface {
	FaceImage(CardID) *ebiten.Image
	BackImage(string) *ebiten.Image
	BackImages() map[string]*ebiten.Image
	ShadowImage() *ebiten.Image
	// MovableImage() *ebiten.Image
}
