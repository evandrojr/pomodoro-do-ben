package pomo

import (
	"testing"
	"time"

	"pomodoro-do-ben/config"
)

func TestNewTimer(t *testing.T) {
	cfg := &config.Config{
		FocusDuration: time.Minute * 25,
	}
	timer := NewTimer(cfg)

	if timer.State != Pomodoro {
		t.Errorf("Expected initial state to be Pomodoro, got %v", timer.State)
	}
	if timer.Duration != cfg.FocusDuration {
		t.Errorf("Expected initial duration to be FocusDuration, got %v", timer.Duration)
	}
	if timer.RemainingTime != cfg.FocusDuration {
		t.Errorf("Expected initial remaining time to be FocusDuration, got %v", timer.RemainingTime)
	}
	if timer.IsRunning != false {
		t.Errorf("Expected IsRunning to be false, got %v", timer.IsRunning)
	}
	if timer.pomodoroCount != 0 {
		t.Errorf("Expected initial pomodoroCount to be 0, got %v", timer.pomodoroCount)
	}
}

func TestTimerStartStopReset(t *testing.T) {
	cfg := &config.Config{
		FocusDuration: time.Second * 2,
	}
	timer := NewTimer(cfg)

	timer.Start()
	if !timer.IsRunning {
		t.Error("Expected timer to be running after Start()")
	}

	// Give it a moment to tick
	time.Sleep(time.Millisecond * 1500)

	timer.Stop()
	if timer.IsRunning {
		t.Error("Expected timer to be stopped after Stop()")
	}

	timer.Reset()
	if timer.RemainingTime != cfg.FocusDuration {
		t.Errorf("Expected remaining time to reset to FocusDuration, got %v", timer.RemainingTime)
	}
}

func TestTimerTick(t *testing.T) {
	cfg := &config.Config{
		FocusDuration: time.Second * 5,
	}
	timer := NewTimer(cfg)

	timer.IsRunning = true // Manually set to running for Tick test
	timer.Tick()
	if timer.RemainingTime != time.Second*4 {
		t.Errorf("Expected remaining time to decrease by 1 second, got %v", timer.RemainingTime)
	}

	timer.IsRunning = false // Should not tick when not running
	timer.Tick()
	if timer.RemainingTime != time.Second*4 {
		t.Errorf("Expected remaining time to not change when not running, got %v", timer.RemainingTime)
	}
}

func TestNextState(t *testing.T) {
	cfg := &config.Config{
		FocusDuration:      time.Minute * 25,
		ShortBreakDuration: time.Minute * 5,
		LongBreakDuration:  time.Minute * 15,
	}

	tests := []struct {
		name          string
		initialState  State
		expectedState State
		expectedDuration time.Duration
	}{
		{
			name:          "Pomodoro to ShortBreak",
			initialState:  Pomodoro,
			expectedState: ShortBreakState,
			expectedDuration: cfg.ShortBreakDuration,
		},
		{
			name:          "ShortBreak to Pomodoro",
			initialState:  ShortBreakState,
			expectedState: Pomodoro,
			expectedDuration: cfg.FocusDuration,
		},
		{
			name:          "LongBreak to Pomodoro",
			initialState:  LongBreakState,
			expectedState: Pomodoro,
			expectedDuration: cfg.FocusDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := NewTimer(cfg)
			timer.State = tt.initialState
			timer.NextState()

			if timer.State != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, timer.State)
			}
			if timer.Duration != tt.expectedDuration {
				t.Errorf("Expected duration %v, got %v", tt.expectedDuration, timer.Duration)
			}
			if timer.RemainingTime != tt.expectedDuration {
				t.Errorf("Expected remaining time to reset to %v, got %v", tt.expectedDuration, timer.RemainingTime)
			}
		})
	}
}
