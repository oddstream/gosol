// $ go mod init oddstream.games/gosol
// $ go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"log"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	sol "oddstream.games/gosol/sol"
)

func main() {
	sol.NewGame() // sets sol.TheGame

	defer func() {
		log.Println("main defer cleanup")
		if !sol.NoGameSave {
			sol.TheBaize.Save()
		}
		sol.TheSettings.Save()
	}()

	if err := ebiten.RunGame(sol.TheGame); err != nil {
		log.Fatal(err)
	}
}
