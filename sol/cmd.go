package sol

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

func (b *Baize) Execute(cmd interface{}) {
	switch v := cmd.(type) {
	case ebiten.Key:
		if fn, ok := b.commandTable[v]; ok {
			b.ui.HideActiveDrawer()
			b.ui.HideFAB()
			fn()
		}

	case ui.ChangeRequest:
		// a widget has sent a change request
		b.ui.HideActiveDrawer()
		b.ui.HideFAB()
		switch v.ChangeRequested {
		case "Variant":
			newVariant := v.Data
			if newVariant == "" {
				println("ChangeRequest unknown variant", v.Data)
				break
			}
			if newVariant != b.Variant {
				b.Save()
				if !TheBaize.LoadVariant(newVariant) {
					b.NewVariant(newVariant)
				}
			}
		case "Retro cards":
			ThePreferences.RetroCards, _ = strconv.ParseBool(v.Data)
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		case "Fixed cards":
			ThePreferences.FixedCards, _ = strconv.ParseBool(v.Data)
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
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
			b.MarkMovable() // re-run this
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
