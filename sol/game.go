// Package sol provides a polymorphic solitaire engine
package sol

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomps5/ui"
)

// Game represents a game state
type Game struct {
}

var (
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoDrawing is set when resizing cards to stop the screen flickering
	// NoDrawing bool = false
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
	// ScreenWidth,ScreenHeight  of screen at startup
	ScreenWidth, ScreenHeight int
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
	// CardShadowImage applies to all cards so is kept globally as an optimization
	CardShadowImage *ebiten.Image
	// CardHighlightImage applies to all cards so is kept globally as an optimization
	CardHighlightImage *ebiten.Image
	// ExitRequested is set when user has had enough
	ExitRequested bool = false
	// InGameLoop is true when ebiten is running
	InGameLoop = false
)

// ThePreferences holds serialized game progress data
// Colors are named from the web extended colors at https://en.wikipedia.org/wiki/Web_colors
var ThePreferences = &Preferences{
	Title:           "Solitaire",
	Variant:         "Klondike",
	BaizeColor:      "BaizeGreen",
	PowerMoves:      true,
	CardFaceColor:   "Ivory",
	CardBackColor:   "CornflowerBlue",
	ExtraColors:     false,
	RedColor:        "Crimson",
	BlackColor:      "Black",
	ClubColor:       "Indigo",
	DiamondColor:    "OrangeRed",
	HeartColor:      "Crimson",
	SpadeColor:      "Black",
	FixedCards:      true,
	Mute:            false,
	Volume:          0.5,
	FixedCardWidth:  90,
	FixedCardHeight: 122,
	CardRatio:       1.357,
}

// TheStatistics holds statistics for all variants
var TheStatistics *Statistics

// TheBaize points to the Baize, so that main can see it
var TheBaize *Baize

// The UI points to the singleton User Interface object
var TheUI *ui.UI

// CardLibrary is the slice where Card objects actually exist, everything else is a *Card
var CardLibrary []Card

// TheError is a global copy of the last error reported, for optional toasting
// var TheError string

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}
	TheUI = ui.New(Execute)
	TheStatistics = NewStatistics()
	TheBaize = NewBaize()
	TheBaize.StartFreshGame()
	return g, nil
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
