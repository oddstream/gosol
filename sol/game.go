// Package sol provides a polymorphic solitaire engine
package sol

import (
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/ui"
)

// Game represents a game state
type Game struct {
}

var (
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoGameLoad is a boolean set by command line flag -noload
	NoGameLoad bool = false
	// NoGameSave is a boolean set by command line flag -nosave
	NoGameSave bool = false
	// NoShuffle stops the cards from being shuffled
	NoShuffle bool = false
	// CardWidth of cards
	CardWidth int = 71
	// CardHeight of cards
	CardHeight int = 96
	// PilePaddingX the gap left to the right of the pile
	PilePaddingX int = CardWidth / 10
	// PilePaddingY the gap left underneath each pile
	PilePaddingY int = CardHeight / 10
	// LeftMargin the gap between the left of the screen and the first pile
	LeftMargin int = (CardWidth / 2) + PilePaddingX
	// TopMargin the gap between top pile and top of baize
	TopMargin int = 48 + CardHeight/3
	// CardBackImage applies to all cards so is kept globally as an optimization
	CardBackImage *ebiten.Image
	// CardShadowImage applies to all cards so is kept globally as an optimization
	CardShadowImage *ebiten.Image
	// CardMovableImage applies to all cards so is kept globally as an optimization
	// CardMovableImage *ebiten.Image
)

// TheGSM provides global access to the game state manager
var TheGSM *GameStateManager = &GameStateManager{}

// ThePreferences holds serialized game progress data
var ThePreferences = &Preferences{Game: "Solitaire", Variant: "Klondike", HighlightMovable: true, PowerMoves: true, SingleTap: true, CardBackPattern: "FlowerBlue", CardBackColor: "CornflowerBlue"}

// TheStatistics holds statistics for all variants
var TheStatistics *Statistics

// TheBaize points to the Baize, so that main can see it
var TheBaize *Baize

// The UI points to the singleton User Interface object
var TheUI *ui.UI

var TheCIP CardImageProvider

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}

	TheStatistics = NewStatistics()

	TheGSM.Switch(NewSplash())

	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	state := TheGSM.Get()
	return state.Layout(outsideWidth, outsideHeight)
}

// Update updates the current game state.
func (g *Game) Update() error {
	state := TheGSM.Get()
	if err := state.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	state := TheGSM.Get()
	state.Draw(screen)
}
