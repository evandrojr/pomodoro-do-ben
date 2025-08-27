package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"pomodoro-do-ben/config"
	"pomodoro-do-ben/notifier"
	"pomodoro-do-ben/player"
	"pomodoro-do-ben/pomo"
)

func Show(cfg *config.Config) {
	myApp := app.New()
	myWindow := myApp.NewWindow(config.AppName)

	timer := pomo.NewTimer(cfg)

	timerStr := binding.NewString()
	timerStr.Set(formatTime(timer.RemainingTime))

	timerLabel := widget.NewLabelWithData(timerStr)
	timerLabel.Alignment = fyne.TextAlignCenter
	timerLabel.TextStyle = fyne.TextStyle{Bold: true}

	sessionLabel := widget.NewLabel("Pomodoro do Ben")
	sessionLabel.Alignment = fyne.TextAlignCenter

	startButton := widget.NewButtonWithIcon("▶️ Start", theme.MediaPlayIcon(), func() {
		if !timer.IsRunning {
			if isInactive(cfg) {
				notifier.Notify("Pomodoro", "Timer is inactive during this period.")
				return
			}
			timer.Start()
			player.Play(getMediaPath("focar/f1.aac"))
			notifier.Notify("Pomodoro", "Timer started!")
		}
	})

	pauseButton := widget.NewButtonWithIcon("⏸️ Pause", theme.MediaPauseIcon(), func() {
		if timer.IsRunning {
			timer.Stop()
			notifier.Notify("Pomodoro", "Timer paused!")
		}
	})

	resetButton := widget.NewButtonWithIcon("🔄 Reset", theme.MediaReplayIcon(), func() {
		timer.Reset()
		timerStr.Set(formatTime(timer.RemainingTime))
		notifier.Notify("Pomodoro", "Timer reset!")
	})

	go func() {
		for range timer.Updates {
			timerStr.Set(formatTime(timer.RemainingTime))
			updateTitle(myWindow, timer)
			updateSessionLabel(sessionLabel, timer)
			myWindow.Canvas().Refresh(myWindow.Content())
			if timer.RemainingTime <= 0 {
				timer.NextState()
				player.Play(getMediaPath("meditar/m1.aac"))
				notifier.Notify("Pomodoro", "Time for a break!")
			}
		}
	}()

	buttons := container.NewHBox(layout.NewSpacer(), startButton, pauseButton, resetButton, layout.NewSpacer())
	pomodoroTab := container.NewVBox(timerLabel, sessionLabel, buttons)

	startOnLaunchBinding := binding.NewBool()
	startOnLaunchBinding.Set(cfg.StartOnLaunch)
	startOnLaunchBinding.AddListener(binding.NewDataListener(func() {
		cfg.StartOnLaunch, _ = startOnLaunchBinding.Get()
		cfg.Save()
	}))

	inactiveStartBinding := binding.NewString()
	inactiveStartBinding.Set(cfg.InactiveStart)
	inactiveStartBinding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveStart, _ = inactiveStartBinding.Get()
		cfg.Save()
	}))

	inactiveEndBinding := binding.NewString()
	inactiveEndBinding.Set(cfg.InactiveEnd)
	inactiveEndBinding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnd, _ = inactiveEndBinding.Get()
		cfg.Save()
	}))

	focusDurationBinding := binding.NewString()
	focusDurationBinding.Set(fmt.Sprintf("%.0f", cfg.FocusDuration.Minutes()))
	focusDurationBinding.AddListener(binding.NewDataListener(func() {
		val, _ := focusDurationBinding.Get()
		mins, _ := strconv.Atoi(val)
		cfg.FocusDuration = time.Duration(mins) * time.Minute
		cfg.Save()
	}))

	shortBreakDurationBinding := binding.NewString()
	shortBreakDurationBinding.Set(fmt.Sprintf("%.0f", cfg.ShortBreakDuration.Minutes()))
	shortBreakDurationBinding.AddListener(binding.NewDataListener(func() {
		val, _ := shortBreakDurationBinding.Get()
		mins, _ := strconv.Atoi(val)
		cfg.ShortBreakDuration = time.Duration(mins) * time.Minute
		cfg.Save()
	}))

	longBreakDurationBinding := binding.NewString()
	longBreakDurationBinding.Set(fmt.Sprintf("%.0f", cfg.LongBreakDuration.Minutes()))
	longBreakDurationBinding.AddListener(binding.NewDataListener(func() {
		val, _ := longBreakDurationBinding.Get()
		mins, _ := strconv.Atoi(val)
		cfg.LongBreakDuration = time.Duration(mins) * time.Minute
		cfg.Save()
	}))

	settingsTab := container.NewVBox(
		widget.NewCheckWithData("Start on launch", startOnLaunchBinding),
		widget.NewLabel("Inactive Period"),
		container.NewHBox(
			widget.NewLabel("Start:"),
			widget.NewEntryWithData(inactiveStartBinding),
		),
		container.NewHBox(
			widget.NewLabel("End:"),
			widget.NewEntryWithData(inactiveEndBinding),
		),
		widget.NewLabel("Durations (minutes)"),
		container.NewHBox(
			widget.NewLabel("Focus:"),
			widget.NewEntryWithData(focusDurationBinding),
		),
		container.NewHBox(
			widget.NewLabel("Short Break:"),
			widget.NewEntryWithData(shortBreakDurationBinding),
		),
		container.NewHBox(
			widget.NewLabel("Long Break:"),
			widget.NewEntryWithData(longBreakDurationBinding),
		),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Pomodoro", pomodoroTab),
		container.NewTabItem("Settings", settingsTab),
	)

	if cfg.StartOnLaunch {
		startButton.OnTapped()
	}

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}

func formatTime(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func isInactive(cfg *config.Config) bool {
	now := time.Now()
	start, err := time.Parse("15:04", cfg.InactiveStart)
	if err != nil {
		return false
	}
	end, err := time.Parse("15:04", cfg.InactiveEnd)
	if err != nil {
		return false
	}
	return now.After(start) && now.Before(end)
}

func updateTitle(w fyne.Window, t *pomo.Timer) {
	emoji := "🍅"
	if t.State != pomo.Pomodoro {
		emoji = "🧘"
	}
	w.SetTitle(fmt.Sprintf("%s %s", emoji, config.AppName))
}

func updateSessionLabel(l *widget.Label, t *pomo.Timer) {
	text := "Pomodoro"
	if t.State == pomo.ShortBreakState {
		text = "Short Break"
	} else if t.State == pomo.LongBreakState {
		text = "Long Break"
	}
	l.SetText(text)
}

func getMediaPath(fileName string) string {
	executable, err := os.Executable()
	if err != nil {
		return filepath.Join("media", fileName)
	}
	dir := filepath.Dir(executable)
	mediaPath := filepath.Join(dir, "media", fileName)

	if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
		mediaPath = filepath.Join(dir, "..", "media", fileName)
	}

	return mediaPath
}
