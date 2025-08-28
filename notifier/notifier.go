
package notifier

import (
	"log"
	"os/exec"
	"pomodoro-do-ben/i18n"
)

func Notify(title, message string) {
	cmd := exec.Command("notify-send", i18n.T(title), i18n.T(message))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error sending notification:", err)
		log.Println("notify-send output:", string(output))
	}
}
