package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	url := b.script.Wikipedia()
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		log.Println(err)
	}
}
