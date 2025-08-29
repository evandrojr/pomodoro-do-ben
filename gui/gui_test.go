package gui

import (
	"testing"
	"time"

	"pomodoro-do-ben/config"
)

func TestIsInactive(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		now  time.Time
		expected bool
	}{
		{
			name: "Inactive period 1 disabled",
			cfg: &config.Config{
				InactiveEnabled1: false,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
			},
			now:      time.Date(2025, time.January, 1, 9, 30, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 1 enabled, within period (same day)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
			},
			now:      time.Date(2025, time.January, 1, 9, 30, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Inactive period 1 enabled, outside period (same day - before)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
			},
			now:      time.Date(2025, time.January, 1, 8, 30, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 1 enabled, outside period (same day - after)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
			},
			now:      time.Date(2025, time.January, 1, 10, 30, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 1 enabled, spanning overnight (within)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "22:00",
				InactiveEnd1:     "06:00",
			},
			now:      time.Date(2025, time.January, 1, 23, 0, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Inactive period 1 enabled, spanning overnight (within - next day)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "22:00",
				InactiveEnd1:     "06:00",
			},
			now:      time.Date(2025, time.January, 2, 5, 0, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Inactive period 1 enabled, spanning overnight (outside - before)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "22:00",
				InactiveEnd1:     "06:00",
			},
			now:      time.Date(2025, time.January, 1, 21, 0, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 1 enabled, spanning overnight (outside - after)",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "22:00",
				InactiveEnd1:     "06:00",
			},
			now:      time.Date(2025, time.January, 2, 7, 0, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 1 enabled, invalid time format",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "invalid",
				InactiveEnd1:     "10:00",
			},
			now:      time.Date(2025, time.January, 1, 9, 30, 0, 0, time.Local),
			expected: false,
		},
		{
			name: "Inactive period 2 enabled, within period",
			cfg: &config.Config{
				InactiveEnabled1: false,
				InactiveEnabled2: true,
				InactiveStart2:   "14:00",
				InactiveEnd2:     "15:00",
			},
			now:      time.Date(2025, time.January, 1, 14, 30, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Both inactive periods enabled, first one active",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
				InactiveEnabled2: true,
				InactiveStart2:   "14:00",
				InactiveEnd2:     "15:00",
			},
			now:      time.Date(2025, time.January, 1, 9, 30, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Both inactive periods enabled, second one active",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
				InactiveEnabled2: true,
				InactiveStart2:   "14:00",
				InactiveEnd2:     "15:00",
			},
			now:      time.Date(2025, time.January, 1, 14, 30, 0, 0, time.Local),
			expected: true,
		},
		{
			name: "Both inactive periods enabled, neither active",
			cfg: &config.Config{
				InactiveEnabled1: true,
				InactiveStart1:   "09:00",
				InactiveEnd1:     "10:00",
				InactiveEnabled2: true,
				InactiveStart2:   "14:00",
				InactiveEnd2:     "15:00",
			},
			now:      time.Date(2025, time.January, 1, 12, 0, 0, 0, time.Local),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock time.Now() for testing purposes
			originalTimeNow := timeNow
			timeNow = func() time.Time { return tt.now }
			defer func() { timeNow = originalTimeNow }()

			result := isInactive(tt.cfg)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for %s", tt.expected, result, tt.name)
			}
		})
	}
}
