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
	"pomodoro-do-ben/i18n"
	"pomodoro-do-ben/notifier"
	"pomodoro-do-ben/player"
	"pomodoro-do-ben/pomo"
)

func Show(cfg *config.Config) {
	myApp := app.New()
	myWindow := myApp.NewWindow(i18n.T("bens_pomodoro"))

	timer := pomo.NewTimer(cfg)

	timerStr := binding.NewString()
	timerStr.Set(formatTime(timer.RemainingTime))

	timerLabel := widget.NewLabelWithData(timerStr)
	timerLabel.Alignment = fyne.TextAlignCenter
	timerLabel.TextStyle = fyne.TextStyle{Bold: true}

	sessionLabel := widget.NewLabel(i18n.T("pomodoro"))
	sessionLabel.Alignment = fyne.TextAlignCenter

	startButton := widget.NewButtonWithIcon("▶️ "+i18n.T("start"), theme.MediaPlayIcon(), func() {
		if !timer.IsRunning {
			if isInactive(cfg) {
				notifier.Notify(i18n.T("pomodoro"), "Timer is inactive during this period.")
				return
			}
			timer.Start()
			player.Play(getMediaPath("focar/f1.aac"))
			notifier.Notify(i18n.T("pomodoro"), i18n.T("time_to_focus"))
		}
	})

	pauseButton := widget.NewButtonWithIcon("⏸️ "+i18n.T("pause"), theme.MediaPauseIcon(), func() {
		if timer.IsRunning {
			timer.Stop()
			notifier.Notify(i18n.T("pomodoro"), "Timer paused!")
		}
	})

	resetButton := widget.NewButtonWithIcon("🔄 "+i18n.T("stop"), theme.MediaReplayIcon(), func() {
		timer.Reset()
		timerStr.Set(formatTime(timer.RemainingTime))
		notifier.Notify(i18n.T("pomodoro"), "Timer reset!")
	})

	go func() {
		for range timer.Updates {
			timerStr.Set(formatTime(timer.RemainingTime))
			updateTitle(myWindow, timer)
			updateSessionLabel(sessionLabel, timer)
			myWindow.Canvas().Refresh(myWindow.Content())
			if timer.RemainingTime <= 0 {
				timer.NextState()
				if cfg.AutoStartCycles {
					timer.Start()
				}
				player.Play(getMediaPath("meditar/m1.aac"))
				notifier.Notify(i18n.T("pomodoro"), i18n.T("time_to_break"))
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

	autoStartCyclesBinding := binding.NewBool()
	autoStartCyclesBinding.Set(cfg.AutoStartCycles)
	autoStartCyclesBinding.AddListener(binding.NewDataListener(func() {
		cfg.AutoStartCycles, _ = autoStartCyclesBinding.Get()
		cfg.Save()
	}))

	inactiveEnabled1Binding := binding.NewBool()
	inactiveEnabled1Binding.Set(cfg.InactiveEnabled1)
	inactiveEnabled1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnabled1, _ = inactiveEnabled1Binding.Get()
		cfg.Save()
	}))

	inactiveStart1Binding := binding.NewString()
	inactiveStart1Binding.Set(cfg.InactiveStart1)
	inactiveStart1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveStart1, _ = inactiveStart1Binding.Get()
		cfg.Save()
	}))

	inactiveEnd1Binding := binding.NewString()
	inactiveEnd1Binding.Set(cfg.InactiveEnd1)
	inactiveEnd1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnd1, _ = inactiveEnd1Binding.Get()
		cfg.Save()
	}))

	inactiveEnabled2Binding := binding.NewBool()
	inactiveEnabled2Binding.Set(cfg.InactiveEnabled2)
	inactiveEnabled2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnabled2, _ = inactiveEnabled2Binding.Get()
		cfg.Save()
	}))

	inactiveStart2Binding := binding.NewString()
	inactiveStart2Binding.Set(cfg.InactiveStart2)
	inactiveStart2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveStart2, _ = inactiveStart2Binding.Get()
		cfg.Save()
	}))

	inactiveEnd2Binding := binding.NewString()
	inactiveEnd2Binding.Set(cfg.InactiveEnd2)
	inactiveEnd2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnd2, _ = inactiveEnd2Binding.Get()
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
		widget.NewCheckWithData(i18n.T("start_on_launch"), startOnLaunchBinding),
		widget.NewCheckWithData(i18n.T("auto_start_cycles"), autoStartCyclesBinding),
		widget.NewCheckWithData(i18n.T("inactive_period_1"), inactiveEnabled1Binding),
		container.NewHBox(
			widget.NewLabel(i18n.T("start_time")),
			widget.NewEntryWithData(inactiveStart1Binding),
		),
		container.NewHBox(
			widget.NewLabel(i18n.T("end_time")),
			widget.NewEntryWithData(inactiveEnd1Binding),
		),
		widget.NewCheckWithData(i18n.T("inactive_period_2"), inactiveEnabled2Binding),
		container.NewHBox(
			widget.NewLabel(i18n.T("start_time")),
			widget.NewEntryWithData(inactiveStart2Binding),
		),
		container.NewHBox(
			widget.NewLabel(i18n.T("end_time")),
			widget.NewEntryWithData(inactiveEnd2Binding),
		),
		widget.NewLabel(i18n.T("durations_in_minutes")),
		container.NewHBox(
			widget.NewLabel(i18n.T("focus_duration")),
			widget.NewEntryWithData(focusDurationBinding),
		),
		container.NewHBox(
			widget.NewLabel(i18n.T("short_break_duration")),
			widget.NewEntryWithData(shortBreakDurationBinding),
		),
		container.NewHBox(
			widget.NewLabel(i18n.T("long_break_duration")),
			widget.NewEntryWithData(longBreakDurationBinding),
		),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("pomodoro"), pomodoroTab),
		container.NewTabItem(i18n.T("settings"), settingsTab),
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
	if cfg.InactiveEnabled1 {
		start, err := time.Parse("15:04", cfg.InactiveStart1)
		if err != nil {
			return false
		}
		end, err := time.Parse("15:04", cfg.InactiveEnd1)
		if err != nil {
			return false
		}
		if now.After(start) && now.Before(end) {
			return true
		}
	}
	if cfg.InactiveEnabled2 {
		start, err := time.Parse("15:04", cfg.InactiveStart2)
		if err != nil {
			return false
		}
		end, err := time.Parse("15:04", cfg.InactiveEnd2)
		if err != nil {
			return false
		}
		if now.After(start) && now.Before(end) {
			return true
		}
	}
	return false
}

func updateTitle(w fyne.Window, t *pomo.Timer) {
	emoji := "🍅"
	if t.State != pomo.Pomodoro {
		emoji = "🧘"
	}
	w.SetTitle(fmt.Sprintf("%s %s", emoji, i18n.T("bens_pomodoro")))
}

func updateSessionLabel(l *widget.Label, t *pomo.Timer) {
	text := i18n.T("pomodoro")
	if t.State == pomo.ShortBreakState {
		text = i18n.T("break")
	} else if t.State == pomo.LongBreakState {
		text = i18n.T("break")
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
