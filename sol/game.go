// Package sol provides a polymorphic solitaire engine
package sol

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

// Game represents a game state
type Game struct {
}

var (
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// NoGameLoad is a boolean set by command line flag -noload
	NoGameLoad = false
	// NoGameSave is a boolean set by command line flag -nosave
	NoGameSave = false
	// WindowWidth of main window in pixels
	WindowWidth int
	// WindowHeight of main window in pixels
	WindowHeight int
	// CardWidth CardHeight of cards
	CardWidth int = 71
	// CardHeight of cards
	CardHeight int = 96
	// PileMarginX the horzontal gaps between piles
	PileMarginX int = 10
	// PileMarginY the vertical gap between piles
	PileMarginY int = 10
	// TopMargin the gap between top pile and top of baize
	TopMargin int = 48 + CardHeight/3
)

// GSM provides global access to the game state manager
var GSM *GameStateManager = &GameStateManager{}

// CTQ provides global access to the Card Transition Queue
//var CTQ *CardTransitionQueue = &CardTransitionQueue{}

// TheUserData holds serialized game progress data
var TheUserData = &UserData{Copyright: "Copyright ©️ 2021 oddstream.games", Game: "Solitaire", Variant: "Klondike", CardBack: "FlowerBlue", CardStyle: "retro", BackColor: "CornflowerBlue"}

// TheStatistics holds statistics for all variants
var TheStatistics *Statistics

// TheBaize points to the Baize, so that main can see it
var TheBaize *Baize

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}

	/*
		71 x 96 = 1:1.352 (Microsoft retro)
		64 x 89 = 1:1.390 (official poker size)
		90 x 130 = 1:1.444 (nice looking scalable)
		89 x 137 = 1:1.539 (measured real card)
		57 x 89 = 1:1.561 (official bridge size)
	*/
	switch TheUserData.CardStyle {
	case "", "default", "scalable":
		CardHeight = int(math.Ceil(float64(CardWidth) * 1.444))
	case "poker":
		CardHeight = int(math.Ceil(float64(CardWidth) * 1.39))
	case "bridge":
		CardHeight = int(math.Ceil(float64(CardWidth) * 1.561))
	}
	log.Printf("card size %s %dx%d", TheUserData.CardStyle, CardWidth, CardHeight)
	schriftbank.MakeCardFonts(CardWidth) // CardWidth/Height have now been set
	TheStatistics = NewStatistics()

	GSM.Switch(NewSplash())

	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	WindowWidth = outsideWidth
	WindowHeight = outsideHeight
	state := GSM.Get()
	return state.Layout(outsideWidth, outsideHeight)
}

// Update updates the current game state.
func (g *Game) Update() error {
	state := GSM.Get()
	if err := state.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	state := GSM.Get()
	state.Draw(screen)
}
