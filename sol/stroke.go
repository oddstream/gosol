package sol

import (
	"time"

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

	startTime time.Time

	released bool
	tapped   bool

	// draggingObject represents a object (like a tile) that is being dragged.
	draggingObject interface{}
}

// NewStroke creates a new Stroke object
func NewStroke(source StrokeSource) *Stroke {
	cx, cy := source.Position()
	return &Stroke{
		source:    source,
		initX:     cx,
		initY:     cy,
		currX:     cx,
		currY:     cy,
		startTime: time.Now(),
	}
}

// Update is called once per frame and updates the Stroke object
func (s *Stroke) Update() {
	if s.released {
		return
	}

	s.currX, s.currY = s.source.Position()

	if s.source.IsJustReleased() {
		s.released = true
		elapsed := time.Since(s.startTime) / 1000 / 1000 // convert nano- to milli- seconds
		s.tapped = elapsed < 150
	}
}

// IsReleased returns true if ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// IsTapped returns true if ...
func (s *Stroke) IsTapped() bool {
	return s.tapped
}

// Position returns the x,y position of the cursor
func (s *Stroke) Position() (int, int) {
	return s.currX, s.currY
}

// PositionDiff returns the x,y difference between the start of the stroke and the stoke's current position
func (s *Stroke) PositionDiff() (int, int) {
	return s.currX - s.initX, s.currY - s.initY
}

// DraggingObject returns a reference to the object currently being dragged
func (s *Stroke) DraggingObject() interface{} {
	return s.draggingObject
}

// SetDraggingObject sets the object currently being dragged
func (s *Stroke) SetDraggingObject(object interface{}) {
	s.draggingObject = object
}
