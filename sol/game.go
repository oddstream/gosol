// Package sol provides a polymorphic solitaire engine
package sol

import (
	"errors"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

// Game represents a game state
type Game struct {
	UI         *ui.UI
	Baize      *Baize
	Statistics *Statistics
	Settings   *Settings
}

var (
	// GosolVersionMajor is the integer version number
	GosolVersionMajor int = 5
	// CsolVersionMinor is the integer version number
	GosolVersionMinor int = 15
	// CSolVersionDate is the ISO 8601 date of bumping the version number
	GosolVersionDate string = "2023-02-26"
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoGameLoad is a boolean set by command line flag -noload
	NoGameLoad bool = false
	// NoGameSave is a boolean set by command line flag -nosave
	NoGameSave bool = false
	// NoScrunch stops cards being scrunched
	NoScrunch bool = false
	// CardWidth of cards, start with a silly value to force a rescale/refan
	CardWidth int = 9
	// CardHeight of cards, start with a silly value to force a rescale/refan
	CardHeight int = 13
	// CardDiagonal float64 = 15.8
	// Card Corner Radius
	CardCornerRadius float64 = float64(CardWidth) / 10.0
	// PilePaddingX the gap left to the right of the pile
	PilePaddingX int = CardWidth / 10
	// PilePaddingY the gap left underneath each pile
	PilePaddingY int = CardHeight / 10
	// LeftMargin the gap between the left of the screen and the first pile
	LeftMargin int = (CardWidth / 2) + PilePaddingX
	// TopMargin the gap between top pile and top of baize
	TopMargin int = ui.ToolbarHeight + CardHeight/3
	// CardFaceImageLibrary
	// thirteen suitless cards,
	// one entry for each face card (4 suits * 13 cards),
	// suits are 1-indexed (eg club == 1) so image to be used for a card is (suit * 13) + (ord - 1).
	// can use (ord - 1) as an index to get suitless card
	TheCardFaceImageLibrary [13 * 5]*ebiten.Image
	// CardBackImage applies to all cards so is kept globally as an optimization
	CardBackImage *ebiten.Image
	// MovableCardBackImage applies to all cards so is kept globally as an optimization
	MovableCardBackImage *ebiten.Image
	// CardShadowImage applies to all cards so is kept globally as an optimization
	CardShadowImage *ebiten.Image
	// ExitRequested is set when user has had enough
	ExitRequested bool = false
)

var TheGame *Game // pointer to object that implements ebiten.Game interface

// NewGame generates a new Game object, which implements ebiten.Game interface
func NewGame() {
	TheGame = &Game{Settings: NewSettings()}
	if TheGame.Settings.Mute {
		sound.SetVolume(0.0)
	} else {
		sound.SetVolume(TheGame.Settings.Volume)
	}
	TheGame.Statistics = NewStatistics()
	TheGame.UI = ui.New(Execute)
	if TheGame.Baize = NewBaize(TheGame.Settings.Variant); TheGame.Baize == nil {
		log.Panic("cannot create Baize")
	}
	TheGame.Baize.StartFreshGame()
	if !NoGameLoad {
		if undoStack := LoadUndoStack(TheGame.Settings.Variant); TheGame.Baize.IsSavableStackOk(undoStack) {
			TheGame.Baize.SetUndoStack(undoStack)
		}
	}

	if TheGame.Settings.LastVersionMajor != GosolVersionMajor || TheGame.Settings.LastVersionMinor != GosolVersionMinor {
		TheGame.UI.Toast("Glass", fmt.Sprintf("Upgraded from %d.%d to %d.%d",
			TheGame.Settings.LastVersionMajor,
			TheGame.Settings.LastVersionMinor,
			GosolVersionMajor,
			GosolVersionMinor))
	}
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Baize.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// Update updates the current game state.
// In Ebiten, the Update function is part of a fixed timestep loop,
// and it's called based on the ticks per second (TPS) of your game.
// By default, Ebiten's TPS is 60, which means that the Update method
// will be called 60 times per second. In other words, unless you modify
// the TPS with SetMaxTPS, the fixed timestep will be 1000/60 = 16.666 milliseconds.
// https://ebitencookbook.vercel.app/blog
func (g *Game) Update() error {
	g.Baize.Update()
	if ExitRequested {
		if !NoGameSave {
			g.Baize.Save(TheGame.Settings.Variant)
		}
		g.Settings.Save()
		return errors.New("exit requested")
	}
	return nil
}

// Draw draws the current game to the given screen.
// Draw will be called based on the refresh rate of the screen (FPS).
// https://ebitencookbook.vercel.app/blog
func (g *Game) Draw(screen *ebiten.Image) {
	g.Baize.Draw(screen)
}
