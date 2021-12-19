package sol

import (
	"log"
	"os/exec"
)

func (b *Baize) Wikipedia() {
	url := b.script.Wikipedia()
	var cmd *exec.Cmd = exec.Command("xdg-open", url)
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}
