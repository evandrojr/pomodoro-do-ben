package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"pomodoro-do-ben/config"
	"pomodoro-do-ben/gui"
	"pomodoro-do-ben/i18n"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	myApp := app.New()
	myWindow := myApp.NewWindow(i18n.T("bens_pomodoro"))

	icon, err := fyne.LoadResourceFromPath("pomodoro-do-ben.png")
	if err != nil {
		fmt.Println("Error loading icon:", err)
	}
	myApp.SetIcon(icon)

	gui.Show(cfg, myWindow)
}