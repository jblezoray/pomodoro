package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WorkDuration             int    `json:"work_duration"`
	ShortBreak               int    `json:"short_break"`
	LongBreak                int    `json:"long_break"`
	PomodorosBeforeLongBreak int    `json:"pomodoros_before_long_break"`
	AutoStartBreaks          bool   `json:"auto_start_breaks"`
	AutoStartPomodoros       bool   `json:"auto_start_pomodoros"`
	SoundEnabled             bool   `json:"sound_enabled"`
	SoundFile                string `json:"sound_file"`
	WorkLabel                string `json:"work_label"`
	ShortBreakLabel          string `json:"short_break_label"`
	LongBreakLabel           string `json:"long_break_label"`
}

var defaultConfig = Config{
	WorkDuration:             25,
	ShortBreak:               5,
	LongBreak:                15,
	PomodorosBeforeLongBreak: 4,
	AutoStartBreaks:          false,
	AutoStartPomodoros:       false,
	SoundEnabled:             true,
	WorkLabel:                "Travail",
	ShortBreakLabel:          "Pause courte",
	LongBreakLabel:           "Grande pause",
}

func loadConfig() Config {
	candidates := []string{"pomodoro.json"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(home, ".pomodoro.json"),
			filepath.Join(home, ".config", "pomodoro", "config.json"),
		)
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		applyDefaults(&cfg)
		return cfg
	}
	return defaultConfig
}

func applyDefaults(cfg *Config) {
	if cfg.WorkDuration == 0 {
		cfg.WorkDuration = defaultConfig.WorkDuration
	}
	if cfg.ShortBreak == 0 {
		cfg.ShortBreak = defaultConfig.ShortBreak
	}
	if cfg.LongBreak == 0 {
		cfg.LongBreak = defaultConfig.LongBreak
	}
	if cfg.PomodorosBeforeLongBreak == 0 {
		cfg.PomodorosBeforeLongBreak = defaultConfig.PomodorosBeforeLongBreak
	}
	if cfg.WorkLabel == "" {
		cfg.WorkLabel = defaultConfig.WorkLabel
	}
	if cfg.ShortBreakLabel == "" {
		cfg.ShortBreakLabel = defaultConfig.ShortBreakLabel
	}
	if cfg.LongBreakLabel == "" {
		cfg.LongBreakLabel = defaultConfig.LongBreakLabel
	}
}
