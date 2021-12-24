package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", b.script.Info().wikipedia).Start()
	if err != nil {
		log.Println(err)
	}
}
