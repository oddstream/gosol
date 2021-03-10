// Copyright ©️ 2021 oddstream.games

package input

import (
	"image"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input records state of mouse and touch, Subject in Observer pattern
type Input struct {
	// pressed        map[ebiten.Key]struct{} // an empty and useless type
	observers   sync.Map
	timePressed time.Time
}

// NewInput Input object constructor
func NewInput() *Input {
	// no fields to initialize, so use the built-in new()
	return new(Input)
}

// Add this observer to the list
func (i *Input) Add(observer Observer) {
	i.observers.Store(observer, struct{}{})
}

// Remove this observer from the list
func (i *Input) Remove(observer Observer) {
	i.observers.Delete(observer)
}

// Notify observers that an event has happened
func (i *Input) Notify(event interface{}) {
	i.observers.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		key.(Observer).NotifyCallback(event)
		return true
	})
}

// Update the state of the Input object
func (i *Input) Update() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		i.timePressed = time.Now()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		elapsed := time.Since(i.timePressed) / 1000 / 1000 // convert nano- to milli- seconds
		if elapsed < 150 {
			x, y := ebiten.CursorPosition()
			i.Notify(image.Point{X: x, Y: y})
		}
	}

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustReleased(k) {
			i.Notify(k)
		}
	}

}
