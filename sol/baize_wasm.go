package sol

import (
	"syscall/js"
)

func (b *Baize) Wikipedia() {
	url := b.script.Wikipedia()
	js.Global().Get("window").Call("open", url)
}
