package main

import (
	"log"

	"pomodoro-do-ben/config"
	"pomodoro-do-ben/gui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gui.Show(cfg)
}
