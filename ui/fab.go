package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// ActionButton is a Widget that sits in the Floating Action Bar.
// It's an icon on a round background, linked to a keyboard-shortcut command.
type ActionButton struct {
	WidgetBase
	iconName string
	key      ebiten.Key
}

func (ab *ActionButton) createImg() *ebiten.Image {
	// WidgetBase doesn't have a default createImg
	if ab.width == 0 || ab.height == 0 {
		return nil
	}
	dc := gg.NewContext(ab.width, ab.height)
	dc.SetColor(BackgroundColor)
	dc.DrawCircle(float64(ab.width/2), float64(ab.height/2), float64(ab.height/2))
	dc.Fill()
	dc.Stroke()
	dc.SetColor(ForegroundColor)
	dc.DrawImageAnchored(IconMap[ab.iconName], ab.width/2, ab.height/2, 0.5, 0.5)
	return ebiten.NewImageFromImage(dc.Image())
}

func NewActionButton(parent Containery, id string, iconName string, key ebiten.Key) *ActionButton {
	ab := &ActionButton{WidgetBase: WidgetBase{parent: parent, id: id, x: 0, y: 0, width: ActionButtonSize, height: ActionButtonSize}, iconName: iconName, key: key}
	ab.Activate()
	return ab
}

func (ab *ActionButton) Tapped() {
	if ab.disabled {
		return
	}
	cmdFn(ab.key)
}

// Activate this widget. Silly, because ActionButtons are never deactivated.
func (ab *ActionButton) Activate() {
	ab.disabled = false
	ab.img = ab.createImg()
}

// Deactivate this widget. Silly, because ActionButtons are never deactivated.
func (ab *ActionButton) Deactivate() {
	ab.disabled = true
	ab.img = ab.createImg()
}

//

type FAB struct {
	WindowBase
}

func NewFAB() *FAB {
	fb := &FAB{WindowBase: WindowBase{x: 0, y: 0, width: FABWidth, height: FABHeight}}
	fb.img = fb.createImg(TransparentColor)
	// no widgets yet
	return fb
}

// StartDrag notifies this container that dragging has started
func (f *FAB) StartDrag() {
	// println("FAB start drag, adding widgets")
}

// DragBy this container
func (f *FAB) DragBy(dx, dy int) {
	// you can't drag a bar
}

// StopDrag notifies the container that dragging has been stopped
func (f *FAB) StopDrag() {
}

// CancelDrag notifies the container that dragging has been cancelled
func (f *FAB) CancelDrag() {
}

func (f *FAB) Tapped() {
}

// Show the bar
func (f *FAB) Show() {
}

// Hide the bar
func (f *FAB) Hide() {
	f.widgets = f.widgets[:0]
}

// Visible is the bar
func (f *FAB) Visible() bool {
	return len(f.widgets) > 0
}

func (f *FAB) LayoutWidgets() {
	var x int = 0
	var y int = f.height - ActionButtonSize
	// ActionButtons are stacked upwards from the bottom of the FAB
	for _, w := range f.widgets {
		w.SetPosition(x, y)
		y -= ActionButtonSize
	}
}

func (f *FAB) Update() {
}

// Layout implements Ebiten's Layout
func (f *FAB) Layout(outsideWidth, outsideHeight int) (int, int) {
	// override BarBase.Layout to get position near bottom right of screen
	f.x = outsideWidth - f.width - (f.width / 2)
	f.y = outsideHeight - f.height - (ActionButtonSize / 2) - StatusbarHeight // statusbar is 24 high
	// println("FAB.Layout() Window=", outsideWidth, outsideHeight, "FAB=", fb.x, fb.y)
	return outsideWidth, outsideHeight
}

//

func (u *UI) AddButtonToFAB(iconName string, key ebiten.Key) {
	u.fab.widgets = append(u.fab.widgets, NewActionButton(u.fab, "", iconName, key))
	u.fab.LayoutWidgets()
}

// HideFAB doesn't actually hide the FAB; it removes all the widegets. The FAB is always shown, albeit transparently.
func (u *UI) HideFAB() {
	u.fab.Hide()
}
