package sol

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/ui"
)

var CommandTable = map[ebiten.Key]func(){
	ebiten.KeyN: func() { TheGame.Baize.NewDeal() },
	ebiten.KeyR: func() { TheGame.Baize.RestartDeal() },
	ebiten.KeyU: func() { TheGame.Baize.Undo() },
	ebiten.KeyB: func() {
		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			TheGame.Baize.LoadPosition()
		} else {
			TheGame.Baize.SavePosition()
		}
	},
	ebiten.KeyL: func() { TheGame.Baize.LoadPosition() },
	ebiten.KeyS: func() { TheGame.Baize.SavePosition() },
	ebiten.KeyC: func() { TheGame.Baize.Collect2() },
	ebiten.KeyH: func() {
		TheGame.Settings.ShowMovableCards = !TheGame.Settings.ShowMovableCards
		if TheGame.Settings.ShowMovableCards {
			if TheGame.Baize.moves+TheGame.Baize.fmoves > 0 {
				TheGame.UI.ToastInfo("Movable cards highlighted")
			} else {
				TheGame.UI.ToastError("There are no movable cards")
			}
		}
	},
	ebiten.KeyM: func() {
		TheGame.Settings.AlwaysShowMovableCards = !TheGame.Settings.AlwaysShowMovableCards
		TheGame.Settings.ShowMovableCards = TheGame.Settings.AlwaysShowMovableCards
		if TheGame.Settings.AlwaysShowMovableCards {
			TheGame.UI.ToastInfo("Movable cards always highlighted")
		}
	},
	ebiten.KeyF: func() { TheGame.UI.ShowVariantPickerEx(VariantGroupNames(), "ShowVariantPicker") },
	ebiten.KeyA: func() { ShowAniSpeedDrawer() },
	ebiten.KeyX: func() { ExitRequested = true },
	// ebiten.KeyTab: func() {
	// 	if DebugMode {
	// 		for _, p := range TheGame.Baize.piles {
	// 			p.Refan()
	// 		}
	// 	}
	// },
	ebiten.KeyF1: func() { TheGame.Baize.Wikipedia() },
	ebiten.KeyF2: func() { ShowStatisticsDrawer() },
	ebiten.KeyF3: func() { ShowSettingsDrawer() },
	ebiten.KeyF5: func() { TheGame.Baize.StartSpinning() }, // debug
	ebiten.KeyF6: func() { TheGame.Baize.StopSpinning() },  // debug
	ebiten.KeyF7: func() {
		TheGame.UI.AddButtonToFAB("restore", ebiten.KeyR)
		TheGame.UI.AddButtonToFAB("done_all", ebiten.KeyC)
	}, // debug
	ebiten.KeyF8:     func() { TheGame.UI.HideFAB() }, // debug
	ebiten.KeyMenu:   func() { TheGame.UI.ToggleNavDrawer() },
	ebiten.KeyEscape: func() { TheGame.UI.HideActiveDrawer() },
}

func Execute(cmd interface{}) {
	TheGame.UI.HideActiveDrawer()
	TheGame.UI.HideFAB()
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := CommandTable[v]; ok {
			fn()
		}
	case ui.Command:
		// a widget has sent a command
		switch v.Command {
		case "ShowVariantGroupPicker":
			TheGame.UI.ShowVariantPickerEx(VariantGroupNames(), "ShowVariantPicker")
		case "ShowVariantPicker":
			TheGame.UI.ShowVariantPickerEx(VariantNames(v.Data), "ChangeVariant")
		case "ChangeVariant":
			if _, ok := Variants[v.Data]; !ok {
				TheGame.UI.ToastError(fmt.Sprintf("Don't know how to play '%s'", v.Data))
			} else if v.Data == TheGame.Settings.Variant {
				TheGame.UI.ToastError(fmt.Sprintf("Already playing '%s'", v.Data))
			} else {
				TheGame.Baize.ChangeVariant(v.Data)
				TheGame.Settings.Save() // save now especially if running in a browser
			}
		case "SaveSettings":
			TheGame.Settings.Save() // save now especially if running in a browser
		default:
			log.Panic("unknown command", v.Command, v.Data)
		}
	default:
		log.Fatal("Baize.Execute unknown command type", cmd)
	}
}
