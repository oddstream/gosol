package sol

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/ui"
)

var CommandTable = map[ebiten.Key]func(){
	ebiten.KeyN: func() { TheBaize.NewDeal() },
	ebiten.KeyR: func() { TheBaize.RestartDeal() },
	ebiten.KeyU: func() { TheBaize.Undo() },
	ebiten.KeyB: func() {
		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			TheBaize.LoadPosition()
		} else {
			TheBaize.SavePosition()
		}
	},
	ebiten.KeyL: func() { TheBaize.LoadPosition() },
	ebiten.KeyS: func() { TheBaize.SavePosition() },
	ebiten.KeyC: func() { TheBaize.Collect2() },
	ebiten.KeyH: func() {
		TheSettings.ShowMovableCards = !TheSettings.ShowMovableCards
		if TheSettings.ShowMovableCards {
			if TheBaize.moves+TheBaize.fmoves > 0 {
				TheUI.ToastInfo("Movable cards highlighted")
			} else {
				TheUI.ToastError("There are no movable cards")
			}
		}
	},
	ebiten.KeyM: func() {
		TheSettings.AlwaysShowMovableCards = !TheSettings.AlwaysShowMovableCards
		TheSettings.ShowMovableCards = TheSettings.AlwaysShowMovableCards
		if TheSettings.AlwaysShowMovableCards {
			TheUI.ToastInfo("Movable cards always highlighted")
		}
	},
	ebiten.KeyF: func() { TheUI.ShowVariantPickerEx(VariantGroupNames(), "ShowVariantPicker") },
	ebiten.KeyA: func() { ShowAniSpeedDrawer() },
	ebiten.KeyX: func() { ExitRequested = true },
	// ebiten.KeyTab: func() {
	// 	if DebugMode {
	// 		for _, p := range TheBaize.piles {
	// 			p.Refan()
	// 		}
	// 	}
	// },
	ebiten.KeyF1: func() { TheBaize.Wikipedia() },
	ebiten.KeyF2: func() { ShowStatisticsDrawer() },
	ebiten.KeyF3: func() { ShowSettingsDrawer() },
	ebiten.KeyF5: func() { TheBaize.StartSpinning() }, // debug
	ebiten.KeyF6: func() { TheBaize.StopSpinning() },  // debug
	ebiten.KeyF7: func() {
		TheUI.AddButtonToFAB("restore", ebiten.KeyR)
		TheUI.AddButtonToFAB("done_all", ebiten.KeyC)
	}, // debug
	ebiten.KeyF8:     func() { TheUI.HideFAB() }, // debug
	ebiten.KeyMenu:   func() { TheUI.ToggleNavDrawer() },
	ebiten.KeyEscape: func() { TheUI.HideActiveDrawer() },
}

func Execute(cmd interface{}) {
	TheUI.HideActiveDrawer()
	TheUI.HideFAB()
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := CommandTable[v]; ok {
			fn()
		}
	case ui.Command:
		// a widget has sent a command
		switch v.Command {
		case "ShowVariantGroupPicker":
			TheUI.ShowVariantPickerEx(VariantGroupNames(), "ShowVariantPicker")
		case "ShowVariantPicker":
			TheUI.ShowVariantPickerEx(VariantNames(v.Data), "ChangeVariant")
		case "ChangeVariant":
			if _, ok := Variants[v.Data]; !ok {
				TheUI.ToastError(fmt.Sprintf("Don't know how to play '%s'", v.Data))
			} else if v.Data == TheSettings.Variant {
				TheUI.ToastError(fmt.Sprintf("Already playing '%s'", v.Data))
			} else {
				TheBaize.ChangeVariant(v.Data)
				TheSettings.Save() // save now especially if running in a browser
			}
		case "SaveSettings":
			TheSettings.Save() // save now especially if running in a browser
		default:
			log.Panic("unknown command", v.Command, v.Data)
		}
	default:
		log.Fatal("Baize.Execute unknown command type", cmd)
	}
}
