package ui

import (
	"syscall/js"
)

func OpenBrowserWindow(url string) {
	js.Global().Get("window").Call("open", url)
}
