package ui

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Toast represents a simple popup window that disappears after a few seconds
type Toast struct {
	img       *ebiten.Image
	ticksLeft int
}

// ToastManager manages the queue of toasts so that many may appear on screen at once
type ToastManager struct {
	toasts []*Toast
}

// Toast creates a new toast message an adds it to the queue of messages
func (u *UI) Toast(message string) {

	dc := gg.NewContext(8, 8)
	dc.SetFontFace(u.fontFace)
	w, h := dc.MeasureString(message)

	w += 20
	h += 20

	dc = gg.NewContext(int(w), int(h))
	dc.SetRGBA(0, 0, 0, 0.5)
	dc.DrawRectangle(0, 0, w, h)
	dc.Fill()
	dc.Stroke()

	dc.SetFontFace(u.fontFace)
	dc.SetRGBA(0.9, 0.9, 0.9, 1)
	dc.DrawStringAnchored(message, w/2, h/2, 0.5, 0.5)
	dc.Stroke()

	t := &Toast{}
	t.img = ebiten.NewImageFromImage(dc.Image())
	t.ticksLeft = int(ebiten.CurrentTPS()) * 5

	u.toastManager.Add(t)
	println("toast:", message, "(", int(w), ",", int(h), ")")
}

// Add a new toast to the queued toasts
func (tm *ToastManager) Add(t *Toast) {
	tm.toasts = append([]*Toast{t}, tm.toasts...)
}

// Update the queue of toasts
func (tm *ToastManager) Update() {
	if len(tm.toasts) == 0 {
		return
	}
	for _, t := range tm.toasts {
		t.ticksLeft--
	}
	t := tm.toasts[0]
	if t.ticksLeft < 0 {
		tm.toasts = tm.toasts[1:]
	}
}

// Draw the toasts
func (tm *ToastManager) Draw(screen *ebiten.Image) {

	if len(tm.toasts) == 0 {
		return
	}
	sx, sy := screen.Size()
	var tx, ty int
	ty = sy - 10
	for _, t := range tm.toasts {
		w, h := t.img.Size()
		tx = (sx - w) / 2
		ty = ty - h - 10
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tx), float64(ty))
		screen.DrawImage(t.img, op)
	}
}
