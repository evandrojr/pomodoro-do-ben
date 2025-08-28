package gui

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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

	sessionBinding := binding.NewString()
	sessionBinding.Set(i18n.T("pomodoro"))
	sessionLabel := widget.NewLabelWithData(sessionBinding)
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

			// updateTitle(myWindow, timer) // FIXME: This is not thread-safe and causes crashes.

			newSessionText := i18n.T("pomodoro")
			if timer.State == pomo.ShortBreakState {
				newSessionText = i18n.T("break")
			} else if timer.State == pomo.LongBreakState {
				newSessionText = i18n.T("break")
			}
			sessionBinding.Set(newSessionText)

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

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(0, 150))
	pomodoroTab := container.NewVBox(topSpacer, timerLabel, sessionLabel, buttons)

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

	// --- Inactive Period Bindings and UI --- 

	nextDayLabel1 := widget.NewLabel("(" + i18n.T("next_day") + ")")
	nextDayLabel1.Hide()
	nextDayLabel2 := widget.NewLabel("(" + i18n.T("next_day") + ")")
	nextDayLabel2.Hide()

	validIcon1s := widget.NewIcon(theme.ConfirmIcon())
	validIcon1s.Hide()
	validIcon1e := widget.NewIcon(theme.ConfirmIcon())
	validIcon1e.Hide()
	validIcon2s := widget.NewIcon(theme.ConfirmIcon())
	validIcon2s.Hide()
	validIcon2e := widget.NewIcon(theme.ConfirmIcon())
	validIcon2e.Hide()

	checkTimeValidity := func(binding binding.String, icon *widget.Icon) {
		timeStr, _ := binding.Get()
		_, err := time.Parse("15:04", timeStr)
		if err == nil {
			icon.Show()
		} else {
			icon.Hide()
		}
	}

	checkOvernight := func(startBinding, endBinding binding.String, label *widget.Label) {
		startStr, _ := startBinding.Get()
		endStr, _ := endBinding.Get()
		s, err1 := time.Parse("15:04", startStr)
		e, err2 := time.Parse("15:04", endStr)

		if err1 == nil && err2 == nil && s.After(e) {
			label.Show()
		} else {
			label.Hide()
		}
	}

	inactiveEnabled1Binding := binding.NewBool()
	inactiveEnabled1Binding.Set(cfg.InactiveEnabled1)
	inactiveEnabled1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnabled1, _ = inactiveEnabled1Binding.Get()
		cfg.Save()
	}))

	inactiveStart1Binding := binding.NewString()
	inactiveStart1Binding.Set(cfg.InactiveStart1)

	inactiveEnd1Binding := binding.NewString()
	inactiveEnd1Binding.Set(cfg.InactiveEnd1)

	inactiveStart1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveStart1, _ = inactiveStart1Binding.Get()
		cfg.Save()
		checkTimeValidity(inactiveStart1Binding, validIcon1s)
		checkOvernight(inactiveStart1Binding, inactiveEnd1Binding, nextDayLabel1)
	}))
	inactiveEnd1Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnd1, _ = inactiveEnd1Binding.Get()
		cfg.Save()
		checkTimeValidity(inactiveEnd1Binding, validIcon1e)
		checkOvernight(inactiveStart1Binding, inactiveEnd1Binding, nextDayLabel1)
	}))

	inactiveEnabled2Binding := binding.NewBool()
	inactiveEnabled2Binding.Set(cfg.InactiveEnabled2)
	inactiveEnabled2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnabled2, _ = inactiveEnabled2Binding.Get()
		cfg.Save()
	}))

	inactiveStart2Binding := binding.NewString()
	inactiveStart2Binding.Set(cfg.InactiveStart2)

	inactiveEnd2Binding := binding.NewString()
	inactiveEnd2Binding.Set(cfg.InactiveEnd2)

	inactiveStart2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveStart2, _ = inactiveStart2Binding.Get()
		cfg.Save()
		checkTimeValidity(inactiveStart2Binding, validIcon2s)
		checkOvernight(inactiveStart2Binding, inactiveEnd2Binding, nextDayLabel2)
	}))
	inactiveEnd2Binding.AddListener(binding.NewDataListener(func() {
		cfg.InactiveEnd2, _ = inactiveEnd2Binding.Get()
		cfg.Save()
		checkTimeValidity(inactiveEnd2Binding, validIcon2e)
		checkOvernight(inactiveStart2Binding, inactiveEnd2Binding, nextDayLabel2)
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

	// Create Entry widgets and disable their default validators
	inactiveStart1Entry := widget.NewEntryWithData(inactiveStart1Binding)
	inactiveStart1Entry.Validator = nil
	inactiveEnd1Entry := widget.NewEntryWithData(inactiveEnd1Binding)
	inactiveEnd1Entry.Validator = nil
	inactiveStart2Entry := widget.NewEntryWithData(inactiveStart2Binding)
	inactiveStart2Entry.Validator = nil
	inactiveEnd2Entry := widget.NewEntryWithData(inactiveEnd2Binding)
	inactiveEnd2Entry.Validator = nil

	inactiveForm1 := widget.NewForm(
		widget.NewFormItem(i18n.T("start_time"), container.NewBorder(nil, nil, nil, validIcon1s, inactiveStart1Entry)),
		widget.NewFormItem(i18n.T("end_time"), container.NewBorder(nil, nil, nil, container.NewHBox(validIcon1e, nextDayLabel1), inactiveEnd1Entry)),
	)

	inactiveForm2 := widget.NewForm(
		widget.NewFormItem(i18n.T("start_time"), container.NewBorder(nil, nil, nil, validIcon2s, inactiveStart2Entry)),
		widget.NewFormItem(i18n.T("end_time"), container.NewBorder(nil, nil, nil, container.NewHBox(validIcon2e, nextDayLabel2), inactiveEnd2Entry)),
	)

	durationForm := widget.NewForm(
		widget.NewFormItem(i18n.T("focus_duration"), widget.NewEntryWithData(focusDurationBinding)),
		widget.NewFormItem(i18n.T("short_break_duration"), widget.NewEntryWithData(shortBreakDurationBinding)),
		widget.NewFormItem(i18n.T("long_break_duration"), widget.NewEntryWithData(longBreakDurationBinding)),
	)

	settingsTab := container.NewVBox(
		widget.NewCheckWithData(i18n.T("start_on_launch"), startOnLaunchBinding),
		widget.NewCheckWithData(i18n.T("auto_start_cycles"), autoStartCyclesBinding),
		widget.NewSeparator(),
		widget.NewCheckWithData(i18n.T("inactive_period_1"), inactiveEnabled1Binding),
		inactiveForm1,
		widget.NewSeparator(),
		widget.NewCheckWithData(i18n.T("inactive_period_2"), inactiveEnabled2Binding),
		inactiveForm2,
		widget.NewLabel(i18n.T("next_day_tip")),
		widget.NewSeparator(),
		widget.NewLabel(i18n.T("durations_in_minutes")),
		durationForm,
	)

	// Perform initial check on load
	checkTimeValidity(inactiveStart1Binding, validIcon1s)
	checkTimeValidity(inactiveEnd1Binding, validIcon1e)
	checkTimeValidity(inactiveStart2Binding, validIcon2s)
	checkTimeValidity(inactiveEnd2Binding, validIcon2e)
	checkOvernight(inactiveStart1Binding, inactiveEnd1Binding, nextDayLabel1)
	checkOvernight(inactiveStart2Binding, inactiveEnd2Binding, nextDayLabel2)

	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("pomodoro"), pomodoroTab),
		container.NewTabItem(i18n.T("settings"), settingsTab),
	)

	if cfg.StartOnLaunch {
		startButton.OnTapped()
	}

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

func formatTime(d time.Duration) string {
	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

func isInactive(cfg *config.Config) bool {
	now := time.Now()

	check := func(enabled bool, startStr, endStr string) bool {
		if !enabled {
			return false
		}
		start, err1 := time.Parse("15:04", startStr)
		end, err2 := time.Parse("15:04", endStr)
		if err1 != nil || err2 != nil {
			return false // Don't be inactive if times are invalid
		}

		// Handle overnight period
		if start.After(end) {
			if now.After(start) || now.Before(end) {
				return true
			}
		} else {
			// Handle same-day period
			if now.After(start) && now.Before(end) {
				return true
			}
		}
		return false
	}

	if check(cfg.InactiveEnabled1, cfg.InactiveStart1, cfg.InactiveEnd1) {
		return true
	}
	if check(cfg.InactiveEnabled2, cfg.InactiveStart2, cfg.InactiveEnd2) {
		return true
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
