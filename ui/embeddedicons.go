package ui

//lint:file-ignore U1000,ST1003 Ignore unused code and underscores in generated code

import (
	"bytes"
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"image"
	"log"
)

//go:embed icons/bookmark.png
var bookmarkIconBytes []byte

//go:embed icons/bookmark_add.png
var bookmark_addIconBytes []byte

//go:embed icons/check_box.png
var check_boxIconBytes []byte

//go:embed icons/check_box_outline_blank.png
var check_box_outline_blankIconBytes []byte

//go:embed icons/close.png
var closeIconBytes []byte

//go:embed icons/done.png
var doneIconBytes []byte

//go:embed icons/done_all.png
var done_allIconBytes []byte

//go:embed icons/info.png
var infoIconBytes []byte

//go:embed icons/list.png
var listIconBytes []byte

//go:embed icons/menu.png
var menuIconBytes []byte

//go:embed icons/radio_button_checked.png
var radio_button_checkedIconBytes []byte

//go:embed icons/radio_button_unchecked.png
var radio_button_uncheckedIconBytes []byte

//go:embed icons/restore.png
var restoreIconBytes []byte

//go:embed icons/search.png
var searchIconBytes []byte

//go:embed icons/settings.png
var settingsIconBytes []byte

//go:embed icons/star.png
var starIconBytes []byte

//go:embed icons/undo.png
var undoIconBytes []byte

//go:embed icons/lightbulb.png
var lightbulbIconBytes []byte

//go:embed icons/poll.png
var pollIconBytes []byte

//go:embed icons/speed.png
var speedIconBytes []byte

//go:embed icons/wikipedia.png
var wikipediaIconBytes []byte

var IconMap = map[string]image.Image{}

func decode(name string, variable []byte) {
	img, _, err := image.Decode(bytes.NewReader(variable))
	if err != nil {
		log.Panic(err)
	}
	IconMap[name] = img
}

// LoadIconMapFromEmbedded loads icons from go:embed vars
func LoadIconMapFromEmbedded() {
	decode("bookmark", bookmarkIconBytes)
	decode("bookmark_add", bookmark_addIconBytes)
	decode("check_box", check_boxIconBytes)
	decode("check_box_outline_blank", check_box_outline_blankIconBytes)
	decode("close", closeIconBytes)
	decode("done", doneIconBytes)
	decode("done_all", done_allIconBytes)
	decode("info", infoIconBytes)
	decode("list", listIconBytes)
	decode("menu", menuIconBytes)
	decode("radio_button_checked", radio_button_checkedIconBytes)
	decode("radio_button_unchecked", radio_button_uncheckedIconBytes)
	decode("restore", restoreIconBytes)
	decode("search", searchIconBytes)
	decode("settings", settingsIconBytes)
	decode("star", starIconBytes)
	decode("undo", undoIconBytes)
	decode("lightbulb", lightbulbIconBytes)
	decode("poll", pollIconBytes)
	decode("speed", speedIconBytes)
	decode("wikipedia", wikipediaIconBytes)
}
