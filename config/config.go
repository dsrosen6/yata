package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type ConfigIn struct {
	MainColor      *uint `json:"main_color"`
	SecondaryColor *uint `json:"secondary_color"`
}

type Config struct {
	MainColor      lipgloss.ANSIColor
	SecondaryColor lipgloss.ANSIColor
}

var defaultConfig = Config{
	MainColor:      lipgloss.ANSIColor(4),
	SecondaryColor: lipgloss.ANSIColor(7),
}

const (
	cfgDirName  = "yata"
	cfgFileName = "config.json"
)

func GetConfig() (*Config, error) {
	uc, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("getting user config directory path: %w", err)
	}
	path := filepath.Join(uc, cfgDirName, cfgFileName)
	cfg, err := readConfig(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	return configInToConfig(cfg), nil
}

func configInToConfig(in *ConfigIn) *Config {
	cfg := defaultConfig
	if in == nil {
		return &cfg
	}

	if in.MainColor != nil {
		cfg.MainColor = lipgloss.ANSIColor(*in.MainColor)
	}

	if in.SecondaryColor != nil {
		cfg.SecondaryColor = lipgloss.ANSIColor(*in.SecondaryColor)
	}

	return &cfg
}

func readConfig(path string) (*ConfigIn, error) {
	cfg := &ConfigIn{}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("checking for config file: %w", err)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := json.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling json: %w", err)
	}

	return cfg, nil
}
