// Package input provides basic tap, stroke and key press input
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
	observers          sync.Map // sync.Map is not type safe, it is similar to a map[interface{}]interface{}
	timePressed        time.Time
	xPressed, yPressed int
	id                 ebiten.TouchID
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

	// TODO refactor this mess

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		i.xPressed, i.yPressed = ebiten.CursorPosition()
		i.timePressed = time.Now()
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		elapsed := time.Since(i.timePressed) / 1000 / 1000 // convert nano- to milli- seconds
		xNow, yNow := ebiten.CursorPosition()
		// distance := util.DistanceInt(i.xPressed, i.yPressed, xNow, yNow)
		// can't use distance < n because card will be animating
		if elapsed < 200 || (i.xPressed == xNow && i.yPressed == yNow) {
			i.Notify(image.Point{X: xNow, Y: yNow})
		}
	} else {
		ids := inpututil.JustPressedTouchIDs()
		if len(ids) > 0 {
			i.id = ids[0]
			i.xPressed, i.yPressed = ebiten.CursorPosition()
			i.timePressed = time.Now()
		}
		if inpututil.IsTouchJustReleased(i.id) {
			elapsed := time.Since(i.timePressed) / 1000 / 1000 // convert nano- to milli- seconds
			xNow, yNow := ebiten.CursorPosition()
			if elapsed < 200 || (i.xPressed == xNow && i.yPressed == yNow) {
				i.Notify(image.Point{X: xNow, Y: yNow})
			}
		}
	}

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustReleased(k) {
			i.Notify(k)
		}
	}

}
