package sol

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

var CommandTable = map[ebiten.Key]func(){
	ebiten.KeyN:      func() { TheBaize.NewGame() },
	ebiten.KeyR:      func() { TheBaize.RestartGame() },
	ebiten.KeyU:      func() { TheBaize.Undo() },
	ebiten.KeyS:      func() { TheBaize.SavePosition() },
	ebiten.KeyL:      func() { TheBaize.LoadPosition() },
	ebiten.KeyC:      func() { TheBaize.Collect() },
	ebiten.KeyF:      func() { TheBaize.ShowVariantPicker() },
	ebiten.KeyF1:     func() { TheBaize.ShowRules() },
	ebiten.KeyF2:     func() { TheUI.ShowCardBackPicker(TheCIP.BackImages()) },
	ebiten.KeyF3:     func() { ShowSettingsDrawer() },
	ebiten.KeyF4:     func() { TheStatistics.ShowStatistics() },
	ebiten.KeyF5:     func() { TheBaize.StartSpinning() },
	ebiten.KeyF6:     func() { TheBaize.StopSpinning() },
	ebiten.KeyF7:     func() { TheUI.ShowFAB("star", ebiten.KeyN) },
	ebiten.KeyF8:     func() { TheUI.HideFAB() },
	ebiten.KeyMenu:   func() { TheUI.ToggleNavDrawer() },
	ebiten.KeyEscape: func() { TheUI.HideActiveDrawer() },
	ebiten.KeyX:      func() { TheBaize.Exit() },
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
				println("ChangeRequest unknown variant", v.Data)
				break
			}
			if newVariant != TheBaize.Variant {
				TheBaize.Save()
				if !TheBaize.LoadVariant(newVariant) {
					TheBaize.NewVariant(newVariant)
				}
			}
		case "Retro cards":
			ThePreferences.RetroCards, _ = strconv.ParseBool(v.Data)
			TheBaize.OldWindowWidth = 0 // force a rescale
			TheBaize.Scale()
		case "Fixed cards":
			ThePreferences.FixedCards, _ = strconv.ParseBool(v.Data)
			TheBaize.OldWindowWidth = 0 // force a rescale
			TheBaize.Scale()
		case "CardBack":
			if ThePreferences.RetroCards {
				ThePreferences.CardBackPattern = v.Data
				CardBackImage = TheCIP.BackImage(ThePreferences.CardBackPattern)
			} else {
				ThePreferences.CardBackColor = v.Data
				CardBackImage = TheCIP.BackImage(ThePreferences.CardBackColor)
			}
		case "Single tap":
			ThePreferences.SingleTap, _ = strconv.ParseBool(v.Data)
		case "Highlights":
			ThePreferences.HighlightMovable, _ = strconv.ParseBool(v.Data)
		case "Power moves":
			ThePreferences.PowerMoves, _ = strconv.ParseBool(v.Data)
			TheBaize.MarkMovable() // re-run this
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
