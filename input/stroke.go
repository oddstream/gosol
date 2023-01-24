package input

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"oddstream.games/gosol/util"
)

type (
	// Observable https://gist.github.com/patrickmn/1549985
	Observable interface {
		Add(Observer)
		Notify(StrokeEvent)
		Remove(Observer)
	}

	// Observer https://gist.github.com/patrickmn/1549985
	Observer interface {
		NotifyCallback(StrokeEvent)
	}
)

// The following is taken from https://ebiten.org/examples/drag.html

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

// Position returns the x,y cordinates of the cursor position
func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

// IsJustReleased returns true if the left mouse button was released in the current frame
func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

// Position returns the x,y cordinates of the cursor position
func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

// IsJustReleased returns true if the first touch was released in the current frame
func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

type WheelStrokeSource struct{}

func (w *WheelStrokeSource) Position() (int, int) {
	x, y := ebiten.Wheel()
	return int(x), int(y)
}

func (w *WheelStrokeSource) IsJustReleased() bool {
	_, y := ebiten.Wheel()
	return y == 0.0
}

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source        StrokeSource
	initX, initY  int
	currX, currY  int
	released      bool
	cancelled     bool
	draggedObject interface{} // object (eg a []*Card) that is being dragged
	observers     sync.Map    // sync.Map is not type safe, it is similar to a map[interface{}]interface{}
}

type EventType int

const (
	Start EventType = iota + 1
	Move
	Tap
	Stop
	Cancel
)

// StrokeEvent is sent to observers when stroke moves or ends
type StrokeEvent struct {
	Event  EventType
	Stroke *Stroke
	Object interface{}
	X, Y   int
}

// NewStroke create a new Stroke object
func NewStroke(source StrokeSource) *Stroke {
	x, y := source.Position()
	return &Stroke{
		source: source,
		initX:  x,
		initY:  y,
		currX:  x,
		currY:  y,
	}
}

// StartStroke returns a pointer to a new Stroke if one is just starting
func StartStroke(observer Observer) *Stroke {
	var s *Stroke
	// Stroke always starts immediately otherwise weird lag (tap will cancel drag)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s = NewStroke(&MouseStrokeSource{})
	}

	if s == nil {
		var ids []ebiten.TouchID = []ebiten.TouchID{}
		ids = inpututil.AppendJustPressedTouchIDs(ids)
		if len(ids) > 0 {
			// println(len(ids), "touch IDs, first is", ids[0])
			s = NewStroke(&TouchStrokeSource{ID: ids[0]})
		}
	}

	if s == nil {
		_, y := ebiten.Wheel()
		if y != 0.0 {
			s = NewStroke(&WheelStrokeSource{})
		}
	}

	if s != nil {
		s.Add(observer)
		s.Notify(StrokeEvent{Event: Start, Stroke: s, X: s.initX, Y: s.initY})
	}
	return s
}

// Update is called once per frame and updates the Stroke object
func (s *Stroke) Update() {

	if s.released || s.cancelled {
		return
	}

	// test release before testing move
	// to make the "dropping while moving" problem better
	if s.source.IsJustReleased() {
		s.released = true
		if util.Abs(s.initX-s.currX) < 4 && util.Abs(s.initY-s.currY) < 4 {
			s.Notify(StrokeEvent{Event: Cancel, Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
			s.Notify(StrokeEvent{Event: Tap, Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
		} else {
			s.Notify(StrokeEvent{Event: Stop, Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
		}
	} else {
		x, y := s.source.Position()
		if s.currX != x || s.currY != y {
			s.currX, s.currY = x, y
			s.Notify(StrokeEvent{Event: Move, Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
		}
	}

}

// Cancel this stroke; observer is not interested
func (s *Stroke) Cancel() {
	s.cancelled = true
}

// IsReleased returns true if ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// IsCancelled returns true if ...
func (s *Stroke) IsCancelled() bool {
	return s.cancelled
}

// Position returns the x,y position of the cursor
func (s *Stroke) Position() (int, int) {
	return s.currX, s.currY
}

// PositionDiff returns the x,y difference between the start of the stroke and the stoke's current position
func (s *Stroke) PositionDiff() (int, int) {
	return s.currX - s.initX, s.currY - s.initY
}

// DraggedObject returns a reference to the object currently being dragged
func (s *Stroke) DraggedObject() interface{} {
	return s.draggedObject
}

// SetDraggedObject sets the object currently being dragged
func (s *Stroke) SetDraggedObject(object interface{}) {
	s.draggedObject = object
}

// Add this observer to the list
func (s *Stroke) Add(observer Observer) {
	// fmt.Printf("Stroke Add() %T\n", observer)
	s.observers.Store(observer, struct{}{})
	// count := 0
	// s.observers.Range(func(key, value interface{}) bool {
	// 	if key == nil {
	// 		return false
	// 	}
	// 	count++
	// 	return true
	// })
	// println(count, "observers")
}

// Remove this observer from the list
func (s *Stroke) Remove(observer Observer) {
	s.observers.Delete(observer)
}

// Notify observers that an event has happened
func (s *Stroke) Notify(event StrokeEvent) {
	// cannot range over a value of type sync.Map
	s.observers.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		// fmt.Printf("Notifying a %T\n", key)
		key.(Observer).NotifyCallback(event)
		return true
	})
}
