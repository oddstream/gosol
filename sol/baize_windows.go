package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", b.vinfo.wikipedia).Start()
	if err != nil {
		log.Println(err)
	}
}
