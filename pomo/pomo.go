package pomo

import (
	"time"

	"pomodoro-do-ben/config"
)

type State int

const (
	Pomodoro State = iota
	ShortBreakState
	LongBreakState
)

type Timer struct {
	State         State
	Duration      time.Duration
	RemainingTime time.Duration
	IsRunning     bool
	ticker        *time.Ticker
	config        *config.Config
	Updates       chan struct{}
	pomodoroCount int // Tracks completed pomodoros
}

func NewTimer(cfg *config.Config) *Timer {
	return &Timer{
		State:         Pomodoro,
		Duration:      cfg.FocusDuration,
		RemainingTime: cfg.FocusDuration,
		config:        cfg,
		Updates:       make(chan struct{}),
		pomodoroCount: 0, // Initialize pomodoro count
	}
}

func (t *Timer) Start() {
	t.IsRunning = true
	t.ticker = time.NewTicker(time.Second)
	go func() {
		for range t.ticker.C {
			t.Tick()
			t.Updates <- struct{}{}
		}
	}()
}

func (t *Timer) Stop() {
	t.IsRunning = false
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.RemainingTime = t.Duration
}

func (t *Timer) Tick() {
	if t.IsRunning {
		t.RemainingTime -= time.Second
	}
}

func (t *Timer) Ticker() *time.Ticker {
	return t.ticker
}

func (t *Timer) NextState() {
	switch t.State {
	case Pomodoro:
		t.pomodoroCount++ // Increment pomodoro count after a completed Pomodoro
		if t.pomodoroCount%t.config.LongBreakInterval == 0 {
			t.State = LongBreakState
			t.Duration = t.config.LongBreakDuration
		} else {
			t.State = ShortBreakState
			t.Duration = t.config.ShortBreakDuration
		}
	case ShortBreakState:
		t.State = Pomodoro
		t.Duration = t.config.FocusDuration
	case LongBreakState:
		t.State = Pomodoro
		t.Duration = t.config.FocusDuration
		t.pomodoroCount = 0 // Reset pomodoro count after a long break
	}
	t.Reset()
}
