//go:build android

package mobilegosol

import (

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/mobile"
	sol "oddstream.games/gosol/sol"
)

func init() {
	sol.NewGame() // sets sol.TheGame
	mobile.SetGame(sol.TheGame)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
