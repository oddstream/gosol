package sol

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/sound"
	"oddstream.games/gomps5/ui"
)

var CommandTable = map[ebiten.Key]func(){
	ebiten.KeyN:   func() { TheBaize.NewDeal() },
	ebiten.KeyR:   func() { TheBaize.RestartDeal() },
	ebiten.KeyU:   func() { TheBaize.Undo() },
	ebiten.KeyS:   func() { TheBaize.SavePosition() },
	ebiten.KeyL:   func() { TheBaize.LoadPosition() },
	ebiten.KeyC:   func() { TheBaize.Collect() },
	ebiten.KeyA:   func() { TheBaize.CollectAll() },
	ebiten.KeyF:   func() { TheBaize.ShowVariantPicker() },
	ebiten.KeyTab: func() { TheBaize.Refan() },
	ebiten.KeyF1:  func() { TheBaize.Wikipedia() },
	ebiten.KeyF3:  func() { ShowSettingsDrawer() },
	// ebiten.KeyF4: func() { TheStatistics.ShowStatistics() },
	ebiten.KeyF5:     func() { TheBaize.StartSpinning() },
	ebiten.KeyF6:     func() { TheBaize.StopSpinning() },
	ebiten.KeyF7:     func() { TheUI.ShowFAB("star", ebiten.KeyN) },
	ebiten.KeyF8:     func() { TheUI.HideFAB() },
	ebiten.KeyMenu:   func() { TheUI.ToggleNavDrawer() },
	ebiten.KeyEscape: func() { TheUI.HideActiveDrawer() },
}

func Execute(cmd interface{}) {
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := CommandTable[v]; ok {
			TheUI.HideActiveDrawer()
			TheUI.HideFAB()
			fn()
		}

	case ui.ChangeRequest:
		// a widget has sent a change request
		TheUI.HideActiveDrawer()
		TheUI.HideFAB()
		switch v.ChangeRequested {
		case "Variant":
			newVariant := v.Data
			if newVariant == "" {
				log.Panic("ChangeRequest empty variant", v.Data)
			}
			if newVariant != ThePreferences.Variant {
				ThePreferences.Variant = newVariant
				TheBaize.NewVariant()
			}
		case "Fixed cards":
			ThePreferences.FixedCards, _ = strconv.ParseBool(v.Data)
			TheBaize.setFlag(dirtyCardSizes | dirtyPileBackgrounds | dirtyPilePositions | dirtyCardPositions)
		case "Power moves":
			ThePreferences.PowerMoves, _ = strconv.ParseBool(v.Data)
		case "Mirror baize":
			ThePreferences.MirrorBaize, _ = strconv.ParseBool(v.Data)
			savedUndoStack := TheBaize.undoStack
			TheBaize.NewVariant()
			TheBaize.SetUndoStack(savedUndoStack)
		case "Mute sounds":
			ThePreferences.MuteSounds, _ = strconv.ParseBool(v.Data)
			sound.Mute(ThePreferences.MuteSounds)
		default:
			log.Panic("unknown change request", v.ChangeRequested, v.Data)
		}
		ThePreferences.Save() // save now especially if running on a browser

	default:
		log.Fatal("Baize.Execute unknown command type", cmd)
	}
}
