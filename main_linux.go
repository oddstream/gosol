// $ go mod init oddstream.games/gosol
// $ go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"flag"
	"fmt"
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

func main() {

	log.SetFlags(0)

	// load userdata before processing flags, because flags can override userdata
	sol.TheUserData.Load()

	// pearl from the mudbank: don't have any flags that will overwrite TheUserData
	flag.BoolVar(&sol.DebugMode, "debug", false, "turn debug graphics on")
	flag.BoolVar(&sol.NoGameLoad, "noload", false, "do not load saved game when starting")
	flag.BoolVar(&sol.NoGameSave, "nosave", false, "do not save game before exit")
	flag.BoolVar(&sol.NoShuffle, "noshuffle", false, "do not shuffle cards")
	flag.BoolVar(&ui.GenerateIcons, "generateicons", false, "generate icon files")

	flag.Parse()

	sound.Mute(sol.TheUserData.MuteSounds)

	if sol.DebugMode {
		for i, a := range os.Args {
			println(i, a)
		}
	}

	game, err := sol.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	// ebiten default window size is 640, 480
	if sol.TheUserData.WindowWidth == 0 || sol.TheUserData.WindowHeight == 0 {
		// not yet set/saved, so use sensible values
		sol.TheUserData.WindowWidth, sol.TheUserData.WindowHeight = ebiten.ScreenSizeInFullscreen()
		sol.TheUserData.WindowWidth /= 2
		sol.TheUserData.WindowHeight /= 2
	}
	ebiten.SetWindowSize(sol.TheUserData.WindowWidth, sol.TheUserData.WindowHeight)

	if sol.TheUserData.WindowX != 0 && sol.TheUserData.WindowY != 0 {
		ebiten.SetWindowPosition(sol.TheUserData.WindowX, sol.TheUserData.WindowY)
	}
	ebiten.SetWindowResizable(true) //ebiten panics if a window to maximize is not resizable
	if sol.TheUserData.WindowMaximized {
		ebiten.MaximizeWindow()
	}
	if sol.DebugMode {
		ebiten.SetWindowTitle(fmt.Sprintf("Oddstream Solitaire (%s/%s)", runtime.GOOS, runtime.GOARCH))
	} else {
		ebiten.SetWindowTitle("Oddstream Solitaire")
	}
	ebiten.SetWindowIcon(sol.WindowIcons())
	// ebiten.SetScreenClearedEveryFrame(false)

	defer func() {
		println("main defer cleanup")
		if !sol.NoGameSave {
			sol.TheBaize.Save()
		}
		// calling ebiten.* functions here causes panic
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
