package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	AppName = "Pomodoro do Ben"
)

type Config struct {
	StartOnLaunch      bool          `json:"start_on_launch"`
	AutoStartCycles    bool          `json:"auto_start_cycles"`
	InactiveStart      string        `json:"inactive_start"`
	InactiveEnd        string        `json:"inactive_end"`
	FocusDuration      time.Duration `json:"focus_duration"`
	ShortBreakDuration time.Duration `json:"short_break_duration"`
	LongBreakDuration  time.Duration `json:"long_break_duration"`
}

func Load() (*Config, error) {
	cfg := &Config{
		StartOnLaunch:      true,
		AutoStartCycles:    true,
		InactiveStart:      "13:00",
		InactiveEnd:        "14:00",
		FocusDuration:      25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
	}

	path, err := configPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(c)
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, AppName, "config.json"), nil
}
