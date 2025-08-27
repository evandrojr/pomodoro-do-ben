
package notifier

import (
	"log"
	"os/exec"
)

func Notify(title, message string) {
	cmd := exec.Command("notify-send", title, message)
	err := cmd.Run()
	if err != nil {
		log.Println("Error sending notification:", err)
	}
}
