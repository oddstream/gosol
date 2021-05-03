package ui

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/input"
	"oddstream.games/gosol/schriftbank"
)

// TextUrl is a widget that displays a clickable url
type TextUrl struct {
	WidgetBase
	text       string
	url        string
	lines      []string
	lineHeight int
}

func (w *TextUrl) createImg() *ebiten.Image {
	dc := gg.NewContext(w.width, w.height)

	dc.SetRGBA(0.5, 0.5, 1, 1)
	// nota bene - text is drawn with y as a baseline
	dc.SetFontFace(schriftbank.RobotoRegular14)
	y := w.lineHeight
	for _, str := range w.lines {
		dc.DrawString(str, 0, float64(y-8)) // move up a little to stop descenders being clipped on last line
		y += w.lineHeight
	}
	// uncomment this line to visualize text box
	// dc.DrawLine(0, 0, float64(w.width), float64(w.height))
	// dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

func (w *TextUrl) calcHeights() {
	dc := gg.NewContext(w.width, 48)
	dc.SetFontFace(schriftbank.RobotoRegular14)
	// MeasureString says this text, requested to be 48 high, is 14 high
	w.lines = dc.WordWrap(w.text, float64(w.width-48)) // 24 padding left and right
	w.lineHeight = 24
	w.height = w.lineHeight * len(w.lines) // + w.lineHeight
}

// NewTextUrl creates a new TextUrl widget
func NewTextUrl(parent Container, text string) *TextUrl {
	width, _ := parent.Size()
	// widget x, y will be set by LayoutWidgets
	// widget height will be set when wordwrapping in createImg
	w := &TextUrl{
		WidgetBase: WidgetBase{parent: parent, img: nil, width: width},
		text:       "Wikipedia",
		url:        text}
	w.calcHeights()
	return w
}

// Activate tells the input we need notifications
func (w *TextUrl) Activate() {
	w.disabled = false
	w.img = w.createImg()
	// w.input.Add(w)
}

// Deactivate tells the input we no longer need notofications
func (w *TextUrl) Deactivate() {
	w.disabled = true
	w.img = w.createImg()
	// w.input.Remove(w)
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}

// NotifyCallback is called by the Subject (Input/Stroke) when something interesting happens
func (w *TextUrl) NotifyCallback(v input.StrokeEvent) {
	if w.disabled {
		return
	}
	switch v.Event {
	case "tap":
		openBrowser(w.url)
	}
}

// Update the state of this widget
func (w *TextUrl) Update() {
}
