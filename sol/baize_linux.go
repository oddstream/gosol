package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	var cmd *exec.Cmd = exec.Command("xdg-open", b.vinfo.wikipedia)
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}
