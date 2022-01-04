//go:build android

package mobilegomps5

import (
	"log"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/mobile"
	sol "oddstream.games/gomps5/sol"
)

func init() {

	game, err := sol.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	mobile.SetGame(game)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
