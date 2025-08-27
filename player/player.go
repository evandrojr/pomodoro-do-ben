
package player

import (
	"log"
	"os/exec"
)

func Play(file string) {
	cmd := exec.Command("mpv", "--no-video", file)
	err := cmd.Run()
	if err != nil {
		log.Println("Error playing audio:", err)
	}
}
