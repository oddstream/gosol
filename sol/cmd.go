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
		case "CardBack":
			if TheUserData.CardStyle == "retro" {
				TheUserData.CardBackPattern = v.Data
				CardBackImage = TheCIP.BackImage(TheUserData.CardBackPattern)
			} else {
				TheUserData.CardBackColor = v.Data
				CardBackImage = TheCIP.BackImage(TheUserData.CardBackColor)
			}
		case "Single tap":
			TheUserData.SingleTap, _ = strconv.ParseBool(v.Data)
		case "Highlights":
			TheUserData.HighlightMovable, _ = strconv.ParseBool(v.Data)
		case "Power moves":
			TheUserData.PowerMoves, _ = strconv.ParseBool(v.Data)
			b.MarkMovable() // re-run this
		case "Scaled cards":
			if scaled, _ := strconv.ParseBool(v.Data); scaled {
				TheUserData.CardStyle = "scaled"
			}
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		case "Fixed cards":
			if fixed, _ := strconv.ParseBool(v.Data); fixed {
				TheUserData.CardStyle = "fixed"
			}
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		case "Retro cards":
			if retro, _ := strconv.ParseBool(v.Data); retro {
				TheUserData.CardStyle = "retro"
			}
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		case "Mute sounds":
			TheUserData.MuteSounds, _ = strconv.ParseBool(v.Data)
			sound.Mute(TheUserData.MuteSounds)
		default:
			log.Panic("unknown change request", v.ChangeRequested, v.Data)
		}
		TheUserData.Save() // save now especially if running on a browser

	default:
		log.Fatal("Baize.Execute unknown command type", cmd)
	}
}
