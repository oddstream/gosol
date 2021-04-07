package input

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// init X,Y represents the position when dragging starts.
	initX, initY int

	// current X,Y represents the current position
	currX, currY int

	// startTime time.Time

	// starting  bool
	released  bool
	cancelled bool

	// draggedObject represents a object (eg a Card) that is being dragged
	// can't have a valid stroke with an object that is being dragged
	draggedObject interface{}

	observers sync.Map // sync.Map is not type safe, it is similar to a map[interface{}]interface{}
}

// StrokeEvent is sent to observers when stroke moves or ends
type StrokeEvent struct {
	Event  string
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
	ids := inpututil.JustPressedTouchIDs()
	if len(ids) > 0 {
		s = NewStroke(&TouchStrokeSource{ID: ids[0]})
	}
	if s != nil {
		s.Add(observer)
		s.Notify(StrokeEvent{Event: "start", Stroke: s, X: s.initX, Y: s.initY})
	}
	return s
}

// Update is called once per frame and updates the Stroke object
func (s *Stroke) Update() {

	if s.released || s.cancelled {
		return
	}

	x, y := s.source.Position()
	if s.currX != x || s.currY != y {
		s.currX, s.currY = x, y
		s.Notify(StrokeEvent{Event: "move", Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
	}

	if s.source.IsJustReleased() {
		s.released = true
		s.Notify(StrokeEvent{Event: "stop", Stroke: s, Object: s.draggedObject, X: s.currX, Y: s.currY})
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
	s.observers.Store(observer, struct{}{})
}

// Remove this observer from the list
func (s *Stroke) Remove(observer Observer) {
	s.observers.Delete(observer)
}

// Notify observers that an event has happened
func (s *Stroke) Notify(event interface{}) {
	// cannot range over a value of type sync.Map
	s.observers.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		key.(Observer).NotifyCallback(event)
		return true
	})
}
