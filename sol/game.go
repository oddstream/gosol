// Package sol provides a polymorphic solitaire engine
package sol

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/ui"
)

// Game represents a game state
type Game struct {
}

var (
	// GosolVersionMajor is the integer version number
	GosolVersionMajor int = 5
	// CsolVersionMinor is the integer version number
	GosolVersionMinor int = 8
	// CSolVersionDate is the ISO 8601 date of bumping the version number
	GosolVersionDate string = "2023-01-02"
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoGameLoad is a boolean set by command line flag -noload
	NoGameLoad bool = false
	// NoGameSave is a boolean set by command line flag -nosave
	NoGameSave bool = false
	// NoShuffle stops the cards from being shuffled
	NoShuffle bool = false
	// NoScrunch stops cards being scrunched
	NoScrunch bool = false
	// NoCardLerp stops the cards from transitioning
	NoCardLerp = false
	// NoCardFlip stops the cards from animating their flip
	NoCardFlip = false
	// CardWidth of cards, start with a silly value to force a rescale/refan
	CardWidth int = 9
	// CardHeight of cards, start with a silly value to force a rescale/refan
	CardHeight int = 13
	// Card Corner Radius
	CardCornerRadius float64 = float64(CardWidth) / 15.0
	// PilePaddingX the gap left to the right of the pile
	PilePaddingX int = CardWidth / 10
	// PilePaddingY the gap left underneath each pile
	PilePaddingY int = CardHeight / 10
	// LeftMargin the gap between the left of the screen and the first pile
	LeftMargin int = (CardWidth / 2) + PilePaddingX
	// TopMargin the gap between top pile and top of baize
	TopMargin int = 48 + CardHeight/3
	// CardFaceImageLibrary
	// thirteen suitless cards,
	// one entry for each face card (4 suits * 13 cards),
	// suits are 1-indexed (eg club == 1) so image to be used for a card is (suit * 13) + (ord - 1).
	// can use (ord - 1) as in index to get suitless card
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

// TheStatistics holds statistics for all variants
var TheStatistics *Statistics

// TheBaize points to the Baize, so that main can see it
var TheBaize *Baize

// The UI points to the singleton User Interface object
var TheUI *ui.UI

// CardLibrary is the slice where Card objects actually exist, everything else is a *Card
var CardLibrary []Card

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	ThePreferences.Load()
	if ThePreferences.Mute {
		sound.SetVolume(0.0)
	} else {
		sound.SetVolume(ThePreferences.Volume)
	}
	TheUI = ui.New(Execute)
	TheStatistics = NewStatistics()
	TheBaize = NewBaize()
	TheBaize.StartFreshGame()
	if ThePreferences.LastVersionMajor != GosolVersionMajor || ThePreferences.LastVersionMinor != GosolVersionMinor {
		TheUI.Toast("Glass", fmt.Sprintf("Upgraded from %d.%d to %d.%d",
			ThePreferences.LastVersionMajor,
			ThePreferences.LastVersionMinor,
			GosolVersionMajor,
			GosolVersionMinor))
	}
	return &Game{}, nil
}

// Layout implements ebiten.Game's Layout.
func (*Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	TheBaize.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (*Game) Update() error {
	TheBaize.Update()
	if ExitRequested {
		if !NoGameSave {
			TheBaize.Save()
		}
		ThePreferences.Save()
		return errors.New("exit requested")
	}
	return nil
}

// Draw draws the current game to the given screen.
func (*Game) Draw(screen *ebiten.Image) {
	TheBaize.Draw(screen)
}
