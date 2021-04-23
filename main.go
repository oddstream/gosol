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
	"runtime"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	sol "oddstream.games/gosol/sol"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

func init() {
	if runtime.GOARCH != "wasm" {
		println("processing command line flags")
		flag.BoolVar(&sol.DebugMode, "debug", false, "turn debug graphics on")
		flag.BoolVar(&sol.NoGameLoad, "noload", false, "do not load saved game when starting")
		flag.BoolVar(&sol.NoGameSave, "nosave", false, "do not save game before exit")
		flag.BoolVar(&sol.NoShuffle, "noshuffle", false, "do not shuffle cards")
		flag.StringVar(&sol.TheUserData.Variant, "v", "Klondike", "set the variant")
		flag.StringVar(&sol.TheUserData.CardStyle, "c", "default", "set the card face to retro")
		flag.BoolVar(&ui.GenerateIcons, "generateicons", false, "generate icon files")
	}
}

func main() {

	if runtime.GOARCH != "wasm" {
		log.SetFlags(0)
	}

	sol.TheUserData.Load()
	sound.Mute(sol.TheUserData.MuteSounds)

	if runtime.GOARCH != "wasm" {
		flag.Parse()
	}

	if sol.DebugMode {
		for i, a := range os.Args {
			println(i, a)
		}
	}

	game, err := sol.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	// default window size is 640, 480
	ebiten.SetWindowSize(1920/2, 1000) // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowResizable(true)    // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowTitle("Solitaire") // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetScreenClearedEveryFrame(false)

	defer func() {
		println("cleanup")
		if !sol.NoGameSave {
			sol.TheBaize.Save()
		}
		sol.TheUserData.Save()
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
