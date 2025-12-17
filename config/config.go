package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

type ConfigIn struct {
	MainColor      *int `json:"main_color"`
	SecondaryColor *int `json:"secondary_color"`
}

type Config struct {
	MainColor      tcell.Color
	SecondaryColor tcell.Color
}

var defaultConfig = Config{
	MainColor:      tcell.ColorBlue,
	SecondaryColor: tcell.ColorDefault,
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
		c := tcell.PaletteColor(*in.MainColor)
		if c.Valid() {
			cfg.MainColor = c
		}
	}

	if in.SecondaryColor != nil {
		c := tcell.PaletteColor(*in.SecondaryColor)
		if c.Valid() {
			cfg.SecondaryColor = c
		}
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
