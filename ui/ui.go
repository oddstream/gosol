package ui

import (
	_ "embed" // go:embed only allowed in Go files that import "embed"

	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Container holds Widgets
// type Container struct {
// 	widgets         []*Widget
// 	BackgroundImage *ebiten.Image
// }

// UI encapsulates a complete user interface that can be rendered onto the screen.
type UI struct {
	toastTextFace   font.Face
	toastManager    *ToastManager
	toolbarTextFace font.Face
	toolbar         *Toolbar
}

//go:embed assets/Roboto-Regular.ttf
var robotoRegularBytes []byte

// New creates a new UI object
func New() *UI {
	ui := &UI{}

	tt, err := truetype.Parse(robotoRegularBytes)
	if err != nil {
		log.Fatal(err)
	}

	ui.toastTextFace = truetype.NewFace(tt, &truetype.Options{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	robotoRegularBytes = nil

	ui.toastManager = &ToastManager{}
	ui.toolbar = NewToolbar()

	return ui
}

// Update is called once per tick and updates the UI's state
func (u *UI) Update() {
	u.toolbar.Update()
	u.toastManager.Update()
}

// Draw is called once per tick and renders the UI to the screen
func (u *UI) Draw(screen *ebiten.Image) {
	u.toolbar.Draw(screen)
	u.toastManager.Draw(screen)
}
