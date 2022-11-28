package ui

import (
	"log"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// RadioButton is a button that displays a single rune
type RadioButton struct {
	WidgetBase
	checked bool
	text    string
}

func (w *RadioButton) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)

	var iconName string
	if w.checked {
		iconName = "radio_button_checked"
	} else {
		iconName = "radio_button_unchecked"
	}
	// same as NavItem
	img, ok := IconMap[iconName]
	if !ok || img == nil {
		log.Fatal(iconName, " not in icon map")
	}
	dc.DrawImage(img, 0, w.height/4)

	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.RobotoMedium24)
	dc.DrawString(w.text, float64(48), float64(w.height)*0.8)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(w.width), float64(0))
	// dc.DrawLine(0, float64(w.height), float64(w.width), float64(w.height))
	// dc.DrawLine(0, float64(0), float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewRadioButton creates a new RadioButton
func NewRadioButton(parent Containery, id string, text string, checked bool) *RadioButton {
	width, _ := parent.Size()
	w := &RadioButton{
		WidgetBase: WidgetBase{parent: parent, id: id, img: nil, x: 0, y: 0, width: width, height: 48},
		text:       text, checked: checked}
	w.Activate()
	return w
}

// Activate tells the input we need notifications
func (w *RadioButton) Activate() {
	w.disabled = false
	w.img = w.createImg()
}

// Deactivate tells the input we no longer need notofications
func (w *RadioButton) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (w *RadioButton) NotifyCallback(v input.StrokeEvent) {
	if w.disabled {
		return
	}
	switch v.Event {
	case input.Tap:
		if util.InRect(v.X, v.Y, w.OffsetRect) {
			w.checked = !w.checked
			w.img = w.createImg()
			cmdFn(ChangeRequest{ChangeRequested: w.text, Data: strconv.FormatBool(w.checked)})
		}
	}
}
