package ui

import (
	"image"
	"log"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
	"oddstream.games/gosol/util"
)

// Checkbox is a button that displays a single rune
type Checkbox struct {
	WidgetBase
	checked bool
	text    string
}

func (w *Checkbox) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)

	var iconName string
	if w.checked {
		iconName = "check_box"
	} else {
		iconName = "check_box_outline_blank"
	}
	// same as NavItem
	img, ok := IconMap[iconName]
	if !ok || img == nil {
		log.Fatal(iconName, " not in icon map")
	}
	dc.DrawImageAnchored(img, 0, w.height/2, 0, 0.5)

	dc.SetRGBA(1, 1, 1, 1)
	dc.SetFontFace(schriftbank.RobotoMedium24)
	dc.DrawString(w.text, float64(48), float64(w.height)*0.7)

	// uncomment this to show the area we expect the text to occupy
	// dc.DrawLine(0, float64(0), float64(w.width), float64(0))
	// dc.DrawLine(0, float64(w.height), float64(w.width), float64(w.height))
	// dc.DrawLine(0, float64(0), float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// NewCheckbox creates a new Checkbox
func NewCheckbox(parent Container, input *input.Input, text string, checked bool) *Checkbox {
	width, _ := parent.Size()
	w := &Checkbox{
		WidgetBase: WidgetBase{parent: parent, input: input, img: nil, x: 0, y: 0, width: width, height: 48},
		text:       text, checked: checked}
	w.Activate()
	return w
}

// Activate tells the input we need notifications
func (w *Checkbox) Activate() {
	w.disabled = false
	w.img = w.createImg()
	w.input.Add(w)
}

// Deactivate tells the input we no longer need notofications
func (w *Checkbox) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
	w.input.Remove(w)
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (w *Checkbox) NotifyCallback(event interface{}) {
	if w.disabled {
		return
	}
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		// println("Checkbox image.Point", v.X, v.Y)
		if util.InRect(v.X, v.Y, w.OffsetRect) {
			w.checked = !w.checked
			w.img = w.createImg()
			w.input.Notify(ChangeRequest{ChangeRequested: w.text, Data: strconv.FormatBool(w.checked)})

		}
	}
}
