// Copyright ©️ 2021 oddstream.games

package sol

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Game represents a game state.
type Game struct{}

var (
	// DebugMode is a boolean set by command line flag -debug
	DebugMode bool = false
	// WindowWidth of main window in pixels
	WindowWidth int
	// WindowHeight of main window in pixels
	WindowHeight int
)

// GSM provides global access to the game state manager
var GSM *GameStateManager = &GameStateManager{}

// CTQ provides global access to the Card Transition Queue
var CTQ *CardTransitionQueue = &CardTransitionQueue{}

// TheAcmeFonts provides access to small, normal, large, huge Acme fonts
var TheAcmeFonts *AcmeFonts = NewAcmeFonts()

// TheUserData holds serialized game progress data
var TheUserData = &UserData{Copyright: "Copyright ©️ 2021 oddstream.games", Game: "Solitaire", Variant: "Limited", CardBack: "Roses"}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}
	// println("LevelData", len(LevelData))

	//	TheUserData.Load()

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
