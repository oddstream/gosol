package sol

import (
	"syscall/js"
)

func (b *Baize) Wikipedia() {
	js.Global().Get("window").Call("open", b.vinfo.wikipedia)
}
