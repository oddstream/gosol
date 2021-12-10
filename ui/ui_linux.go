package ui

import (
	"log"
	"os/exec"
)

func OpenBrowserWindow(url string) {
	var cmd *exec.Cmd = exec.Command("xdg-open", url)
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}
