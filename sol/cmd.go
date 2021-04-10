package sol

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
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
				println("unknown variant", v.Data)
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
		case "Show hints":
			TheUserData.HighlightMovable, _ = strconv.ParseBool(v.Data)
		case "Power moves":
			TheUserData.PowerMoves, _ = strconv.ParseBool(v.Data)
			b.MarkMovable() // re-run this
		case "Retro cards":
			retro, _ := strconv.ParseBool(v.Data)
			if retro {
				TheUserData.CardStyle = "retro"
			} else {
				TheUserData.CardStyle = "default"
			}
			b.OldWindowWidth = 0 // force a rescale
			b.Scale()
		case "Mute sounds":
			b.ui.Toast("Not implemented yet")
		default:
			log.Panic("unknown change request", v.ChangeRequested, v.Data)
		}
		TheUserData.Save() // save now especially if running on a browser
	}
}
