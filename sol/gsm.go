// Copyright ©️ 2021 oddstream.games

package maze

import "github.com/hajimehoshi/ebiten/v2"

// GameState interface defines the API for each game state
// each separate game state (eg Splash, Menu, Grid, GameOver &c) must implement these
type GameState interface {
	Layout(int, int) (int, int)
	Update() error
	Draw(*ebiten.Image)
}

// GameStateManager does what it says on the tin
type GameStateManager struct {
	// TODO implement a stack with Push(), Pop() methods
	currentState GameState
}

// Switch changes to a different GameState
func (gsm *GameStateManager) Switch(state GameState) {
	gsm.currentState = state
}

// Get returns the current GameState
func (gsm *GameStateManager) Get() GameState {
	return gsm.currentState
}
