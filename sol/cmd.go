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
	ebiten.KeyS: func() { TheBaize.SavePosition() },
	ebiten.KeyL: func() { TheBaize.LoadPosition() },
	ebiten.KeyC: func() { TheBaize.Collect2() },
	ebiten.KeyH: func() {
		ThePreferences.ShowMovableCards = !ThePreferences.ShowMovableCards
		if ThePreferences.ShowMovableCards {
			if TheBaize.moves+TheBaize.fmoves > 0 {
				TheUI.ToastInfo("Movable cards highlighted")
			} else {
				TheUI.ToastError("There are no movable cards")
			}
		}
	},
	ebiten.KeyM: func() {
		ThePreferences.AlwaysShowMovableCards = !ThePreferences.AlwaysShowMovableCards
		ThePreferences.ShowMovableCards = ThePreferences.AlwaysShowMovableCards
		if ThePreferences.AlwaysShowMovableCards {
			TheUI.ToastInfo("Movable cards always highlighted")
		}
	},
	ebiten.KeyF: func() { TheUI.ShowVariantPickerEx(VariantGroupNames(), "VariantGroup") },
	ebiten.KeyA: func() { ShowAniSpeedDrawer() },
	ebiten.KeyX: func() { ExitRequested = true },
	// ebiten.KeyTab: func() {
	// 	if DebugMode {
	// 		for _, p := range TheBaize.piles {
	// 			p.Refan()
	// 		}
	// 		ThePreferences.Save()
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
	case ui.ChangeRequest:
		// a widget has sent a change request
		switch v.ChangeRequested {
		case "Variant":
			if _, ok := Variants[v.Data]; !ok {
				TheUI.ToastError(fmt.Sprintf("Don't know how to play '%s'", v.Data))
			} else {
				if v.Data != ThePreferences.Variant {
					TheBaize.ChangeVariant(v.Data)
				}
			}
		case "VariantGroup":
			TheUI.ShowVariantPickerEx(VariantNames(v.Data), "Variant")
			// TheBaize.ShowVariantPicker(v.Data)
		default:
			log.Panic("unknown change request", v.ChangeRequested, v.Data)
		}
		ThePreferences.Save() // save now especially if running on a browser
	default:
		log.Fatal("Baize.Execute unknown command type", cmd)
	}
}
