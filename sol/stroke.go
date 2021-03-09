package sol

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

	// draggingObject represents a object (like a tile) that is being dragged.
	// draggingObject interface{}
	draggedCard *Card

	observer sync.Map
}

// StrokeEvent is sent to observers when stroke moves or ends
type StrokeEvent struct {
	Event  string
	Stroke *Stroke
	X, Y   int
}

// StartStroke returns a pointer to a new Stroke if one is just starting
func StartStroke(observer Observer) *Stroke {
	// Stroke always starts immediately otherwise weird lag (tap will cancel drag)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		var source StrokeSource = &MouseStrokeSource{}
		x, y := source.Position()
		s := &Stroke{
			source: source,
			initX:  x,
			initY:  y,
			currX:  x,
			currY:  y,
		}
		s.Add(observer)
		s.Notify(StrokeEvent{Event: "start", Stroke: s, X: s.initX, Y: s.initY})
		return s
	}
	return nil
}

// Update is called once per frame and updates the Stroke object
func (s *Stroke) Update() {

	if s.released || s.cancelled {
		return
	}

	x, y := s.source.Position()
	if s.currX != x || s.currY != y {
		s.currX, s.currY = x, y
		s.Notify(StrokeEvent{Event: "move", Stroke: s, X: s.currX, Y: s.currY})
	}

	if s.source.IsJustReleased() {
		s.released = true
		s.Notify(StrokeEvent{Event: "end", Stroke: s, X: s.currX, Y: s.currY})
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

// // DraggingObject returns a reference to the object currently being dragged
// func (s *Stroke) DraggingObject() interface{} {
// 	return s.draggingObject
// }

// // SetDraggingObject sets the object currently being dragged
// func (s *Stroke) SetDraggingObject(object interface{}) {
// 	s.draggingObject = object
// }

// DraggedCard returns the card being dragged
func (s *Stroke) DraggedCard() *Card {
	return s.draggedCard
}

// SetDraggedCard sets the object currently being dragged
func (s *Stroke) SetDraggedCard(c *Card) {
	s.draggedCard = c
}

// Add this observer to the list
func (s *Stroke) Add(observer Observer) {
	s.observer.Store(observer, struct{}{})
}

// Remove this observer from the list
func (s *Stroke) Remove(observer Observer) {
	s.observer.Delete(observer)
}

// Notify observers that an event has happened
func (s *Stroke) Notify(event interface{}) {
	s.observer.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		key.(Observer).NotifyCallback(event)
		return true
	})
}
