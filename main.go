//go:build linux || windows

// https://go.googlesource.com/proposal/+/master/design/draft-gobuild.md
// $ go mod init oddstream.games/gosol
// $ go mod tidy

// the package defining a command (an executable Go program) always has the name main
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
	"oddstream.games/gosol/ui"
	"oddstream.games/gosol/util"
)

func main() {

	log.SetFlags(0)

	// pearl from the mudbank: don't have any flags that will overwrite ThePreferences
	flag.BoolVar(&sol.DebugMode, "debug", false, "turn debug graphics on")
	flag.BoolVar(&sol.NoGameLoad, "noload", false, "do not load saved game when starting")
	flag.BoolVar(&sol.NoGameSave, "nosave", false, "do not save game before exit")
	flag.BoolVar(&sol.NoScrunch, "noscrunch", false, "do not scrunch cards")
	flag.BoolVar(&ui.GenerateIcons, "generateicons", false, "generate icon files")

	flag.Parse()

	if sol.DebugMode {
		for i, a := range os.Args {
			log.Println(i, a)
		}
	}

	// ebiten panics if a window to maximize is not resizable
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if ebiten.IsWindowMaximized() || ebiten.IsWindowMinimized() {
		// GNOME (maybe) annoyingly keeps maximizing the window
		ebiten.RestoreWindow()
	}
	{
		x, y := ebiten.ScreenSizeInFullscreen()
		n := util.Max(x, y)
		ebiten.SetWindowSize(n/2, n/2)
	}
	ebiten.SetWindowIcon(sol.WindowIcons())
	ebiten.SetWindowTitle("Go Solitaire")

	game, err := sol.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if !sol.NoGameLoad {
		if undoStack := sol.LoadUndoStack(); sol.TheBaize.IsSavableStackOk(undoStack) {
			sol.TheBaize.SetUndoStack(undoStack) // TODO doing this before Baize.Layout sets WindowWidth,Height
		}
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	// we come here if the user closed the window with the x button
	// println("main exit")

	if !sol.NoGameSave {
		sol.TheBaize.Save()
	}

	sol.TheSettings.Save()
}
