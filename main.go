// Copyright ©️ 2021 oddstream.games

// $ go mod init oddstream.games/gosol
// $ go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"flag"
	"log"
	"os"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	sol "oddstream.games/gosol/sol"
)

func init() {
	println("processing command line flags")
	flag.BoolVar(&sol.DebugMode, "debug", false, "turn debug graphics on")
	flag.IntVar(&sol.WindowWidth, "width", 1000, "width of window in pixels")
	flag.IntVar(&sol.WindowHeight, "height", 900, "height of window in pixels")
	flag.IntVar(&sol.CardWidth, "cw", 71, "width of a card in pixels")
	flag.IntVar(&sol.CardHeight, "ch", 96, "height of a card in pixels")
	flag.StringVar(&sol.TheUserData.Variant, "v", "Klondike", "set the variant")
	flag.StringVar(&sol.TheUserData.CardStyle, "c", "retro", "set the card face to retro or scalable")
}

func main() {
	flag.Parse()

	if sol.DebugMode {
		for i, a := range os.Args {
			println(i, a)
		}
	}

	game, err := sol.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowTitle("Solitaire")                      // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowSize(sol.WindowWidth, sol.WindowHeight) // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowResizable(true)                         // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
