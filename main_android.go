//go:build android

package main

import (
	"log"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/mobile"
	sol "oddstream.games/gomps5/sol"
)

func init() {

	sol.ThePreferences.Load()

	if sol.ThePreferences.Mute {
		sound.SetVolume(0.0)
	} else {
		sound.SetVolume(sol.ThePreferences.Volume)
	}

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
