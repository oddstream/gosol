// $ go mod init oddstream.games/gomps5
// $ go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"log"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	sol "oddstream.games/gomps5/sol"
	"oddstream.games/gomps5/sound"
)

func main() {

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

	defer func() {
		println("main defer cleanup")
		if !sol.NoGameSave {
			sol.TheBaize.Save()
		}
		sol.ThePreferences.Save()
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
